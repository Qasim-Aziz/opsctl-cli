package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "opsctl",
	Short: "A command-line tool for operations.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to opsctl! Use subcommands like 'cpu-stress' or 'net-tools' or 'api-stress' or 'dns-check")
	},
}

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize()
	// Add any global flags or setup for the root command here.
}
