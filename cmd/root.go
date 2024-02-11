package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/mborders/logmatic"
	"github.com/spf13/cobra"
	"github.com/yieldray/middleman/cmd/flags"
)

var rootCmd = &cobra.Command{
	Use:   "middleman",
	Short: "Middleman is a http(s) interceptor",
	Long:  `CA (Certificate Authority) uses locally generated self-signed certificates and keys. Make sure the system trusts the self-signed certificates.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		l.SetLevel(logmatic.LogLevel(flags.LogLevel))
	},
}

var l = logmatic.NewLogger()

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flags.CaCrt, "ca-crt", "ca.crt", "CA cert path (openssl req -x509 -new -key ca.key -out ca.crt -days 3650)")
	rootCmd.PersistentFlags().StringVar(&flags.CaKey, "ca-key", "ca.key", "CA key path (openssl genpkey -algorithm RSA -out ca.key)")
	rootCmd.PersistentFlags().Uint8Var(&flags.LogLevel, "log-level", 2, "Set the log level, TRACE|DEBUG|INFO|WARN|ERROR|FATAL") // default is INFO
	rootCmd.PersistentFlags().IntVar(&flags.Port, "port", 9980, "http proxy local port")
	rootCmd.PersistentFlags().BoolVar(&flags.Expose, "expose", false, "expose local server")
	rootCmd.PersistentFlags().BoolVar(&flags.Log, "log", false, "write log to file")
	rootCmd.PersistentFlags().StringVar(&flags.LogPath, "log-path", fmt.Sprintf("middleman_%s.log", time.Now().Format("2006-01-02")), "path to log file of request")
}
