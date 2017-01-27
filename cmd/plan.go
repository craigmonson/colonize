package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craigmonson/colonize/plan"
	"github.com/craigmonson/colonize/prep"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Plan a Terraform run for a specific environment",
	Long: `
This command will perform a "terraform plan" command on your project for the
specified environment. This will generate a plan of changes that will be be
applied to your environment when the "apply" command is run.

# Example usage to plan the "dev" environment:
$ colonize plan -e dev

# Plan, but don't sync to remote
$ colonize plan -e dev --skip-remote

# Plan, skip initial remote sync, but sync after apply has finished
$ colonize plan -e dev --skip-remote --remote-state-after-apply
	`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := GetConfig(true)
		if err != nil {
			CompleteFail(err.Error())
		}

		err = Run("PREP", prep.Run, conf, Log, false, nil)
		if err != nil {
			CompleteFail("Prep failed to run: " + err.Error())
		}

		err = Run("PLAN", plan.Run, conf, Log, false, plan.RunArgs{
			SkipRemote: SkipRemote,
		})

		if err != nil {
			CompleteFail("Plan failed to run: " + err.Error())
		}

		CompleteSucceed()
	},
}

func init() {
	RootCmd.AddCommand(planCmd)
}
