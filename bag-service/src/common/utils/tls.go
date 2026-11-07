package utils

import (
	"bag-service/src/common/log"
	"crypto/tls"
	"crypto/x509"
)

func CreateClientTlsConf(caCert []byte, clientCert tls.Certificate) *tls.Config {
	certroot := x509.NewCertPool()
	if !certroot.AppendCertsFromPEM(caCert) {
		log.Logger.Fatal("Failed to append CA certificate")
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            certroot,
		InsecureSkipVerify: false,
	}
}
