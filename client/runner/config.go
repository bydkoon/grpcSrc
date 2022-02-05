package runner

import "time"

type Duration time.Duration

type Config struct {
	Host          string `json:"host" toml:"host" yaml:"host"`
	Port          int    `json:"port" toml:"port" yaml:"port"`
	TotalRequest  uint   `json:"TotalRequest" toml:"TotalRequest" yaml:"TotalRequest"`
	SkipTLSVerify bool   `json:"skipTLS" toml:"skipTLS" yaml:"skipTLS"`
	Cert          string `json:"cert" toml:"cert" yaml:"cert"`
	BlockMode     bool   `json:"blockMode" toml:"blockMode" yaml:"blockMode"`
	//TimeOut    			   int    `json:"timeOut" toml:"timeOut" yaml:"timeOut"`
	RPS     uint     `json:"rps" toml:"rps" yaml:"rps"`
	TimeOut Duration `json:"timeout" toml:"timeout" yaml:"timeout" default:"1s"`
}
