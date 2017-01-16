package initialize

import (
  "bufio"
  "bytes"
  "errors"
  "os"
  "fmt"
  "reflect"
  "strings"

  "github.com/craigmonson/colonize/config"
  "github.com/craigmonson/colonize/log"
)

func Run(c *config.Config, l log.Logger, acceptDefaults bool) error {

  var configBuffer bytes.Buffer

  scan := bufio.NewScanner(os.Stdin)
  defaults := reflect.ValueOf(&config.ConfigFileDefaults).Elem()
  actuals := reflect.ValueOf(&c.ConfigFile).Elem()
  typeof := defaults.Type()

  // Iterate through each field in the config file struct
  // asking for user inout or setting to a default
  for i := 0; i< defaults.NumField(); i++ {

    fieldName := strings.ToLower(typeof.Field(i).Name)
    defaultField := defaults.Field(i)
    actualField := actuals.Field(i)

    if acceptDefaults || strings.HasPrefix(fieldName,"autogenerate_"){
      actualField.Set(reflect.Value(defaultField))
    } else {
      l.Print(fmt.Sprintf("Enter '%s' [%v]: ", fieldName, defaultField.Interface()))
      scan.Scan()
      val := scan.Text()

      if val == "" {
        actualField.Set(reflect.Value(defaultField))
      } else {
        actualField.SetString(val)
      }
    }

    configBuffer.WriteString(fmt.Sprintf("%-30s => %v\n",fieldName,actualField.Interface()))
  }

  if !acceptDefaults {
    // print configuration details
    l.Log(fmt.Sprintf("\n\nInitializing Colonize using the following config:\n%s\n",configBuffer.String()))

    // Confirm initialization
    accept := ""
    for accept == "" {
      l.Print("Please enter [y] to accept this configuration or [n] to cancel: ")
      scan.Scan()
      accept = scan.Text()
    }

    if accept != "y" {
      return errors.New("Colonize initialization cancelled by user")
    }
  }

  return nil
}
