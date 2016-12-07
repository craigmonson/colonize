package cmd

import (
	"github.com/spf13/cobra"

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
		err := prep.Run(Config, Environment)
		if err != nil {
			panic("Prep Failed to Run: " + err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(prepCmd)
}
