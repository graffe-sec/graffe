package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "graffe",
	Short: "Reviews cloud configurations to identify security risks",
	Long: `Graffe is a CLI app that allows dev/sec/ops professionals review their cloud security
without having to pay an arm and a leg for an enterprise report generator`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() error {
	return rootCmd.Execute()
}
