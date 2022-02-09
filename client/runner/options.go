package runner

import (
	"google.golang.org/grpc/credentials"
	"strings"
	"time"
)

type RunConfig struct {
	// call settings
	Host              string
	Port              int
	proto             string
	enableCompression bool

	Creds        credentials.TransportCredentials
	Timeout      time.Duration
	TotalRequest int

	// security settings
	Cert       string
	SkipVerify bool
	Insecure   bool
	Block      bool

	// concurrency
	C int

	// load
	Rps              uint
	LoadStart        uint
	LoadEnd          uint
	LoadStep         int
	LoadSchedule     string
	LoadDuration     time.Duration
	LoadStepDuration time.Duration
}

type Option func(*RunConfig) error

func NewConfig(host string, port int, options ...Option) (*RunConfig, error) {

	// init with defaults
	c := &RunConfig{}
	// apply options
	for _, option := range options {
		err := option(c)

		if err != nil {
			return nil, err
		}
	}

	// host and call may have been applied via options
	// only override if not present
	if c.Host == "" {
		c.Host = strings.TrimSpace(host)
	}

	c.Port = port

	if c.Cert != "" {
		creds, err := LoadTLSCredentials(
			c.SkipVerify,
			c.Cert,
		)
		if err != nil {
			return nil, err
		}
		c.Creds = creds
	}

	return c, nil
}

func WithConfig(cfg *Config) Option {
	return func(o *RunConfig) error {

		// init / fix up durations

		for _, option := range fromConfig(cfg) {
			if err := option(o); err != nil {
				return err
			}
		}
		return nil
	}
}

func fromConfig(cfg *Config) []Option {
	options := make([]Option, 0, 7)
	options = append(options,
		WithTimeout(time.Duration(cfg.TimeOut)),
		WithTotalRequest(cfg.TotalRequest),
		WithBlock(cfg.BlockMode),
		WithInsecure(cfg.SkipTLSVerify),
		WithHost(cfg.Host),
		WithPort(cfg.Port),
		WithCert(cfg.Cert),
		WithRps(cfg.RPS),
	)
	return options
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *RunConfig) error {
		o.Timeout = timeout

		return nil
	}
}
func WithHost(host string) Option {
	return func(o *RunConfig) error {
		o.Host = host
		return nil
	}
}

func WithPort(port int) Option {
	return func(o *RunConfig) error {
		o.Port = port
		return nil
	}
}

func WithCert(cert string) Option {
	return func(o *RunConfig) error {
		o.Cert = cert
		return nil
	}
}

func WithRps(rps uint) Option {
	return func(o *RunConfig) error {
		o.Rps = rps
		return nil
	}
}

func WithTotalRequest(n uint) Option {
	return func(o *RunConfig) error {
		o.TotalRequest = int(n)

		return nil
	}
}
func WithBlock(block bool) Option {
	return func(o *RunConfig) error {
		o.Block = block
		return nil
	}
}
func WithInsecure(v bool) Option {
	return func(o *RunConfig) error {
		o.Insecure = v
		return nil
	}

}
