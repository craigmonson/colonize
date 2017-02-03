package cmd

import (
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
or branch that apply command is run under.  You may also run this command alone,
and a plan will be executed before the apply.

# Example usage to apply changes:
$ colonize apply --environment dev

# Example usage to apply changes but skip setting up the remote
$ colonize apply --environment dev --skip-remote

# Example usage to apply changes, skip initial remote sync, then sync after.
$ colonize apply --environment dev --skip-remote --remote-state-after-apply
        `,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := GetConfig(true)
		if err != nil {
			CompleteFail(err.Error())
		}

		err = Run("APPLY", apply.Run, conf, Log, false, apply.RunArgs{
			SkipRemote:            SkipRemote,
			RemoteStateAfterApply: RemoteStateAfterApply,
		})
		if err != nil {
			CompleteFail("Apply failed to run: " + err.Error())
		}

		CompleteSucceed()
	},
}

func init() {
	RootCmd.AddCommand(applyCmd)
}
