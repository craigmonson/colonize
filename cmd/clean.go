package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/craigmonson/colonize/clean"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes all Colonize prep data",
	Long: `
This command will remove all temproary and derived files that Colonize
has created, via its "prep" command. This can be run at the leaf or the
branch level.

Example usage to clean a project:
$ colonize clean
	`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := GetConfig(false)
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
