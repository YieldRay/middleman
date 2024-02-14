package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yieldray/middleman/cli/impl"
	"github.com/yieldray/middleman/cli/interceptor"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy <proxy-url>",
	Short: "Using a proxy server, For example: https://cros.deno.dev/",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Not specifying a proxy server!")
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

		fatalErrorChan, shutdown := impl.Proxy(proxyServer, "", "", "")

		<-fatalErrorChan
		shutdown()
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)
}
