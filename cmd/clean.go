package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/craigmonson/colonize/clean"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes all Colonize prep data, and .terraform data",
	Long: `
This command will remove all temproary and derived files that Colonize
has created, via its "prep" command. This can be run at the leaf or the
branch level.  Note that this will also remove the .terraform directory, so
all module and state data will be removed.

Example usage to clean a project:
$ colonize clean
	`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := GetConfig(false)
		if err != nil {
			Log.Log(err.Error())
			os.Exit(-1)
		}
		err = Run("CLEAN", clean.Run, conf, Log, false, nil)
		if err != nil {
			Log.Log("Clean failed to run: " + err.Error())
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(cleanCmd)
}
