package runner

import (
	pb "Src1/proto"
	"Src1/utils"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Requester struct {
	reporter *Reporter
	config   *RunConfig
	start    time.Time
	End      time.Time
	//stopwatch *utils.StopWatch
	//uuid string
}

// NewRequester creates a new requestor from the passed RunConfig
func NewRequester(c *RunConfig) *Requester {
	r := &Requester{
		config: c,
	}
	return r
}

func (report *SubWorker) ErrorHandler(code string, message string, err error) *SubWorker {
	report.Error.ErrorCode = code
	report.Error.ErrorMessage = message
	report.Error.Error = err
	return report
}

func (report *SubWorker) GetError() Error {

	return report.Error
}

func (b *Requester) Run() (*MainWorker, error) {
	//runWorkers()
	err := ConnectionCheck(b.config)
	if err != nil {
		return nil, err
	}

	b.start = time.Now()
	reporter, err := b.runWorkers()
	if err != nil {
		return nil, err
	}
	b.End = time.Now()
	return reporter, nil
}

func (b *Requester) worker(wID string) *SubWorker {
	reporter := newReporter(wID)
	ctx, cancel := context.WithTimeout(context.Background(), b.config.Timeout*time.Second)
	defer cancel()
	sw := utils.NewStopWatch()
	tlsCredentials, err := LoadTLSCredentials(b.config.SkipVerify, b.config.Cert)
	opts := GrpcOption(b.config, tlsCredentials)
	reporter.Start = sw.Start()
	fmt.Printf("%v", b.config)
	conn, err := grpc.DialContext(ctx,
		fmt.Sprintf("%s:%d", b.config.Host, b.config.Port),
		opts...,
	)
	if err != nil {
		return reporter.ErrorHandler("did not connect", "DialContext", err)
	}
	if !b.config.Block {
		checkConnectivityStatusChan(ctx, conn, connectivity.Idle)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
		}
	}(conn)

	c := pb.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "hello"})
	if err != nil {
		reporter.ErrorHandler("Procedure call Error", fmt.Sprintf("could not greet: %v", err), err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
	reporter.EndTime, reporter.End = sw.Stop()
	return reporter
}

func (b *Requester) runWorkers() (*MainWorker, error) {
	reporter := newLoadReporter()
	var wg sync.WaitGroup
	wg.Add(b.config.TotalRequest)
	n := 0
	wc := 0
	sleep := time.Second / time.Duration(b.config.Rps)
	for i := 0; i < b.config.TotalRequest; i++ {
		go func(b *Requester) {
			wID := "g" + strconv.Itoa(wc) + "c" + strconv.Itoa(n)
			worker := b.worker(wID)
			reporter.addWorker(*worker)
			wg.Done()
		}(b)
		time.Sleep(sleep)
	}
	wg.Wait()
	reporter.TotalCount = b.config.TotalRequest
	reporter.EndTime = time.Now()
	reporter.FinishTime = time.Since(b.start)
	reporter.Finish()
	return reporter, nil
}

func (r *MainWorker) Finish() {
	// Slowest / Fastest
	var LatenciesSec float64
	var okLats []float64
	var errReport []Error
	var SuccssCount uint64
	errCount := 0
	for _, worker := range r.Workers {

		okLats = append(okLats, float64(worker.EndTime.Seconds()))
		if worker.Error.ErrorCode != "" {
			errCount += 1
			errReport = append(errReport, worker.GetError())
		}
		SuccssCount += 1
	}
	LatenciesSec = float64(r.FinishTime.Seconds())
	r.SuccssCount = SuccssCount
	r.Rps = int(float64(r.TotalCount) / LatenciesSec)
	sort.Float64s(okLats)
	if len(okLats) > 0 {
		var fastestNum, slowestNum float64
		fastestNum = okLats[0]
		slowestNum = okLats[len(okLats)-1]

		r.Fastest = time.Duration(fastestNum * float64(time.Second))
		r.Slowest = time.Duration(slowestNum * float64(time.Second))
		r.Histogram = histogram(okLats, slowestNum, fastestNum)
		r.LatencyDistribution = latencies(okLats)
	}

	average := LatenciesSec / float64(r.TotalCount)
	r.Average = time.Duration(average * float64(time.Second))
	r.ErrorCount = errCount
	r.ErrorReport = errReport

}
