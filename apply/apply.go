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

func Run(c *config.ColonizeConfig, l log.Logger, skipRemote bool, remoteAfterApply bool) error {
	os.Chdir(c.TmplPath)

	if skipRemote {
		l.Log("Skipping remote setup")
	} else {
		l.Log("Running remote setup")
		util.RunCmd("./" + c.CombinedRemoteFilePath)
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
			l.Log("Running remote after apply")
			util.RunCmd("mv terraform.tfstate .terraform/.")
			util.RunCmd("rm terraform.tfstate.backup")
			util.RunCmd("./" + c.CombinedRemoteFilePath)
		} else {
			l.Log("Skipping remote setup post-apply. REMOTE NOT SYNC'D!")
		}
	}

	return err
}
