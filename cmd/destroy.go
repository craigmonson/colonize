package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/craigmonson/colonize/destroy"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy all defined resources in an environment",
	Long: `
This command will perform a "terraform destroy" command on your project for the 
specified environment. In effect, this will destroy all managed resources in the
give leaf or branch that the destroy command is run under

Example usage to destroy the "dev" environment: 
$ colonize destroy -e dev`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := GetConfig(true)
		if err != nil {
			Log.Log(err.Error())
			os.Exit(-1)
		}

                err,output := destroy.Run(conf, Log, SkipRemote)
                Log.Log(output)

		if err != nil {
			Log.Log("Destroy failed to run: " + err.Error())
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(destroyCmd)
}
