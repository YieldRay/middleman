package interceptor

import (
	"net"
	"net/http"
)

func handleProxyHTTP(conn net.Conn, reqOriginal *http.Request, httpClient http.Client) {
	var r = reqOriginal

	res, err := httpClient.Do(&http.Request{
		Method:        r.Method,
		URL:           r.URL,
		Header:        r.Header,
		Body:          r.Body,
		ContentLength: r.ContentLength,
		Close:         r.Close,
		Trailer:       r.Trailer,
	})

	if err != nil {
		l.Error("%s", err)
		return
	}

	res.Write(conn)
}
