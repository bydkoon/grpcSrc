package driver

import (
	"__TestSrc/gRpc_Study/Src1/client/cmd"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func GrpcOption(opts []grpc.DialOption, cfg *cmd.Config, tlsCredentials credentials.TransportCredentials) []grpc.DialOption {
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
