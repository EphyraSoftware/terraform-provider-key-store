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

func CreateBundle(certPEM string, keyPEM string, caCertsPEM []string, outputPath string, name string) error {
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

	caCerts := make([]*x509.Certificate, len(caCertsPEM))
	for i, cert := range caCertsPEM {
		block, _ := pem.Decode([]byte(cert))
		if block == nil || block.Type != "CERTIFICATE" {
			return errors.New("could not decode CA certificate")
		}

		caCert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return errors.New("could not parse CA certificate")
		}

		caCerts[i] = caCert
	}

	result, err := pkcs12.Encode(
		rand.Reader,
		privateKey,
		publicKey,
		caCerts,
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

func SliceOfString(slice []interface{}) []string {
	result := make([]string, len(slice), len(slice))
	for i, s := range slice {
		result[i] = s.(string)
	}
	return result
}
