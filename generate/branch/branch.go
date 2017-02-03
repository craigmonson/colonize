package branch

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/generate/leaf"
	"github.com/craigmonson/colonize/log"
	"github.com/craigmonson/colonize/util"
)

type RunArgs struct {
	Name  string
	Leafs []string
}

func Run(c *config.Config, l log.Logger, args interface{}) error {

	runArgs := args.(RunArgs)
	l.LogPretty(util.PadRight(fmt.Sprintf("\nGENERATE [Branch | %s] ", runArgs.Name), "*", 79), color.Bold)

	os.Chdir(c.TmplPath)

	parent_build_order, err := os.OpenFile(c.ConfigFile.Branch_Order_File, os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to add branch '%s' to '%s'", runArgs.Name, c.ConfigFile.Branch_Order_File))
		os.Exit(-1)
	}
	defer parent_build_order.Close()

	l.Log("Creating branch directory...")
	os.Mkdir(runArgs.Name, 0755)

	l.Log("Creating branch environment directory...")
	os.Chdir(runArgs.Name)
	os.Mkdir(c.ConfigFile.Environments_Dir, 0755)

	l.Log("Creating environment variables files...")
	matches, _ := filepath.Glob(path.Join(c.RootPath, c.ConfigFile.Environments_Dir, "*.tfvars"))
	for _, match := range matches {
		util.Touch(c.ConfigFile.Environments_Dir, path.Base(match))
	}

	l.Log("Creating branch build order file...")
	build_order, err := os.Create(c.ConfigFile.Branch_Order_File)
	if err != nil {
		return err
	}
	defer build_order.Close()

	l.Log("Updating parent build order file...")
	parent_build_order.WriteString(runArgs.Name + "\n")

	if len(runArgs.Leafs) > 0 {
		for _, leafName := range runArgs.Leafs {
			if leafName != "" {
				err = leaf.Run(c, l, leaf.RunArgs{
					Name:       leafName,
					BuildOrder: build_order,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
