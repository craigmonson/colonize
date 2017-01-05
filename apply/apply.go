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

func Run(c *config.Config, l log.Logger, skipRemote bool, remoteAfterApply bool) error {
	os.Chdir(c.TmplPath)

	if skipRemote {
		l.Log("Skipping remote setup")
	} else {
		l.Log("Running remote setup")
		util.RunCmd(c.CombinedRemoteFilePath)
		_,output := util.RunCmd(c.CombinedRemoteFilePath)
		l.Log(output)
	}

	l.Log("Executing terraform apply")
	err,output := util.RunCmd(
		"terraform",
		"apply",
		"-parallelism", "1",
		"-var-file", c.CombinedValsFilePath,
		"-var-file", c.CombinedDerivedValsFilePath,
	)

	l.Log(output)

	if skipRemote {
		if remoteAfterApply {
			remoteUpdate(c,l)
		} else {
			l.Log("Skipping remote setup post-apply. REMOTE NOT SYNC'D!")
		}
	} else {
	  remoteUpdate(c,l)
	}

	return err
}


func remoteUpdate(c *config.Config, l log.Logger) {
	l.Log("Running remote after apply")
	os.Rename("terraform.tfstate", ".terraform/terraform.tfstate")
	util.RunCmd("rm terraform.tfstate.backup")
	_,output := util.RunCmd(c.CombinedRemoteFilePath)
	l.Log(output)
}
