package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/log"
)

var Environment string
var Config *config.ColonizeConfig
var Log = log.Log{}
var SkipRemote bool
var RemoteStateAfterApply bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "colonize",
	Short: "A terraform tool to manage environment driven templating",
	Long: `
Colonize is a configurable, albeit opinionated way to organize and manage your
terraform templates. It revolves around the idea of environments, and allows
you to organize templates, and template data around that common idiom.

Once it's been configured, it allows for hierarical templates and variables,
and the ability to organize them in a defined manageable way.`,
}

// This is available for all the subcommands
func GetConfig(requireEnvironment bool) (*config.ColonizeConfig, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if requireEnvironment && Environment == "" {
		return nil, errors.New("environment can not be empty")
	}

	config, err := config.LoadConfigInTree(cwd, Environment)
	if err != nil {
		return nil, err
	}

	return config, err
}


// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVarP(
		&Environment,
		"environment",
		"e",
		"",
		"The environment to colonize")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.PersistentFlags().BoolVarP(&SkipRemote, "skip_remote", "r", false, "skip execution of remote configuration.")
	RootCmd.PersistentFlags().BoolVarP(&SkipRemote, "remote_state_after_apply", "a", false, "Run remote state after terraform apply (if it was skipped).")
}
