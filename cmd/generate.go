package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/craigmonson/colonize/generate"
	"github.com/craigmonson/colonize/generate/branch"
	"github.com/craigmonson/colonize/generate/leaf"
	"github.com/spf13/cobra"
)

var branchLeafs string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Colonize project structure and resources",
	Long: `
Use this Command suite to generate Colonzie resources and project structures.
This command can be very useful for quickly creating and stubbing out project
structures.

Example usage:
colonize generate <resource-type> [options]

Example branch generation:
colonize generate branch mybranch
colonize generate branch mybranch --leafs myleaf1,myleaf2

Example leaf generation:
colonize generate leaf myleaf
  `,
}

var genBranchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Generate a Colonize Branch in your project",
	Long: `
The branch generation sub-command is used to generate a Colonize branch,
including build order file, environment directory, environment tfvars, and
optionally a list of leafs underneath the branch.

Example usage:
colonize generate branch mybranch
colonize generate branch mybranch --leafs myleaf1,myleaf2
  `,
	Run: func(cmd *cobra.Command, args []string) {

		err := generate.ValidateArgsLength("Branch", args, 1, 1)
		if err != nil {
			CompleteFail(err.Error())
		}

		name := args[0]
		leafs := strings.Split(branchLeafs, ",")

		err = generate.ValidateNameAvailable("Branch", name)
		if err != nil {
			CompleteFail(err.Error())
		}

		conf, err := GetConfig(false)
		if err != nil {
			CompleteFail(err.Error())
		}

		err = branch.Run(conf, Log, branch.RunArgs{
			Name:  name,
			Leafs: leafs,
		})
		if err != nil {
			CompleteFail("Generate Branch failed to run: " + err.Error())
		}

		CompleteSucceed()

	},
}

var genLeafCmd = &cobra.Command{
	Use:   "leaf",
	Short: "Generate a Colonize Leaf in your project",
	Long: `
Create a Colonize Leaf underneath the current Colonize branch. This will
create the leaf directory, a main.tf file inside the leaf, and append the
leaf name to the parent branches build order file.

Example:
colonize generate leaf myleaf
  `,
	Run: func(cmd *cobra.Command, args []string) {

		err := generate.ValidateArgsLength("Leaf", args, 1, 1)
		if err != nil {
			CompleteFail(err.Error())
		}

		name := args[0]
		err = generate.ValidateNameAvailable("Leaf", name)
		if err != nil {
			CompleteFail(err.Error())
		}

		conf, err := GetConfig(false)
		if err != nil {
			CompleteFail(err.Error())
		}

		build_order, err := os.OpenFile(conf.ConfigFile.Branch_Order_File, os.O_APPEND|os.O_WRONLY, 0664)
		if err != nil {
			CompleteFail(fmt.Sprintf("Failed to add leaf '%s' to '%s'", name, "build_order.txt"))
		}
		defer build_order.Close()

		err = leaf.Run(conf, Log, leaf.RunArgs{
			Name:       name,
			BuildOrder: build_order,
		})
		if err != nil {
			CompleteFail("Generate Leaf failed to run: " + err.Error())
		}

		CompleteSucceed()

	},
}

func constructGenerateBranchCommand() {

	genBranchCmd.Flags().StringVarP(
		&branchLeafs,
		"leafs",
		"l",
		"",
		"CSV list of leaf names to create under the created branch",
	)
	generateCmd.AddCommand(genBranchCmd)
}

func constructGenerateLeafCommand() {
	generateCmd.AddCommand(genLeafCmd)
}

func init() {
	constructGenerateBranchCommand()
	constructGenerateLeafCommand()
	RootCmd.AddCommand(generateCmd)
}
