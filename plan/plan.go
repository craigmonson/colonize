package plan

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
	"github.com/craigmonson/colonize/prep"
	"github.com/craigmonson/colonize/util"
)

func Run(c *config.ColonizeConfig, l log.Logger, skipRemote bool) error {
	os.Chdir(c.TmplPath)

	// always run prep first
	prep.Run(c, l)

	if skipRemote {
		l.Log("Skipping remote setup")
	} else {
		l.Log("Running remote setup")
		util.RunCmd("./" + c.CombinedRemoteFilePath)
		l.Log("Disabling remote")
		util.RunCmd("terraform", "remote", "config", "-disable")
	}

	l.Log("Executing terraform plan")
	return util.RunCmd(
		"terraform",
		"plan",
		"-var-file", c.CombinedValsFilePath,
		"-var-file", c.CombinedDerivedValsFilePath,
		"-out", "terraform.tfplan",
	)
}
