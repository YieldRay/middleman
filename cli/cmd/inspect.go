package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yieldray/middleman/cli/impl"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [<sqlite-file-path>]",
	Short: "Inspect http(s) traffic and save as sqlite file",
	Run: func(cmd *cobra.Command, args []string) {
		var dbPath string
		if len(args) < 1 {
			dbPath = "middleman.db"
		} else {
			dbPath = args[0]
		}

		fatalErrorChan, shutdown := impl.Inspect(dbPath, "", "", "")

		<-fatalErrorChan
		shutdown()
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
