package cmd

import (
  "os"

  "github.com/spf13/cobra"
  "github.com/craigmonson/colonize/config"
  "github.com/craigmonson/colonize/initialize"
)

var initCmd = &cobra.Command{
  Use: "init",
  Short: "Initialize a Colonize project",
  Long: `This command is used to aid in the generation the .colonize.yaml configuration file and project directory structure.`,
  Run: func(cmd *cobra.Command, args []string) {

    _, err := GetConfig(false)
    if err == nil {
      Log.Log("Colonize project already initialized")
      os.Exit(0)
    }

    var blankConfig config.ColonizeConfig
    err = initialize.Run(&blankConfig, Log)
    if err != nil {
      Log.Log("Init failed to run: " + err.Error())
      os.Exit(-1)
    }
  },
}

func init() {
  RootCmd.AddCommand(initCmd)
}
