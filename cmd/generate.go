package cmd

import (
  "fmt"
  "os"
  "strings"
  "github.com/spf13/cobra"
  "github.com/craigmonson/colonize/generate"
)


var branchLeafs string

var generateCmd = &cobra.Command{
  Use: "generate",
  Short: "Generate Colonize project structure and resources",
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
    leafs := strings.Split(branchLeafs,",")

    err = generate.RunBranch(conf, Log, branch, leafs)
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

    if len(args) < 1 {
      Log.Log("You must specify a leaf name to create")
      os.Exit(-1)
    } else if len(args) > 1 {
      Log.Log("You may only create a single leaf at a time")
      os.Exit(-1)
    }

    conf, err := GetConfig(false)
    if err != nil {
      Log.Log(err.Error())
      os.Exit(-1)
    }

    leaf := args[0]
    err = generate.RunLeaf(conf, Log, leaf)
    if err != nil {
      Log.Log("Generate Leaf failed to run: " + err.Error())
      os.Exit(-1)
    }

    build_order,err := os.Create("build_order.txt")
    if err != nil {
      // TODO: Pull file name from config struct (requires run-on-branches code)
      Log.Log(fmt.Sprintf("Failed to add leaf '%s' to '%s'", leaf, "build_order.txt"))
      os.Exit(-1)
    }
    defer build_order.Close()

    build_order.WriteString(leaf + "\n")
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
