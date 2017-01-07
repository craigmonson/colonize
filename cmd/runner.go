package cmd

import (
	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/log"
)

// reverse is here because the destroy command will need to destroy stuff in
// reverse order.
func Run(f func(*config.Config, log.Logger, interface{}) error, c *config.Config, l log.Logger, reverse bool, args interface{}) error {
	if c.IsBranch() {
		return RunBranch(f, c, l, reverse, args)
	}

	l.Log("Running " + c.TmplName)
	return f(c, l, args)
}

func RunBranch(f func(*config.Config, log.Logger, interface{}) error, c *config.Config, l log.Logger, reverse bool, args interface{}) error {
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
		newConf, err := config.LoadConfigInTree(p, c.Environment)
		if err != nil {
			return err
		}

		if err := Run(f, newConf, l, reverse, args); err != nil {
			return err
		}
	}

	return nil
}
