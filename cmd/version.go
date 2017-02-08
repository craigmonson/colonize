package cmd

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Colonize",
	Long:  `Retrieve the version information for Colonize`,
	Run: func(cmd *cobra.Command, args []string) {
		Log.Log("Colonize v0.1.0-alpha")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
