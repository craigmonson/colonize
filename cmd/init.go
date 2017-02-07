package cmd

import (
	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/initialize"
	"github.com/spf13/cobra"
)

var acceptDefaults bool
var initEnvironments string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a Colonize project",
	Long:  `This command is used to aid in the generation the .colonize.yaml configuration file and project directory structure.`,
	Run: func(cmd *cobra.Command, args []string) {

		_, err := GetConfigWithoutEnvironment()
		if err == nil {
			CompleteFail("Colonize project already initialized. Exiting.")
		}

		var blankConfig config.Config
		err = initialize.Run(&blankConfig, Log, initialize.RunArgs{
			AcceptDefaults:   acceptDefaults,
			InitEnvironments: initEnvironments,
		})
		if err != nil {
			CompleteFail(err.Error())
		}

		CompleteSucceed()
	},
}

func init() {
	initCmd.Flags().BoolVarP(
		&acceptDefaults,
		"accept-defaults",
		"",
		false,
		"automatically accepts default values, skipping manual setup",
	)
	initCmd.Flags().StringVarP(
		&initEnvironments,
		"environments",
		"",
		"",
		"specify environment files to initialize the project with",
	)
	RootCmd.AddCommand(initCmd)
}
