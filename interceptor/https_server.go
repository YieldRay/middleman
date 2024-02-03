package interceptor

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"
)

// Return http.Server, but need to serve it:
// server.ListenAndServeTLS("", "")
func createHttpsServer(caKeyPath, caCertPath, addr string, handler func(http.ResponseWriter, *http.Request)) *http.Server {

	caKeyPEM, err := os.ReadFile(caKeyPath)
	if err != nil {
		l.Warn("要生成 CA 密钥，运行 $ openssl genpkey -algorithm RSA -out ca.key")
		panic(err)
	}

	caKeyBlock, _ := pem.Decode(caKeyPEM)

	caKey, err := x509.ParsePKCS8PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		l.Fatal("%s", err)
	}

	caCertPEM, err := os.ReadFile(caCertPath)
	if err != nil {
		l.Warn("要生成 CA 证书，运行 $ openssl req -x509 -new -key ca.key -out ca.crt -days 3650")
		panic(err)
	}

	caCertBlock, _ := pem.Decode(caCertPEM)

	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr: addr,
		TLSConfig: &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				host := strings.ToLower(info.ServerName)
				l.Info("[Generate Certificate] %s", host)

				crtPEM, keyPEM, err := generateCertificate(caKey, caCert, host)
				if err != nil {
					l.Error("Error generating certificate for %s: %v", host, err)
					return nil, err
				}

				cert, err := tls.X509KeyPair(crtPEM, keyPEM)
				if err != nil {
					l.Error("Error loading certificate for %s: %v", host, err)
					return nil, err
				}

				return &cert, nil
			},
		},
	}

	http.HandleFunc("/", handler)
	return server
}

func generateCertificate(caKey any, caCert *x509.Certificate, host string) (certPEM, keyPEM []byte, err error) {
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
