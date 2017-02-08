package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/log"
	"github.com/craigmonson/colonize/util"
)

var Config *config.Config
var Log = log.Log{}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "colonize",
	Short: "A terraform tool to manage environment driven templating",
	Long: `
Colonize is a configurable, albeit opinionated way to organize and manage your
terraform templates. It revolves around the idea of environments, and allows
you to organize templates, and template data around that common organizational
structure.

Once it's been configured, it allows for hierarical templates and variables,
and the ability to organize them in a defined manageable way.`,
}

// This is available for all the subcommands
func GetConfigWithoutEnvironment() (*config.Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	Log.LogPretty(util.PadRight("\nCOLONIZE ", "*", 79), color.Bold)

	config, err := config.LoadConfigInTree(cwd, "")
	if err != nil {
		return nil, err
	}

	return config, err
}

func GetConfig(environment string) (*config.Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if environment == "" {
		return nil, errors.New("environment can not be empty")
	}
	Log.LogPretty(util.PadRight(fmt.Sprintf("\nCOLONIZE [%s] ", environment), "*", 79), color.Bold)

	config, err := config.LoadConfigInTree(cwd, environment)
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

func CompleteSucceed() {
	Log.LogPretty(util.PadRight("\nRECAP ", "*", 79), color.Bold)
	Log.LogPretty("Completed successfully!", color.FgGreen)
	os.Exit(0)
}

func CompleteFail(err string) {
	Log.LogPretty(util.PadRight("\nRECAP ", "*", 79), color.Bold)
	Log.LogPretty(err, color.FgRed)
	os.Exit(-1)
}

func addEnvironmentFlag(cmd *cobra.Command, environment *string) {

	cmd.Flags().StringVarP(
		environment,
		"environment",
		"e",
		"",
		"The environment to colonize",
	)
}

func addSkipRemoteFlag(cmd *cobra.Command, skipRemote *bool) {

	cmd.Flags().BoolVarP(
		skipRemote,
		"skip-remote",
		"k",
		false,
		"skip execution of remote configuration.",
	)
}

func init() {

}
