package generate

import (
  "fmt"
  "errors"
  "os"

  "github.com/craigmonson/colonize/config"
  "github.com/craigmonson/colonize/log"
)

func RunLeaf(c *config.ColonizeConfig, l log.Logger, name string) error {

  if _, err := os.Stat(name); err == nil {
    return errors.New("Leaf already exists")
  }
  // TODO: Check that leaf doesn't exist

  l.Log("Creating Leaf: " + name)
  os.Mkdir(name, 0755)

  // TODO: Cleaner touch?
  fn, err := os.Create(fmt.Sprintf("%s/main.tf", name))
  if err != nil {
    return err
  }
  defer fn.Close()


  return nil
}
