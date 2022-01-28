package main

import (
	"Src1/client/printer"
	"os"

	"Src1/client/runner"

	"github.com/alecthomas/kingpin"
	"google.golang.org/grpc"
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
	var opts []grpc.DialOption
	var cfg runner.Config

	args := os.Args[1:]

	if len(args) > 1 {
		//var cmdCfg cmd.Config
		err := createConfigFromArgs(&cfg)
		if err != nil {
			kingpin.FatalIfError(err, "")
			os.Exit(1)
		}

	}

	p := printer.ReportPrinter{
		Report: report,
		Out:    output,
	}

}
