package utility

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/ywengineer/smart-kit/pkg/logk"
	"go.uber.org/zap"
	"os"
)

func CreateTlsConfiguration(certFile, keyFile, caFile *string, verifySsl bool) (t *tls.Config) {
	if StringWithNoSpace(*certFile) && StringWithNoSpace(*keyFile) && StringWithNoSpace(*caFile) {
		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			logk.Fatal("failed to LoadX509KeyPair(cert, key). %v", zap.Error(err))
		}
		caCert, err := os.ReadFile(*caFile)
		if err != nil {
			logk.Fatal("failed to read ca file. %v", zap.Error(err))
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: verifySsl,
		}
	}
	// will be nil by default if nothing is provided
	return t
}
