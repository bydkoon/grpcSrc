package runner

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"time"
)

type key int

const (
	requestIDKey key = 0
)

var (
	listenAddr string
	healthy    int32
)

func NewHttpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * time.Duration(10),
		Transport: &http.Transport{
			DisableKeepAlives: false,
			IdleConnTimeout:   time.Minute,
		},
	}
}

func ConnectionCheck(c *RunConfig) error {
	//log.Printf("%s", status)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	tlsCredentials, err := LoadTLSCredentials(c.SkipVerify, c.Cert)
	opts := GrpcOption(c, tlsCredentials)
	_, err = grpc.DialContext(ctx,
		fmt.Sprintf("%s:%d", c.Host, c.Port),
		opts...,
	)
	defer cancel()
	if err != nil {
		return err
	}
	log.Println("Server started.")
	return nil

}

// REQUIRES: conn is a pointer to a valid, open TCP connection
// MODIFIES: conn
// EFFECTS:	 Writes a sample message to the connection.
func handleConn(conn *net.TCPConn) {

	for {
		connIsClosed(conn)
		sampleMessage := []byte("Hello!\n")
		_, err := conn.Write(sampleMessage)
		checkErr(err)
		time.Sleep(1000 * time.Millisecond)
	}
}

// REQUIRES: conn is a pointer to a valid, open TCP connection
// EFFECTS:  Logs the event where a client joins.
func logClientJoined(conn *net.TCPConn) {
	log.Println("server.go: Client joined from %s", conn.RemoteAddr())
}

// EFFECTS:	 Handles any non-nil errors by printing them.
func checkErr(err error) {
	if err != nil {
		log.Println("Error: server.go: %s", err.Error())
	}
}

func connIsClosed(c *net.TCPConn) {
	c.SetReadDeadline(time.Now())
	var one []byte
	if _, err := c.Read(one); err == io.EOF {
		log.Println("Client disconnect: %s", c.RemoteAddr())
		c.Close()
		c = nil
	} else {
		var zero time.Time
		c.SetReadDeadline(zero)
	}
}

func ServerCheck4(c *RunConfig) error {
	client := NewHttpClient()
	url := fmt.Sprintf("http://%s:%d", c.Host, c.Port)
	_, err := client.Head(url)

	if err != nil {
		return err
	}
	return nil
}

func ServerCheck3(c *RunConfig) error {
	server := &http.Server{
		Addr:         listenAddr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	listenAddr = fmt.Sprintf("%s:%d", c.Host, c.Port)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	return server.Serve(ln)
}

func ServerCheck2(c *RunConfig) error {
	listenAddr = fmt.Sprintf("%s:%d", c.Host, c.Port)
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is starting...")

	router := http.NewServeMux()
	router.Handle("/", index())
	router.Handle("/healthz", healthz())

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:         listenAddr,
		Handler:      tracing(nextRequestID)(logging(logger)(router)),
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready to handle requests at", listenAddr)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
		return err
	}

	<-done
	logger.Println("Server stopped")
	return nil

}

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, World!")
	})
}

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&healthy) == 1 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Println(requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
