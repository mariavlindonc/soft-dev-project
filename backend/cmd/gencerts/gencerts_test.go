package main

import (
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCertAndKey(t *testing.T) {
	template, key, err := generateCertAndKey()
	require.NoError(t, err)
	require.NotNil(t, template)
	require.NotNil(t, key)

	assert.Equal(t, "localhost", template.Subject.CommonName)
	assert.Contains(t, template.DNSNames, "localhost")
	assert.Contains(t, template.DNSNames, "backend")
	assert.Equal(t, []string{"AR"}, template.Subject.Country)
	assert.Equal(t, []string{"Ceibo"}, template.Subject.Organization)
	assert.True(t, template.BasicConstraintsValid)
	assert.True(t, template.NotAfter.After(template.NotBefore))
	assert.Equal(t, 2048, key.N.BitLen())
}

func TestCreateCertDER(t *testing.T) {
	template, key, err := generateCertAndKey()
	require.NoError(t, err)

	certDER, err := createCertDER(template, key)
	require.NoError(t, err)
	require.NotEmpty(t, certDER)

	block, _ := pem.Decode(certDER)
	assert.Nil(t, block, "DER bytes should not be PEM encoded")

	cert, err := x509.ParseCertificate(certDER)
	require.NoError(t, err)
	assert.Equal(t, "localhost", cert.Subject.CommonName)
	assert.Contains(t, cert.DNSNames, "localhost")
	assert.Contains(t, cert.DNSNames, "backend")
}

func TestWriteCertFiles(t *testing.T) {
	dir := t.TempDir()

	template, key, err := generateCertAndKey()
	require.NoError(t, err)

	certDER, err := createCertDER(template, key)
	require.NoError(t, err)

	err = writeCertFiles(dir, certDER, key)
	require.NoError(t, err)

	t.Run("cert file exists and is valid PEM", func(t *testing.T) {
		certPath := filepath.Join(dir, "server.crt")
		data, err := os.ReadFile(certPath)
		require.NoError(t, err)

		block, _ := pem.Decode(data)
		require.NotNil(t, block)
		assert.Equal(t, "CERTIFICATE", block.Type)

		cert, err := x509.ParseCertificate(block.Bytes)
		require.NoError(t, err)
		assert.Equal(t, "localhost", cert.Subject.CommonName)
	})

	t.Run("key file exists and is valid PEM", func(t *testing.T) {
		keyPath := filepath.Join(dir, "server.key")
		data, err := os.ReadFile(keyPath)
		require.NoError(t, err)

		block, _ := pem.Decode(data)
		require.NotNil(t, block)
		assert.Equal(t, "RSA PRIVATE KEY", block.Type)

		parsedKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		require.NoError(t, err)
		assert.Equal(t, 2048, parsedKey.N.BitLen())
	})
}
