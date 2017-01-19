package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/craigmonson/colonize/apply"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies the Terraform plan to the target environment",
	Long: `
This command will perform a "terraform apply" command on your project for the
existing Terraform plan. In effect, this will create/update/remove any managed
resources according to the output of the "plan" command, for the given leaf
or branch that apply command is run under

Example usage to apply changes:
$ colonize apply
        `,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := GetConfig(true)
		if err != nil {
			Log.Log(err.Error())
			os.Exit(-1)
		}
		err = Run(apply.Run, conf, Log, false, apply.RunArgs{
			SkipRemote:            SkipRemote,
			RemoteStateAfterApply: RemoteStateAfterApply,
		})
		if err != nil {
			Log.Log("Apply failed to run: " + err.Error())
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(applyCmd)
}
