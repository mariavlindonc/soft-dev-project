package main

import (
	"os"

	"backend/logger"
)

func main() {
	template, key, err := generateCertAndKey()
	if err != nil {
		panic(err)
	}

	certDER, err := createCertDER(template, key)
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll("certs", 0755); err != nil {
		panic(err)
	}

	if err := writeCertFiles("certs", certDER, key); err != nil {
		panic(err)
	}

	logger.Info("certificates generated in certs/")
}
