package cmd

import (
        "fmt"
	"os"

	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/log"
        "github.com/craigmonson/colonize/util"

)

// reverse is here because the destroy command will need to destroy stuff in
// reverse order.
func Run(name string, f func(*config.Config, log.Logger, interface{}) error, c *config.Config, l log.Logger, reverse bool, args interface{}) error {
	if c.IsBranch() {
		return RunBranch(name, f, c, l, reverse, args)
	}

        l.Log(util.PadRight(fmt.Sprintf("\n%s [%s] ", name, c.TmplPath),"*",79))
	return f(c, l, args)
}

func RunBranch(name string, f func(*config.Config, log.Logger, interface{}) error, c *config.Config, l log.Logger, reverse bool, args interface{}) error {
	buildPaths, err := c.GetBuildOrderPaths()
	if err != nil {
		return err
	}

	if reverse {
		for i, j := 0, len(buildPaths)-1; i < j; i, j = i+1, j-1 {
			buildPaths[i], buildPaths[j] = buildPaths[j], buildPaths[i]
		}
	}

	for _, p := range buildPaths {
		if err := os.Chdir(p); err != nil {
			return err
		}
		newConf, err := config.LoadConfigInTree(p, c.Environment)
		if err != nil {
			return err
		}

		if err := Run(name, f, newConf, l, reverse, args); err != nil {
			return err
		}
	}

	return nil
}
