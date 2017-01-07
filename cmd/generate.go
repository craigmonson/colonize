package cmd

import (
  "os"
  "strings"
  "github.com/spf13/cobra"
//  "github.com/craigmonson/colonize/config"
  "github.com/craigmonson/colonize/generate"
)


var branchLeafs string

var generateCmd = &cobra.Command{
  Use: "generate",
  Short: "Generate Colonize project structure and ",
  Long:`

  `,
}

var genBranchCmd = &cobra.Command{
  Use: "branch",
  Short: "Generate a Colonize Branch in your project",
  Long: `
Placeholder text
  `,
  Run: func(cmd *cobra.Command, args []string) {

    if len(args) < 1 {
      Log.Log("You must specificy a branch name to create")
      os.Exit(-1)
    } else if len(args) > 1 {
      Log.Log("You may only specify a single branch to create at a time")
      os.Exit(-1)
    }

    conf, err := GetConfig(false)
    if err != nil {
      Log.Log(err.Error())
      os.Exit(-1)
    }

    branch := args[0]
    Log.Log("Creating branch: " + branch)

    leafs := strings.Split(branchLeafs,",")
    Log.Log("Creating leafs: " + strings.Join(leafs,","))

    err = generate.RunBranch(conf, Log, leafs)
    if err != nil {
      Log.Log("Generate Branch failed to run: " + err.Error())
      os.Exit(-1)
    }

  },
}

var genLeafCmd = &cobra.Command{
  Use: "leaf",
  Short: "Generate a Colonize Leaf in your project",
  Long: `
Placeholder text
  `,
  Run: func(cmd *cobra.Command, args []string) {

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
