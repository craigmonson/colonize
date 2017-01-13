package generate

import (
  "fmt"
  "errors"
  "os"

  "github.com/craigmonson/colonize/config"
  "github.com/craigmonson/colonize/log"
)


func RunBranch(c *config.ColonizeConfig, l log.Logger, name string, leafs []string) error {

  if _, err := os.Stat(name); err == nil {
    return errors.New(fmt.Sprintf("Branch with name '%s' already exists", name))
  }

  l.Log("Creating branch: " + name)

  os.Chdir(c.TmplPath)
  os.Mkdir(name, 0755)
  os.Chdir(name)
  os.Mkdir(c.ConfigFile.Environments_Dir, 0755)

  // TODO: Create <env>.tfvars in <branch_dir>/<env_dir>

  // TODO: Create build_order.txt file from config struct (needs `run-on-branches` branch)
  build_order,err := os.Create("build_order.txt")
  if err != nil {
    return err
  }
  defer build_order.Close()


  if len(leafs) > 0 {
    for _,leaf := range leafs {
      err = RunLeaf(c, l, leaf)
      if err != nil {
        return err
      }
      build_order.WriteString(leaf + "\n")
    }
  }


  return nil
}

