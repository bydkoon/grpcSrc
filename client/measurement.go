package main

import (
	"Src1/client/driver"
	errors2 "Src1/client/errors"
	"Src1/client/printer"
	"Src1/client/tls"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"Src1/client/runner"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

func mesasurement() {
	var report printer.Report
	report.MakeReport(cfg)
	report.TrackReport.Init()
	//sw := utils.NewStopWatch()
	//sw.Start()
	/*	pid := getGID()*/

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

func createConfigFromArgs(config *runner.Config) error {
	if config == nil {
		return errors2.New("config cannot be nil", "err")
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

func worker(opts []grpc.DialOption, cfg runner.Config, report printer.Report) *errors2.Errors {
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
	// c := pb.NewGreeterClient(conn)

	// r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: defaultMessage})
	// if err != nil {
	// return errors2.New("could not greet:", "SayHello", err)
	// }

	report.TrackReport.TrackCheck("sendEnd", time.Since(start_time))
	// log.Printf("Greeting: %s", r.GetMessage())
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
