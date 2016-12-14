package cmd

import (
	"github.com/spf13/cobra"
	"os"

	"github.com/craigmonson/colonize/prep"
)

// prepCmd represents the prep command
var prepCmd = &cobra.Command{
	Use:   "prep",
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
		err = prep.Run(conf, Log)
		if err != nil {
			Log.Log("Prep Failed to Run: " + err.Error())
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(prepCmd)
}
