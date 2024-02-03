package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/yieldray/middleman/cmd/flags"
	"github.com/yieldray/middleman/interceptor"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [<sqlite-file-path>]",
	Short: "Inspect http(s) traffic",
	Run: func(cmd *cobra.Command, args []string) {
		var dbPath string
		if len(args) < 1 {
			dbPath = "middleman.db"
		} else {
			dbPath = args[0]
		}

		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			l.Fatal("failed to connect database")
		}
		db.AutoMigrate(&RowSchema{})

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
				now := time.Now()

				var reqAllBody []byte
				if req.Body != nil {
					if reqAllBody, err = io.ReadAll(req.Body); err != nil {
						l.Error("req.Body %s", err)
					} else {
						req.Body = io.NopCloser(bytes.NewReader(reqAllBody)) // clone body
					}
				}
				var reqHeaders map[string][]string = req.Header.Clone()

				res, err := httpClient.Do(req)

				if err != nil {
					l.Error("Do %s", err)
					return res, err
				}

				var resAllBody []byte
				if res.Body != nil {
					if resAllBody, err = io.ReadAll(res.Body); err != nil {
						l.Error("res.Body %s", err)
					} else {
						res.Body = io.NopCloser(bytes.NewReader(resAllBody)) // clone body
					}
				}
				var resHeaders map[string][]string = res.Header.Clone()

				db.Create(&RowSchema{
					Time:            now,
					URL:             req.URL.String(),
					RequestTarget:   fmt.Sprintf("%s %s %s", req.Method, req.RequestURI, req.Proto),
					RequestHeaders:  reqHeaders,
					RequestBody:     string(reqAllBody),
					ResponseTarget:  fmt.Sprintf("%s %s", res.Proto, res.Status),
					ResponseHeaders: resHeaders,
					ResponseBody:    string(resAllBody),
				})

				return res, err
			}),
		}

		interceptor.Entry(flags.GetAddr(), httpProxyClient, flags.CaKey, flags.CaCrt)

	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}

type RowSchema struct {
	Time            time.Time
	URL             string
	RequestTarget   string
	RequestHeaders  map[string][]string `gorm:"serializer:json"`
	RequestBody     string
	ResponseTarget  string
	ResponseHeaders map[string][]string `gorm:"serializer:json"`
	ResponseBody    string
}
