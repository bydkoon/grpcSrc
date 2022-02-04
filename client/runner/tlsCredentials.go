package runner

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

func LoadTLSCredentials(skipVerify bool, clientkeyFile string) (credentials.TransportCredentials, error) {
	//Load certificate of the CA who signed server's certificate7
	var tlsConf tls.Config

	if skipVerify {
		tlsConf.InsecureSkipVerify = true
		pemServerCA, err := ioutil.ReadFile(clientkeyFile)
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
