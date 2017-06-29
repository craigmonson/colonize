package plan

import (
	"os"

	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/log"
	"github.com/craigmonson/colonize/util"
)

type RunArgs struct {
	SkipRemote bool
}

func Run(c *config.Config, l log.Logger, args interface{}) error {
	runArgs := args.(RunArgs)

	if runArgs.SkipRemote {
		l.Log("Skipping remote setup")
	} else {
		l.Log("Running remote setup")
		util.RunCmd(c.CombinedRemoteFilePath)
	}

	l.Log("Executing terraform plan")
	l.Log(c.CombinedValsFilePath)
	d, _ := os.Getwd()
	l.Log(d)
	err := util.RunCmd(
		"terraform",
		"plan",
		"-var-file", c.CombinedValsFilePath,
		"-var-file", c.CombinedDerivedValsFilePath,
		"-out", "terraform.tfplan",
	)
	return err
}
