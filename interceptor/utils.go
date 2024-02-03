package interceptor

import (
	"io"
	"net/http"

	"golang.org/x/sys/windows/registry"
)

func copyCloser(dst io.WriteCloser, src io.ReadCloser) {
	defer dst.Close()
	defer src.Close()
	io.Copy(dst, src)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func setProxySettings(proxyServer string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	if err = k.SetStringValue("AutoConfigURL", ""); err != nil {
		return err
	}
	if err = k.SetDWordValue("ProxyEnable", uint32(1)); err != nil {
		return err
	}
	if err = k.SetStringValue("ProxyServer", proxyServer); err != nil {
		return err
	}
	if _, _, err = k.GetIntegerValue("ProxyOverride"); err == nil {
		if err = k.DeleteValue("ProxyOverride"); err != nil {
			return err
		}
	}
	return nil
}

func disableProxySettings() error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()

	return k.SetDWordValue("ProxyEnable", uint32(0))
}
