package cmd

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/yieldray/middleman/cmd/flags"
	"github.com/yieldray/middleman/interceptor"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy <proxy-url>",
	Short: "使用代理服务器，例如：https://cros.deno.dev/",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "未指定代理服务器")
			cmd.Help()
			os.Exit(1)
		}
		proxyServer, err := interceptor.NewProxyServer(args[0])
		l.Debug("proxyServer=%s", proxyServer)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			cmd.Help()
			os.Exit(1)
		}

		// the real client to send request
		// stateless, no jar, no auto redirect
		httpClient := http.Client{
			Jar: nil,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		httpProxyClient := http.Client{
			Transport: interceptor.NewRoundTripper(func(req *http.Request) (*http.Response, error) {
				targetUrl := req.URL.String()
				u := proxyServer.ProxyURL(targetUrl)

				l.Debug("%s", targetUrl)

				request := &http.Request{
					Method:        req.Method,
					URL:           u,
					Header:        req.Header,
					Body:          req.Body,
					ContentLength: req.ContentLength,
					Close:         req.Close,
					Trailer:       req.Trailer,
				}

				response, err := httpClient.Do(request)

				if flags.Log {
					file, err := os.OpenFile(flags.LogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
					if err != nil {
						l.Error("OpenFile: %s", err)
						return response, err
					}

					now := time.Now()
					file.WriteString(fmt.Sprintf("[%s]\n", now.Format("2006-01-02 15:03:04")))

					//req
					if buf, err := httputil.DumpRequestOut(request, false); err == nil {
						file.Write(buf)
					} else {
						l.Error("DumpRequestOut: %s", err)
					}

					//res
					if buf, err := httputil.DumpResponse(response, false); err == nil {
						file.Write(buf)
					} else {
						l.Error("DumpResponse: %s", err)
					}

					file.WriteString("\n\n")
				}

				return response, err
			}),
		}

		interceptor.Entry(flags.GetAddr(), httpProxyClient, flags.CaKey, flags.CaCrt)
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)
}
