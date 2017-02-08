package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craigmonson/colonize/prep"
)

type PrepFlags struct {
	Environment string
}

var prepFlags = PrepFlags{}

// prepCmd represents the prep command
var prepCmd = &cobra.Command{
	Use:   "prep",
	Short: "Generates files that Terraform utilizes",
	Long: `
The prep command is the workhorse of the colonize command. It does all of the
combining and tree walking to generate files that the installed terraform will
utilize in it's plan / apply / destroy runs. As one would expect, this prepares
terraform for the given environment <env>  Note that colonize doesn't do
anything "special" beyond creating files that terraform will consume.  You
can safely skip the other colonize commands (prep/apply/destroy) and run ANY
terraform commands "by hand" after running prep.  In fact, while debugging
your template code, you may find that pattern very helpful.

Note that it can be run on it's own, however, it is run automatically by
running the plan/destroy commands, as there is no equivalent command in
terraform.

# Example preparing a colonization
$ colonize prep -e dev
	`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := GetConfig(prepFlags.Environment)
		if err != nil {
			CompleteFail(err.Error())
		}
		err = Run("PREP", prep.Run, conf, Log, false, nil)
		if err != nil {
			CompleteFail("Prep Failed to Run: " + err.Error())
		}

		CompleteSucceed()
	},
}

func init() {
	addEnvironmentFlag(prepCmd, &prepFlags.Environment)
	RootCmd.AddCommand(prepCmd)
}
