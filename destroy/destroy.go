package destroy

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

func Run(c *config.Config, l log.Logger, skipRemote bool) (error,string) {
	os.Chdir(c.TmplPath)

	// always run prep first
	prep.Run(c, l)

	if skipRemote {
		l.Log("Skipping remote setup")
	} else {
		l.Log("Running remote setup")
		util.RunCmd("./" + c.CombinedRemoteFilePath)
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
