package runner

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func GrpcOption(cfg *RunConfig, tlsCredentials credentials.TransportCredentials) []grpc.DialOption {
	//opts = append(opts) grpc.WithReturnConnectionError(),
	//grpc.FailOnNonTempDialError(true),
	//grpc.WithBlock(),
	var opts []grpc.DialOption
	if cfg.SkipVerify {
		opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	if cfg.Block {
		opts = append(opts, grpc.WithBlock())
	}

	return opts
}
