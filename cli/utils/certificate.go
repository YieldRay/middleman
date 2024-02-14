package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

type keyCert struct {
	key  any
	cert *x509.Certificate
}

var cacheKeyCert = make(map[string]map[string]keyCert, 1)

// load key and cert by given path
func LoadKeyCert(caKeyPath, caCertPath string) (caKey any, caCert *x509.Certificate, err error) {

	if keyCert, ok := cacheKeyCert[caKeyPath][caCertPath]; ok {
		// only load once
		return keyCert.key, keyCert.cert, nil
	}

	caKeyPEM, err := os.ReadFile(caKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("%v\n    ($ openssl geeky -algorithm RSA -out ca.key)", err)
	}

	caKeyBlock, _ := pem.Decode(caKeyPEM)

	caKey, err = x509.ParsePKCS8PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	caCertPEM, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, nil, fmt.Errorf("%v\n    ($ openssl req -x509 -new -key ca.key -out ca.crt -days 3650)", err)
	}

	caCertBlock, _ := pem.Decode(caCertPEM)

	caCert, err = x509.ParseCertificate(caCertBlock.Bytes)

	if err != nil {
		return nil, nil, err
	}

	return
}

// gen new certificate for host, derived by ca
func GenerateCertificate(host string, caKey any, caCert *x509.Certificate) (*tls.Certificate, error) {
	crtPEM, keyPEM, err := LoadCertificate(host, caKey, caCert)
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair(crtPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	return &cert, nil
}

// load certificate from local path
func LoadCertificate(host string, caKey any, caCert *x509.Certificate) (certPEM, keyPEM []byte, err error) {
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	notBefore := caCert.NotBefore
	notAfter := notBefore.Add(365 * 24 * time.Hour) // Valid for 1 year

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Middleman Certificates"},
			CommonName:   host,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{host},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, caCert, &priv.PublicKey, caKey)
	if err != nil {
		return nil, nil, err
	}
	certPEMBlock := &pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	keyDER, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}
	keyPEMBlock := &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}

	certPEM = pem.EncodeToMemory(certPEMBlock)
	keyPEM = pem.EncodeToMemory(keyPEMBlock)

	return certPEM, keyPEM, nil
}
