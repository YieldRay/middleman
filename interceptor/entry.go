package interceptor

import (
	"bufio"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mborders/logmatic"
	"github.com/yieldray/middleman/cmd/flags"
)

var l = logmatic.NewLogger()

func init() {
	l.SetLevel(logmatic.LogLevel(flags.LogLevel))
}

func Entry(httpProxyAddr string, httpClient http.Client, caKeyPath string, caCrtPath string) {
	ln, err := net.Listen("tcp", httpProxyAddr)
	if err != nil {
		l.Fatal("%s", err)
	}
	l.Info("代理服务器运行在 %s", ln.Addr().String())

	// set system proxy
	if err := setProxySettings(httpProxyAddr); err == nil {
		l.Info("已设置系统代理为 %s", httpProxyAddr)

		var turnOffSystemProxy = func() {
			if err := disableProxySettings(); err != nil {
				l.Warn("%s", err)
			} else {
				l.Info("已关闭系统代理")
			}
		}

		// when error
		defer func() {
			if err := recover(); err != nil {
				turnOffSystemProxy()
				l.Fatal("%s", err)
			}
		}()

		// when Ctrl-C
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigChan
			turnOffSystemProxy()
			os.Exit(0)
		}()
	} else {
		l.Warn("%s", err)
	}

	// run https server in the background
	go func() {
		var server = startHttpsServer(caKeyPath, caCrtPath, httpClient)
		server.ListenAndServeTLS("", "")
		defer server.Close()
	}()

	// listen tcp server
	for {
		conn, err := ln.Accept()
		if err != nil {
			l.Error("%#v", err)
			continue
		}

		// handle it
		go func(conn net.Conn) {
			req, err := http.ReadRequest(bufio.NewReader(conn))

			if err != nil {
				l.Error("%s", err)
				conn.Close()
				return
			}
			if req.Method == "CONNECT" {
				handleProxyHTTPS(conn, req, "127.0.0.1:9443")
			} else {
				handleProxyHTTP(conn, req, httpClient)
			}
		}(conn)
	}
}
