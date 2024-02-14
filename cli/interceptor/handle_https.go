package interceptor

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/yieldray/middleman/cli/utils"
	// "fmt"
	// "net/http/httputil"
)

const httpsServerAddr = "127.0.0.1:9443"

type Handler func(http.ResponseWriter, *http.Request)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}

func confHttpsServer(caKeyPath string, caCertPath string, httpClient http.Client) (*http.Server, error) {
	l.Debug("Now launching https server, keyPath=%s certPath=%s", caKeyPath, caCertPath)

	caKey, caCert, err := utils.LoadKeyCert(caKeyPath, caCertPath)

	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr: httpsServerAddr,
		TLSConfig: &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				host := strings.ToLower(info.ServerName)
				l.Info("[Generate Certificate] %s", host)
				return utils.GenerateCertificate(host, caKey, caCert)
			},
		},
		Handler: Handler(func(w http.ResponseWriter, r *http.Request) {
			u, _ := url.Parse("https://" + r.Host + r.URL.String())

			res, err := httpClient.Do(&http.Request{
				Method:        r.Method,
				URL:           u,
				Header:        r.Header,
				Body:          r.Body,
				ContentLength: r.ContentLength,
				Close:         r.Close,
				Trailer:       r.Trailer,
			})

			// buf, _ := httputil.DumpRequest(r, false)
			// fmt.Printf("%s\n\n", buf)

			if err != nil {
				l.Error("%s", err)
				return
			}

			utils.CopyHeader(w.Header(), res.Header)
			w.WriteHeader(res.StatusCode)
			io.Copy(w, res.Body)
		}),
	}

	return server, nil

}

func handleProxyHTTPS(conn net.Conn, reqOriginal *http.Request) {
	io.WriteString(conn, "HTTP/1.1 200 Connection Established\r\n\r\n")

	localConn, _ := net.Dial("tcp", httpsServerAddr)

	go utils.CopyCloser(localConn, conn)
	go utils.CopyCloser(conn, localConn)
}
