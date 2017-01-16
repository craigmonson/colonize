package branch

import (
  "fmt"
  "os"
  "path"
  "path/filepath"

  "github.com/craigmonson/colonize/config"
  "github.com/craigmonson/colonize/log"
  "github.com/craigmonson/colonize/util"
  "github.com/craigmonson/colonize/generate/leaf"
)

type RunArgs struct {
  Name  string
  Leafs []string
}


func Run(c *config.Config, l log.Logger, args interface{}) error {

  runArgs := args.(RunArgs)
  l.Log("Creating Branch: " + runArgs.Name)

  os.Chdir(c.TmplPath)
  os.Mkdir(runArgs.Name, 0755)

  parent_build_order,err := os.OpenFile(c.ConfigFile.Branch_Order_File, os.O_APPEND|os.O_WRONLY, 0664)
  if err != nil {
    l.Log(fmt.Sprintf("Failed to add branch '%s' to '%s'", runArgs.Name, "build_order.txt"))
    os.Exit(-1)
  }
  parent_build_order.WriteString(runArgs.Name + "\n")
  parent_build_order.Close()


  os.Chdir(runArgs.Name)
  os.Mkdir(c.ConfigFile.Environments_Dir, 0755)

  matches,_ := filepath.Glob(path.Join(c.RootPath, c.ConfigFile.Environments_Dir, "*.tfvars"))
  for _,match := range matches {
    util.Touch(c.ConfigFile.Environments_Dir,path.Base(match))
  }

  build_order,err := os.Create(c.ConfigFile.Branch_Order_File)
  if err != nil {
    return err
  }
  defer build_order.Close()


  if len(runArgs.Leafs) > 0 {
    for _,leafName := range runArgs.Leafs {
      err = leaf.Run(c, l, leaf.RunArgs{
        Name: leafName,
        BuildOrder: build_order,
      })
      if err != nil {
        return err
      }
    }
  }


  return nil
}

