package plan

import (
	"os"

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
	os.Chdir(c.TmplPath)

	// always run prep first
	prep.Run(c, l, nil)

	if runArgs.SkipRemote {
		l.Log("Skipping remote setup")
	} else {
		l.Log("Running remote setup")
		util.RunCmd("./" + c.CombinedRemoteFilePath)
		l.Log("Disabling remote")
		util.RunCmd("terraform", "remote", "config", "-disable")
	}

	l.Log("Executing terraform plan")
	err, out := util.RunCmd(
		"terraform",
		"plan",
		"-var-file", c.CombinedValsFilePath,
		"-var-file", c.CombinedDerivedValsFilePath,
		"-out", "terraform.tfplan",
	)
	l.Log(out)
	return err
}
