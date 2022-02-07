package runner

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

func Run(c *RunConfig) (*Reporter, error) {
	report := &Reporter{}
	r := NewRequester(c)
	mainWorker, err := r.Run()
	report.MainWorker = mainWorker
	if err != nil {
		return nil, err
	}
	return report, nil

}

// Check the change of connectivity state.
// Wait for state change, then notify and recurse
func checkConnectivityStatusChan(ctx context.Context, conn *grpc.ClientConn, sourceState connectivity.State) {
	ch := make(chan bool, 1)
	ch <- conn.WaitForStateChange(ctx, sourceState)

	select {
	case <-ctx.Done():
	case <-ch:
		curState := conn.GetState()
		close(ch)
		go checkConnectivityStatusChan(ctx, conn, curState)
	}
}
