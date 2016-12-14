package cmd

import (
	"github.com/spf13/cobra"

	"github.com/craigmonson/colonize/plan"
)

var SkipRemote bool

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := GetConfig()
		if err != nil {
			Log.Log(err.Error())
			os.Exit(-1)
		}
		err = plan.Run(conf, Log, SkipRemote)
		if err != nil {
			Log.Log("Plan failed to run: " + err.Error())
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(planCmd)

	RootCmd.Flags().BoolVarP(&SkipRemote, "skip_remote", "r", false, "skip execution of remote configuration.")
}
