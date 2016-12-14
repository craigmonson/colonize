package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/craigmonson/colonize/clean"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
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
		err = clean.Run(conf, Log)
		if err != nil {
			Log.Log("Clean failed to run: " + err.Error())
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(cleanCmd)
}
