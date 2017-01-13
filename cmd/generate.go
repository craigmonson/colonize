package cmd

import (
  "fmt"
  "os"
  "strings"

  "github.com/spf13/cobra"
  "github.com/craigmonson/colonize/generate"
  "github.com/craigmonson/colonize/generate/branch"
  "github.com/craigmonson/colonize/generate/leaf"
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

    err := generate.ValidateArgsLength("Branch", args, 1, 1)
    if err != nil {
      Log.Log(err.Error())
      os.Exit(-1)
    }

    name := args[0]
    leafs := strings.Split(branchLeafs,",")

    err = generate.ValidateNameAvailable("Branch", name)
    if err != nil {
      Log.Log(err.Error())
      os.Exit(-1)
    }

    conf, err := GetConfig(false)
    if err != nil {
      Log.Log(err.Error())
      os.Exit(-1)
    }

    err = branch.Run(conf, Log, branch.RunArgs{
      Name: name,
      Leafs: leafs,
    })
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

    err := generate.ValidateArgsLength("Leaf", args, 1, 1)
    if err != nil {
      Log.Log(err.Error())
      os.Exit(-1)
    }

    name := args[0]
    err = generate.ValidateNameAvailable("Leaf", name)
    if err != nil {
      Log.Log(err.Error())
      os.Exit(-1)
    }

    conf, err := GetConfig(false)
    if err != nil {
      Log.Log(err.Error())
      os.Exit(-1)
    }

    // TODO: Pull file name from config struct (requires run-on-branches code)
    build_order,err := os.OpenFile("build_order.txt", os.O_APPEND|os.O_WRONLY, 0664)
    if err != nil {
      Log.Log(fmt.Sprintf("Failed to add leaf '%s' to '%s'", name, "build_order.txt"))
      os.Exit(-1)
    }
    defer build_order.Close()

    err = leaf.Run(conf, Log, leaf.RunArgs{
      Name: name,
      BuildOrder: build_order,
    })
    if err != nil {
      Log.Log("Generate Leaf failed to run: " + err.Error())
      os.Exit(-1)
    }

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
