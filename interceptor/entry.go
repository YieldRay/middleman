package interceptor

import (
	"bufio"
	"github.com/yieldray/middleman/utils"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// httpProxyAddr: where the http proxy server listen to
//
// httpClient: the client that handle the intercepted request
func Entry(httpProxyAddr string, httpClient http.Client,
	caKeyPath string, caCrtPath string) (fatalErrorChan chan error, shutdown func()) {
	fatalErrorChan = make(chan error, 1)

	var tcpListener net.Listener = nil
	var httpsServer *http.Server = nil

	// listen tcp
	tcpListener, err := net.Listen("tcp", httpProxyAddr)
	if err != nil {
		fatalErrorChan <- err
		return fatalErrorChan, func() {}
	}

	l.Info("HTTP proxy server running at %s", tcpListener.Addr().String())

	// set system proxy
	if err := utils.SetProxySettings(httpProxyAddr); err == nil {
		l.Info("System proxy has been set to %s", httpProxyAddr)

		var turnOffSystemProxy = func() {
			if err := utils.DisableProxySettings(); err != nil {
				l.Warn("%s", err)
			} else {
				l.Info("System proxy has been closed")
			}
		}

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

	// listen https
	go func() {
		httpsServer, err := confHttpsServer(caKeyPath, caCrtPath, httpClient)
		if err != nil {
			fatalErrorChan <- err
			return
		}
		httpsServer.ListenAndServeTLS("", "")
	}()

	// listen tcp server
	go func() {
		for {
			conn, err := tcpListener.Accept()
			if err != nil {
				// stop accept when closed
				if opErr, ok := err.(*net.OpError); ok && opErr.Err == net.ErrClosed {
					return
				}

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
					handleProxyHTTPS(conn, req)
				} else {
					handleProxyHTTP(conn, req, httpClient)
				}
			}(conn)
		}
	}()

	return fatalErrorChan, func() {
		if tcpListener != nil {
			utils.DisableProxySettings()
			tcpListener.Close()
		}
		if httpsServer != nil {
			httpsServer.Close()
		}
	}
}
