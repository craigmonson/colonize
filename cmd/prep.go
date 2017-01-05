package cmd

import (
	"github.com/spf13/cobra"
	"os"

	"github.com/craigmonson/colonize/prep"
)

// prepCmd represents the prep command
var prepCmd = &cobra.Command{
	Use:   "prep",
	Short: "Generates files that Terraform utilizes",
	Long: `
The prep command is the workhorse of the colonize command. It does all of the combining and tree walking to generate files that the installed terraform will utilize in it's plan / apply / destroy runs. As one would expect, this prepares terraform for the given environment <env>

It can be run on it's own, however, it is run automatically by running the plan/destroy commands.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := GetConfig(true)
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
