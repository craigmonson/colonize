package generate

import (
  "fmt"
  "errors"
  "os"

  "github.com/craigmonson/colonize/config"
  "github.com/craigmonson/colonize/log"
)

type RunLeafArgs struct {
  Name string
}

func RunLeaf(c *config.ColonizeConfig, l log.Logger, args interface{}) error {
  runArgs := args.(RunLeafArgs)

  //TODO: Make sure we're working in a branch

  if _, err := os.Stat(runArgs.Name); err == nil {
    return errors.New("Leaf already exists")
  }

  l.Log("Creating Leaf: " + runArgs.Name)
  os.Mkdir(runArgs.Name, 0755)

  // TODO: Cleaner touch?
  fn, err := os.Create(fmt.Sprintf("%s/main.tf", runArgs.Name))
  if err != nil {
    return err
  }
  defer fn.Close()


  return nil
}
