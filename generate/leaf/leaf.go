package leaf

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/log"
	"github.com/craigmonson/colonize/util"
)

type RunArgs struct {
	Name       string
	BuildOrder *os.File
}

func Run(c *config.Config, l log.Logger, args interface{}) error {
	runArgs := args.(RunArgs)

	cwd, _ := os.Getwd()
	branch := util.GetBasename(cwd)
	l.LogPretty(util.PadRight(fmt.Sprintf("\nGENERATE [Leaf | %s/%s] ", branch, runArgs.Name), "*", 79), color.Bold)

	l.Log("Creating leaf directory...")
	os.Mkdir(runArgs.Name, 0755)

	l.Log("Creating main leaf terraform template...")
	util.Touch(runArgs.Name, "main.tf")

	l.Log("Updating branch build order file...")
	runArgs.BuildOrder.WriteString(runArgs.Name + "\n")

	return nil
}
