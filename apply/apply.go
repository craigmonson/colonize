package apply

import (
	//	"bytes"
	//	//"fmt"
	//	"io/ioutil"
	"os"
	//	"regexp"
	//	"strings"
	//
	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/log"
	"github.com/craigmonson/colonize/util"
)

type RunArgs struct {
	SkipRemote            bool
	RemoteStateAfterApply bool
}

func Run(c *config.Config, l log.Logger, args interface{}) error {
	runArgs := args.(RunArgs)

	if runArgs.SkipRemote {
		l.Log("Skipping remote setup")
	} else {
		l.Log("Running remote setup")
		util.RunCmd(c.CombinedRemoteFilePath)
	}

	l.Log("Executing terraform apply")
	err := util.RunCmd(
		"terraform",
		"apply",
		"-parallelism", "1",
		"-var-file", c.CombinedValsFilePath,
		"-var-file", c.CombinedDerivedValsFilePath,
	)

	if runArgs.SkipRemote {
		if runArgs.RemoteStateAfterApply {
			remoteUpdate(c, l)
		} else {
			l.Log("Skipping remote setup post-apply. REMOTE NOT SYNC'D!")
		}
	} else {
		remoteUpdate(c, l)
	}

	return err
}

func remoteUpdate(c *config.Config, l log.Logger) {
	l.Log("Running remote after apply")
	os.Rename("terraform.tfstate", ".terraform/terraform.tfstate")
	util.RunCmd("rm terraform.tfstate.backup")
	util.RunCmd(c.CombinedRemoteFilePath)
}
