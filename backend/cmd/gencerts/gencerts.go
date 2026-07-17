package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

func generateCertAndKey() (*x509.Certificate, *rsa.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Country:      []string{"AR"},
			Province:     []string{"Cordoba"},
			Locality:     []string{"Cordoba"},
			Organization: []string{"Ceibo"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost", "backend"},
	}
	return &template, key, nil
}

func createCertDER(template *x509.Certificate, key *rsa.PrivateKey) ([]byte, error) {
	return x509.CreateCertificate(rand.Reader, template, template, key.Public(), key)
}

func writeCertFiles(dir string, certDER []byte, key *rsa.PrivateKey) error {
	certPath := filepath.Join(dir, "server.crt")
	keyPath := filepath.Join(dir, "server.key")

	certOut, err := os.Create(certPath)
	if err != nil {
		return err
	}
	defer certOut.Close()
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	keyOut, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer keyOut.Close()
	pem.Encode(keyOut, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	return nil
}
