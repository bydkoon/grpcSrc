package runner

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"time"
)

func Run(c *RunConfig) *Reporter {
	report := &Reporter{}
	reqr := NewRequester(c)

	report.MainWorker = reqr.Run()

	return report

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
