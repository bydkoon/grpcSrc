package runner

import (
	"Src1/utils"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"strconv"
	"sync"
	"time"
)

type callResult struct {
	err       error
	status    string
	duration  time.Duration
	timestamp time.Time
}

type Requester struct {
	reporter  *Reporter
	config    *RunConfig
	start     time.Time
	End       time.Time
	stopwatch *utils.StopWatch
	uuid      string
}

// NewRequester creates a new requestor from the passed RunConfig
func NewRequester(c *RunConfig) *Requester {
	uid := uuid.NewString()
	reqr := &Requester{
		config:    c,
		uuid:      uid,
		stopwatch: utils.NewStopWatchUUID(uid),
	}

	return reqr
}

func (report *Worker) ErrorHandler(code string, message string, err error) *Worker {
	report.ErrorCode = code
	report.ErrorMessage = message
	report.Error = err
	return report
}
func (b *Requester) Run() *MainWorker {
	//runWorkers()
	b.start = time.Now()
	reporter := b.runWorkers()
	b.End = time.Now()
	return reporter
}

func (b *Requester) worker(wID string) *SubWorker {
	reporter := newReporter(wID)
	ctx, cancel := context.WithTimeout(context.Background(), b.config.timeout*time.Second)
	defer cancel()
	sw := b.stopwatch
	tlsCredentials, err := LoadTLSCredentials(b.config.skipVerify, b.config.cert)
	opts := GrpcOption(b.config, tlsCredentials)
	reporter.Start = sw.Start()
	dur := sw.Track()
	reporter.StartTime = dur

	conn, err := grpc.DialContext(ctx,
		fmt.Sprintf("%s:%d", b.config.host, b.config.port),
		opts...,
	)
	if err != nil {
		return reporter.ErrorHandler("did not connect", "DialContext", err)
	}

	if b.config.block {
		go checkConnectivityStatusChan(ctx, conn, connectivity.Idle)
	}
	err = conn.Close()
	if err != nil {
		return reporter.ErrorHandler("connection close error", "CloseError", err)
	}
	reporter.EndTime, reporter.End = sw.Stop()
	return reporter
	//go checkConnectivityStatusChan(ctx, conn, connectivity.Idle)
	//conn, err := driver.Dial(address, driver.WithInsecure(), driver.WithBlock())
	// c := pb.NewGreeterClient(conn)

	// r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: defaultMessage})
	// if err != nil {
	// return errors2.New("could not greet:", "SayHello", err)
	// }

	//report.PrinterReport(err)
}

func (b *Requester) runWorkers() *MainWorker {
	reporter := newLoadReporter()
	var wg sync.WaitGroup
	wg.Add(b.config.TotalRequest)

	fmt.Printf("%s:%d", b.config.host, b.config.port)

	for i := 0; i < b.config.TotalRequest; i++ {
		go func(b *Requester) {
			n := 0
			wc := 0
			wID := "g" + strconv.Itoa(wc) + "c" + strconv.Itoa(n)
			worker := b.worker(wID)
			reporter.addSimpleReports(*worker)
			wc++
			n++
			defer wg.Done()

		}(b)
		time.Sleep(time.Duration(b.config.rps))

	}
	wg.Wait()
	reporter.EndTime = time.Now()
	reporter.Rps = b.config.rps
	reporter.Finish()

	return reporter
}

func (r *MainWorker) Finish() {
	// Slowest / Fastest
	var s, f, totalLatenciesSec time.Duration
	var latencies []float64
	//okLats := make([]float64, 0)
	var errReport []SubWorker

	for _, ar := range r.Workers {
		sec := ar.EndTime

		latencies = append(latencies, float64(sec))
		totalLatenciesSec += sec

		if ar.Error != nil {
			errReport = append(errReport, ar)
		}

	}

	r.Slowest = s
	r.Fastest = f

	//average := totalLatenciesSec / float64(r.TotalCount)
	//r.Average = sum / time.Duration(len(ltr.SimpleTestReports))
	//
	//r.Histogram = calHistogram(latencies, float64(s), float64(f))
	//
	////
	//// StatusDistribute
	//var statusDistribute []Status
	//statusDistribute = append(statusDistribute, Status{
	//	Mark:    "OK",
	//	Count:   len(altr.AcceptReports) - len(errReport),
	//	Message: "reponse",
	//})
	//for _, er := range errReport {
	//	exist := false
	//	for idx, s := range statusDistribute {
	//		if s.Message == er.ErrorMessage {
	//			statusDistribute[idx].Count++
	//			exist = true
	//			break
	//		}
	//	}
	//	if !exist {
	//		statusDistribute = append(statusDistribute, Status{
	//			Mark:    er.ErrorCode,
	//			Count:   1,
	//			Message: er.ErrorMessage,
	//		})
	//	}
	//}
	//altr.StatusDistribute = statusDistribute
}

//
//
//	if len(r.details) > 0 {
//		average := r.totalLatenciesSec / float64(r.totalCount)
//		rep.Average = time.Duration(average * float64(time.Second))
//
//		rep.Rps = float64(r.totalCount) / total.Seconds()
//
//		okLats := make([]float64, 0)
//		//for _, d := range r.details {
//		//	if d.Error == "" || rep.Options.CountErrors {
//		//		okLats = append(okLats, d.Latency.Seconds())
//		//	}
//		//}
//		sort.Float64s(okLats)
//		if len(okLats) > 0 {
//			var fastestNum, slowestNum float64
//			fastestNum = okLats[0]
//			slowestNum = okLats[len(okLats)-1]
//
//			rep.Fastest = time.Duration(fastestNum * float64(time.Second))
//			rep.Slowest = time.Duration(slowestNum * float64(time.Second))
//			rep.Histogram = histogram(okLats, slowestNum, fastestNum)
//		}
//
//		rep.Details = r.details
//	}
//
//	return rep
//}
