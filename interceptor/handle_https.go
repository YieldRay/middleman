package interceptor

import (
	"io"
	"net"
	"net/http"
	"net/url"
	// "fmt"
	// "net/http/httputil"
)

func startHttpsServer(caKeyPath string, caCrtPath string, httpClient http.Client) *http.Server {
	// the server may choose a random port to listen at
	server := createHttpsServer(caKeyPath, caCrtPath, "127.0.0.1:9443", func(w http.ResponseWriter, r *http.Request) {

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

		copyHeader(w.Header(), res.Header)
		w.WriteHeader(res.StatusCode)
		io.Copy(w, res.Body)
	})

	return server
}

func handleProxyHTTPS(conn net.Conn, reqOriginal *http.Request, httpsServerAddr string) {
	io.WriteString(conn, "HTTP/1.1 200 Connection Established\r\n\r\n")

	localConn, _ := net.Dial("tcp", httpsServerAddr)

	go copyCloser(localConn, conn)
	go copyCloser(conn, localConn)
}
