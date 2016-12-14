package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/craigmonson/colonize/apply"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := GetConfig()
		if err != nil {
			Log.Log(err.Error())
			os.Exit(-1)
		}
		err = apply.Run(conf, Log, SkipRemote, RemoteStateAfterApply)
		if err != nil {
			Log.Log("Apply failed to run: " + err.Error())
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(applyCmd)
}
