package runner

import (
	"google.golang.org/grpc/credentials"
	"strings"
	"time"
)

type RunConfig struct {
	// call settings
	host              string
	port              int
	proto             string
	enableCompression bool

	creds        credentials.TransportCredentials
	timeout      time.Duration
	TotalRequest int
	cWorker      int

	// security settings
	cert       string
	skipVerify bool
	insecure   bool
	block      bool

	// concurrency
	c int

	// load
	rps              int
	loadStart        uint
	loadEnd          uint
	loadStep         int
	loadSchedule     string
	loadDuration     time.Duration
	loadStepDuration time.Duration
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
	if c.host == "" {
		c.host = strings.TrimSpace(host)
	}

	c.port = port
	creds, err := LoadTLSCredentials(
		c.skipVerify,
		c.cert,
	)

	if err != nil {
		return nil, err
	}

	c.creds = creds

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
	options := make([]Option, 0, 4)
	options = append(options,
		WithTimeout(time.Duration(cfg.TimeOut)),
		WithTotalRequest(cfg.TotalRequest),
		WithBlock(cfg.BlockMode),
		WithInsecure(cfg.SkipTLSVerify),
		WithHost(cfg.Host),
		//WithPort(cfg.Port)
	)
	return options
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *RunConfig) error {
		o.timeout = timeout

		return nil
	}
}
func WithHost(v string) Option {
	return func(o *RunConfig) error {
		o.host = v
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
		o.block = block
		return nil
	}
}
func WithInsecure(v bool) Option {
	return func(o *RunConfig) error {
		o.insecure = v
		return nil
	}

}
