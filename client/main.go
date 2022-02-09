package main

import (
	"Src1/client/printer"
	"Src1/client/runner"
	"fmt"
	"github.com/alecthomas/kingpin"
	"os"
)

const (
	defaultHost    = "localhost"
	defaultPort    = "50051"
	defaultMessage = "hello man"
	//defaultCertPem = "C:\\Users\\K\\gopath\\src\\Src1\\cert\\ca-cert.pem"
)

var (
	skipVerify = kingpin.Flag("skipTLS", "Skip TLS client verification of the server's certificate chain and host name.").
			Default("false").Short('s').Bool()

	host = kingpin.Flag("host", "Hostname").Default(defaultHost).Short('h').String()

	port = kingpin.Flag("port", "Port number").Default(defaultPort).Short('p').Int()

	cert = kingpin.Flag("cert", "TLS cert.pem").PlaceHolder(" ").String()

	call = kingpin.Flag("call", "Procedure call").String()

	protoPath = kingpin.Flag("protoPath", "proto file path").String()

	totalRequest = kingpin.Flag("totalCount", "total count").Default("1").Short('n').Uint()

	blockMode = kingpin.Flag("blockMode", "Dial BlockMode").Default("true").Short('b').Bool()

	timeOut = kingpin.Flag("timeOut", "Time Out option").Default("20s").Short('t').Duration()

	rps = kingpin.Flag("rps", "Requests per second (RPS) rate limit for constant load schedule. Default is no rate limit.").
		Default("1").Short('r').Uint()
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
	var cfg runner.Config
	//args := os.Args[1:]
	err := createConfigFromArgs(&cfg)
	if err != nil {
		kingpin.FatalIfError(err, "")
		os.Exit(1)
	}
	options := []runner.Option{runner.WithConfig(&cfg)}
	c, _ := runner.NewConfig(cfg.Host, cfg.Port, options...)
	report, err := runner.Run(c)

	if err != nil {
		handleError(err)
	}
	p := printer.ReportPrinter{
		Report: report,
		Config: c,
	}

	p.Print()
}

func createConfigFromArgs(config *runner.Config) error {
	if config == nil {
		return nil
	}

	config.Host = *host
	config.Port = *port
	config.SkipTLSVerify = *skipVerify
	config.Cert = *cert
	config.TotalRequest = *totalRequest
	config.TimeOut = runner.Duration(*timeOut)
	config.BlockMode = *blockMode
	config.RPS = *rps
	config.Call = *call
	config.ProtoPath = *protoPath

	//config.KeyPem = *keyPem
	return nil
}
func handleError(err error) {
	if err != nil {
		if errString := err.Error(); errString != "" {
			fmt.Fprintln(os.Stderr, errString)
		}
		os.Exit(1)
	}
}
