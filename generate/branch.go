package generate

import (
  "fmt"
  "errors"
  "os"

  "github.com/craigmonson/colonize/config"
  "github.com/craigmonson/colonize/log"
)

type RunBranchArgs struct {
  Name  string
  Leafs []string
}


func RunBranch(c *config.ColonizeConfig, l log.Logger, args interface{}) error {

  runArgs := args.(RunBranchArgs)

  if _, err := os.Stat(runArgs.Name); err == nil {
    return errors.New(fmt.Sprintf("Branch with name '%s' already exists", runArgs.Name))
  }

  l.Log("Creating Branch: " + runArgs.Name)

  os.Chdir(c.TmplPath)
  os.Mkdir(runArgs.Name, 0755)
  os.Chdir(runArgs.Name)
  os.Mkdir(c.ConfigFile.Environments_Dir, 0755)

  // TODO: Create <env>.tfvars in <branch_dir>/<env_dir>

  // TODO: Create build_order.txt file from config struct (needs `run-on-branches` branch)
  build_order,err := os.Create("build_order.txt")
  if err != nil {
    return err
  }
  defer build_order.Close()


  if len(runArgs.Leafs) > 0 {
    for _,leaf := range runArgs.Leafs {
      err = RunLeaf(c, l, RunLeafArgs{
        Name: leaf,
      })
      if err != nil {
        return err
      }
      build_order.WriteString(leaf + "\n")
    }
  }


  return nil
}

