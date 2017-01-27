package initialize

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/log"
	"github.com/craigmonson/colonize/util"
)

type RunArgs struct {
	AcceptDefaults   bool
	InitEnvironments string
}

func Run(c *config.Config, l log.Logger, args interface{}) error {
	var configBuffer bytes.Buffer
	runArgs := args.(RunArgs)

	scan := bufio.NewScanner(os.Stdin)
	defaults := reflect.ValueOf(&config.ConfigFileDefaults).Elem()
	actuals := reflect.ValueOf(&c.ConfigFile).Elem()
	typeof := defaults.Type()

	if !runArgs.AcceptDefaults {
		l.Log("Starting interactive setup:")
	}

	// Iterate through each field in the config file struct
	// asking for user inout or setting to a default
	for i := 0; i < defaults.NumField(); i++ {

		fieldName := strings.ToLower(typeof.Field(i).Name)
		defaultField := defaults.Field(i)
		actualField := actuals.Field(i)

		if runArgs.AcceptDefaults || strings.HasPrefix(fieldName, "autogenerate_") {
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

		configBuffer.WriteString(fmt.Sprintf("%-30s => %v\n", fieldName, actualField.Interface()))
	}

	// print configuration details
	l.Log(fmt.Sprintf("\nInitializing Colonize using the following config:\n%s\n", configBuffer.String()))

	if !runArgs.AcceptDefaults {
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

		if runArgs.InitEnvironments == "" {
			l.Print("\nProvide a comma-separated list of environment names to initialize []: ")
			scan.Scan()
			runArgs.InitEnvironments = scan.Text()
		}

	}

	err := c.ConfigFile.WriteToFile(".colonize.yaml")
	if err != nil {
		return errors.New("Failed to create .colonize.yaml file: " + err.Error())
	} else {
		l.Log("Created Configuration file")
	}

	err = os.Mkdir(c.ConfigFile.Environments_Dir, 0755)
	if err != nil {
		os.Remove(".colonize.yaml")
		return errors.New("Failed to create environments directory: " + err.Error())
	} else {
		l.Log("Created Environments directory")
	}

	util.Touch(c.ConfigFile.Branch_Order_File)
	l.Log("Created Branch Order File")

	envs := strings.Split(runArgs.InitEnvironments, ",")
	for _, env := range envs {
		fn := env + ".tfvars"
		util.Touch(c.ConfigFile.Environments_Dir, fn)
		l.Log(fmt.Sprintf("Created %s environment variable file", env))
	}

	return nil
}
