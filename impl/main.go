package impl

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"software.sslmate.com/src/go-pkcs12"
)

func CreateBundle(certPEM string, keyPEM string) (*string, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("could not decode public certificate")
	}

	publicKey, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.New("could not parse public certificate")
	}

	keyBlock, _ := pem.Decode([]byte(keyPEM))
	if keyBlock == nil || keyBlock.Type != "CERTIFICATE" {
		return nil, errors.New("could not decode private key")
	}

	var privateKey interface{}
	privateKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		privateKey, err = x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	}
	if err != nil {
		return nil, errors.New("could not parse private key")
	}

	result, err := pkcs12.Encode(
		rand.Reader,
		privateKey,
		publicKey,
		[]*x509.Certificate{},
		pkcs12.DefaultPassword,
	)

	if err != nil {
		return nil, errors.New("failed to build PKCS12 bundle")
	}

	str := base64.StdEncoding.EncodeToString(result)
	return &str, nil
}
