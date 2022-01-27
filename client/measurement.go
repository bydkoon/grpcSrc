package main

import (
	"Src1/client/cmd"
	"Src1/client/driver"
	errors2 "Src1/client/errors"
	"Src1/client/printer"
	pb "Src1/proto"

	"Src1/client/tls"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/alecthomas/kingpin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

const (
	defaultHost    = "localhost"
	defaultPort    = "50051"
	defaultMessage = "hello man"
	defaultCertPem = "C:\\Users\\K\\gopath\\src\\Src1\\cert\\ca-cert.pem"
)

var (
	skipVerify = kingpin.Flag("skipTLS", "Skip TLS client verification of the server's certificate chain and host name.").
			Default("false").Short('s').Bool()

	host = kingpin.Flag("host", "Hostname").Default(defaultHost).Short('h').String()

	port = kingpin.Flag("port", "Port number").Default(defaultPort).Short('p').Int()

	certPem = kingpin.Flag("cert", "TLS cert.pem").Default(defaultCertPem).PlaceHolder(" ").String()

	totalCount = kingpin.Flag("totalCount", "total count").Default("1").Short('t').Int()

	blockMode = kingpin.Flag("blockMode", "Dial BlockMode").Default("true").Short('b').Bool()

	timeOut = kingpin.Flag("timeOut", "Time Out option").Default("1").Int()
)

// 프로그램 실행시 호출
func init() {
	// 커맨드 라인 명령: cmd> *.exe -name [value] : https://gobyexample.com/command-line-flags
	//log.SetFlags(log.LstdFlags | log.Lshortfile)
	//glog := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	//grpclog.SetLoggerV2(glog)
	kingpin.Parse()

}

func main() {

	//sw := utils.NewStopWatch()
	//sw.Start()
	/*	pid := getGID()*/

	var opts []grpc.DialOption
	var cfg cmd.Config
	var report printer.Report
	report.TrackReport.Init()
	args := os.Args[1:]

	if len(args) > 1 {
		//var cmdCfg cmd.Config
		err := createConfigFromArgs(&cfg)
		if err != nil {
			kingpin.FatalIfError(err, "")
			os.Exit(1)
		}

	}
	start := time.Now()
	fullAddr := fmt.Sprintf("%s:%v", cfg.Host, cfg.Port)
	report.MakeReport(fullAddr, &cfg, start)

	wg := new(sync.WaitGroup)
	//for i :=0 ; i < cfg.TotalCount ; i++ {

	wg.Add(cfg.TotalCount)
	for i := 0; i < cfg.TotalCount; i++ {
		go func() {
			err := worker(opts, cfg, report)
			report.PrinterReport(err, &cfg)
			defer wg.Done()
		}()
	}

	wg.Wait()

}

// Check the change of connectivity state.
// Wait for state change, then notify and recurse
func checkConnectivityStatusChan(ctx context.Context, conn *grpc.ClientConn, sourceState connectivity.State) {
	ch := make(chan bool, 1)
	ch <- conn.WaitForStateChange(ctx, sourceState)

	select {
	case <-ctx.Done():
		log.Println("Context is Done")
	case <-ch:
		curState := conn.GetState()
		log.Printf("Change channel state : %s > %s [%s]\r\n", sourceState.String(), curState.String(), time.Now())
		close(ch)
		go checkConnectivityStatusChan(ctx, conn, curState)
	}
}

func createConfigFromArgs(config *cmd.Config) error {
	if config == nil {
		return errors.New("config cannot be nil")
	}
	config.Host = *host
	config.Port = *port
	config.SkipTLSVerify = *skipVerify
	config.CertPem = *certPem
	config.TotalCount = *totalCount
	config.TimeOut = *timeOut
	config.BlockMode = *blockMode
	//config.KeyPem = *keyPem
	return nil
}

func worker(opts []grpc.DialOption, cfg cmd.Config, report printer.Report) *errors2.Errors {
	start_time := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.TimeOut)*time.Second)
	defer cancel()

	tlsCredentials, err := tls.LoadTLSCredentials(&cfg)
	opts = driver.GrpcOption(opts, &cfg, tlsCredentials)
	conn, err := grpc.DialContext(ctx,
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		opts...,
	)
	if err != nil {
		return errors2.New("did not connect", "DialContext", err)
		//log.Fatalf("did not connect: %v", err)
		//return fmt.Errorf("did not connect: %v", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	ctx.Done()
	go checkConnectivityStatusChan(ctx, conn, connectivity.Idle)
	//conn, err := driver.Dial(address, driver.WithInsecure(), driver.WithBlock())
	conn.GetState()
	c := pb.NewGreeterClient(conn)

	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: defaultMessage})
	if err != nil {
		return errors2.New("could not greet:", "SayHello", err)
	}

	report.TrackReport.TrackCheck("sendEnd", time.Since(start_time))
	log.Printf("Greeting: %s", r.GetMessage())
	return nil
}

//func mergeConfig(config *cmd.Config, src *cmd.Config) errors {
//	if src == nil || config == nil {
//		return errors.New("config cannot be nil")
//	}
//
//	if isHostSet {
//		config.Host = src.Host
//	}
//	if isPortSet {
//		config.Port = src.Port
//	}
//
//	if isSkipSet {
//		config.SkipTLSVerify = src.SkipTLSVerify
//	}
//
//	if isCertSet {
//		config.CertPem = src.CertPem
//	}
//
//	//if isKeySet {
//	//	config.KeyPem = src.KeyPem
//	//}
//	return nil
//}
