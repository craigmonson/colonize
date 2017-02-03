package destroy

import (
	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/log"
	"github.com/craigmonson/colonize/prep"
	"github.com/craigmonson/colonize/util"
)

type RunArgs struct {
	SkipRemote bool
}

func Run(c *config.Config, l log.Logger, args interface{}) error {
	runArgs := args.(RunArgs)

	// always run prep first
	prep.Run(c, l, nil)

	if runArgs.SkipRemote {
		l.Log("Skipping remote setup")
	} else {
		l.Log("Running remote setup")
		util.RunCmd(c.CombinedRemoteFilePath)
	}

	l.Log("Executing terraform destroy")
	return util.RunCmd(
		"terraform",
		"destroy",
		"-force",
		"-var-file", c.CombinedValsFilePath,
		"-var-file", c.CombinedDerivedValsFilePath,
	)
}
