package cmd

import (
	"bufio"
	"os"

	"github.com/spf13/cobra"

	"github.com/craigmonson/colonize/destroy"
)

var yesToDestroy bool

const WARNING_MSG string = `All managed infrastructure will be deleted.
There is no undo. Only entering 'yes' will confirm this operation.

Do you wish to proceed with destroy: `

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy all defined resources in an environment",
	Long: `
This command will perform a "terraform destroy" command on your project for the 
specified environment. In effect, this will destroy all managed resources in the
given leaf or branch that the destroy command is run under.

# Example usage to destroy the "dev" environment: 
$ colonize destroy -e dev

# Example usage to destroy the "dev" environment, say yes to prompt
$ colonize destroy -e dev -y
`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := GetConfig(true)
		if err != nil {
			Log.Log(err.Error())
			os.Exit(-1)
		}

		if !yesToDestroy {
			scan := bufio.NewScanner(os.Stdin)
			Log.Print(WARNING_MSG)
			scan.Scan()

			if scan.Text() != "yes" {
				Log.Log("Destroy operation cancelled by user")
				os.Exit(0)
			}
		}

		err = Run("DESTROY", destroy.Run, conf, Log, true, destroy.RunArgs{
			SkipRemote: SkipRemote,
		})

		if err != nil {
			Log.Log("Destroy failed to run: " + err.Error())
			os.Exit(-1)
		}
	},
}

func init() {
	destroyCmd.Flags().BoolVarP(
		&yesToDestroy,
		"accept",
		"y",
		false,
		"bypass 'accept' prompt by automatically accepting the destruction",
	)
	RootCmd.AddCommand(destroyCmd)
}
