package tls

import (
	"Src1/client/runner"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

func LoadTLSCredentials(c *runner.Config) (credentials.TransportCredentials, error) {
	//Load certificate of the CA who signed server's certificate7
	if c.SkipTLSVerify {
		pemServerCA, err := ioutil.ReadFile(c.CertPem)
		if err != nil {
			return nil, err
		}

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(pemServerCA) {
			return nil, fmt.Errorf("failed to add server CA's certificate")
		}
		// Create the credentials and return it
		config := &tls.Config{
			RootCAs: certPool,
		}
		return credentials.NewTLS(config), nil
	}
	return nil, nil
}
