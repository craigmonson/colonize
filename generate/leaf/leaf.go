package leaf

import (
  "os"

  "github.com/craigmonson/colonize/config"
  "github.com/craigmonson/colonize/log"
  "github.com/craigmonson/colonize/util"
)

type RunArgs struct {
  Name string
  BuildOrder *os.File
}

func Run(c *config.ColonizeConfig, l log.Logger, args interface{}) error {
  runArgs := args.(RunArgs)
  l.Log("Creating Leaf: " + runArgs.Name)

  os.Mkdir(runArgs.Name, 0755)
  util.Touch(runArgs.Name,"main.tf")
  runArgs.BuildOrder.WriteString(runArgs.Name + "\n")

  return nil
}
