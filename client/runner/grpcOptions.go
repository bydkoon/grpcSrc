package driver

import (
	"Src1/client/runner"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func GrpcOption(opts []grpc.DialOption, cfg *runner.Config, tlsCredentials credentials.TransportCredentials) []grpc.DialOption {
	//opts = append(opts) grpc.WithReturnConnectionError(),
	//grpc.FailOnNonTempDialError(true),
	//grpc.WithBlock(),

	if cfg.SkipTLSVerify {
		opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	if cfg.BlockMode {
		opts = append(opts, grpc.WithBlock())
	}

	return opts
}
