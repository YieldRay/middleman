package utils

import (
	"io"
	"net/http"
)

func CopyCloser(dst io.WriteCloser, src io.ReadCloser) {
	if src != nil {
		defer src.Close()
	}
	if dst != nil {
		defer dst.Close()
	}
	if src != nil && dst != nil {
		io.Copy(dst, src)
	}
}

func CopyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
