package impl

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"path"
	"software.sslmate.com/src/go-pkcs12"
)

func CreateBundle(certPEM string, keyPEM string, outputPath string, name string) error {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil || block.Type != "CERTIFICATE" {
		return errors.New("could not decode public certificate")
	}

	publicKey, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return errors.New("could not parse public certificate")
	}

	keyBlock, _ := pem.Decode([]byte(keyPEM))
	if keyBlock == nil || keyBlock.Type != "RSA PRIVATE KEY" {
		return errors.New("could not decode private key")
	}

	var privateKey interface{}
	privateKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		privateKey, err = x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	}
	if err != nil {
		return errors.New("could not parse private key")
	}

	result, err := pkcs12.Encode(
		rand.Reader,
		privateKey,
		publicKey,
		[]*x509.Certificate{},
		pkcs12.DefaultPassword,
	)

	if err != nil {
		return errors.New("failed to build PKCS12 bundle")
	}

	var outFile = path.Join(outputPath, name+".p12")
	err = ioutil.WriteFile(outFile, result, 0644)
	if err != nil {
		return errors.New("failed to save bundle to file")
	}

	return nil
}
