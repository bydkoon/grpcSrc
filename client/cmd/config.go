package cmd

type Config struct {
	Host          string `json:"host" toml:"host" yaml:"host"`
	Port          int    `json:"port" toml:"port" yaml:"port"`
	TotalCount    int    `json:"TotalCount" toml:"TotalCount" yaml:"TotalCount"`
	SkipTLSVerify bool   `json:"skipTLS" toml:"skipTLS" yaml:"skipTLS"`
	CertPem       string `json:"certPem" toml:"certPem" yaml:"certPem"`
	KeyPem        string `json:"keyPem" toml:"keyPem" yaml:"keyPem"`
	BlockMode     bool   `json:"blockMode" toml:"blockMode" yaml:"blockMode"`
	TimeOut       int    `json:"timeOut" toml:"timeOut" yaml:"timeOut"`
}
