package main

import (
	"context"
	"crypto/tls"
	"flag"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"log"
	"net"
	"time"

	pb "Src1/proto"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
	name string
}

var kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second,
	PermitWithoutStream: true,
}

var kasp = keepalive.ServerParameters{

	MaxConnectionIdle:     15 * time.Second,
	MaxConnectionAge:      30 * time.Second,
	MaxConnectionAgeGrace: 5 * time.Second,
	Time:                  5 * time.Second,
	Timeout:               1 * time.Second,
}
var (
	system = ""
	sleep  = flag.Duration("sleep", time.Second*5, "duration between changes in health")
)

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {

	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: in.GetName()}, nil
}

const (
	defaultServerCertPem = "C:\\Users\\K\\gopath\\src\\__TestSrc\\gRpc_Study\\Src1\\cert\\server-cert.pem"
	defaultServerKeyFile = "C:\\Users\\K\\gopath\\src\\__TestSrc\\gRpc_Study\\Src1\\cert\\server-key.pem"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	//Load server's certificate and private key
	serverCrt, err := tls.LoadX509KeyPair(defaultServerCertPem, defaultServerKeyFile)
	if err != nil {
		return nil, err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCrt},
		ClientAuth:   tls.NoClientCert,
	}
	return credentials.NewTLS(config), nil
}

func orderUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("====== [Server Interceptor] %s", info.FullMethod)

	m, err := handler(ctx, req)

	return m, err
}

func main() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials : ", err)

	}

	opts = append(opts) //grpc.InitialConnWindowSize(0),
	//grpc.UnaryInterceptor(orderUnaryServerInterceptor),
	grpc.Creds(tlsCredentials)
	s := grpc.NewServer(opts...)

	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	//reflection.Register(s) // grpcurl 명령을 사용하게 하기 위해
	//healthcheck := health.NewServer()
	//healthpb.RegisterHealthServer(s, healthcheck)

}
