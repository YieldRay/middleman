package interceptor

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func NewRoundTripper(roundTrip func(req *http.Request) (*http.Response, error)) http.RoundTripper {
	return transportRoundTrip{
		RoundTripImpl: roundTrip,
	}
}

// this type is just a wrapper for the interface: RoundTrip
type transportRoundTrip struct {
	RoundTripImpl func(req *http.Request) (*http.Response, error)
}

func (trt transportRoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	return trt.RoundTripImpl(req)
}

// just wrap string, making sure string is correct proxy server address
type ProxyServer string

func (ps ProxyServer) Proxy(targetUrl string) string {
	return string(ps) + targetUrl
}
func (ps ProxyServer) ProxyURL(targetUrl string) *url.URL {
	u, _ := url.Parse(ps.Proxy(targetUrl))
	return u
}

// just a wrapper of string
func NewProxyServer(urlLike string) (ProxyServer, error) {
	if !strings.HasPrefix(urlLike, "http://") && !strings.HasPrefix(urlLike, "https://") {
		return "", fmt.Errorf("malformed url")
	}
	if !strings.HasSuffix(urlLike, "/") {
		return ProxyServer(urlLike) + "/", nil
	}
	return ProxyServer(urlLike), nil
}
