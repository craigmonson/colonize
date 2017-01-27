package initialize

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/color"

	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/log"
	"github.com/craigmonson/colonize/util"
)

type RunArgs struct {
	AcceptDefaults   bool
	InitEnvironments string
	ConfigBuffer     bytes.Buffer
}

func Run(c *config.Config, l log.Logger, args interface{}) error {
	runArgs := args.(RunArgs)

	err := initTask("Config File", configFileTask, c, l, &runArgs)
	if err != nil {
		return err
	}

	err = initTask("Environments", environmentsTask, c, l, &runArgs)
	if err != nil {
		return err
	}

	err = initTask("Summary", summaryTask, c, l, &runArgs)
	if err != nil {
		return err
	}

	return initTask("Create Artifacts", createArtifactsTask, c, l, &runArgs)
}

func initTask(name string, f func(*config.Config, log.Logger, *RunArgs) error, c *config.Config, l log.Logger, args *RunArgs) error {
	l.LogPretty(util.PadRight(fmt.Sprintf("\nINIT [%s] ", name), "*", 79), color.Bold)
	err := f(c, l, args)
	return err
}

func configFileTask(c *config.Config, l log.Logger, args *RunArgs) error {
	if !args.AcceptDefaults {
		l.Log("Starting interactive setup...\n")
	} else {
		l.LogPretty("Skipping interactive setup...", color.FgBlue)
	}

	scan := bufio.NewScanner(os.Stdin)
	defaults := reflect.ValueOf(&config.ConfigFileDefaults).Elem()
	actuals := reflect.ValueOf(&c.ConfigFile).Elem()
	typeof := defaults.Type()

	// Iterate through each field in the config file struct
	// asking for user inout or setting to a default
	for i := 0; i < defaults.NumField(); i++ {

		fieldName := strings.ToLower(typeof.Field(i).Name)
		defaultField := defaults.Field(i)
		actualField := actuals.Field(i)

		if args.AcceptDefaults || strings.HasPrefix(fieldName, "autogenerate_") {
			actualField.Set(reflect.Value(defaultField))
		} else {
			l.PrintPretty(fmt.Sprintf("Enter '%s' [%v]: ", fieldName, defaultField.Interface()), color.Bold)
			scan.Scan()
			val := scan.Text()

			if val == "" {
				actualField.Set(reflect.Value(defaultField))
			} else {
				actualField.SetString(val)
			}
		}

		args.ConfigBuffer.WriteString(fmt.Sprintf("%-30s => %v\n", fieldName, actualField.Interface()))
	}

	return nil
}

func environmentsTask(c *config.Config, l log.Logger, args *RunArgs) error {

	scan := bufio.NewScanner(os.Stdin)

	if args.InitEnvironments != "" {
		l.LogPretty("Skipping interactive environments setup...", color.FgBlue)
		l.Log("Environments configured from command line: " + args.InitEnvironments)
	} else {

		if !args.AcceptDefaults {
			l.LogPretty(util.PadRight("\nInitialize [Environments] ", "*", 79), color.Bold)
			l.Log("Starting environment setup\n")
			l.PrintPretty("Provide a comma-separated list of environment names to initialize []: ", color.Bold)
			scan.Scan()
			args.InitEnvironments = scan.Text()
		}

		if args.InitEnvironments == "" {
			l.LogPretty("No environments will be setup...", color.FgBlue)
		}
	}

	args.ConfigBuffer.WriteString(fmt.Sprintf("%-30s => %s\n", "Environments", args.InitEnvironments))

	return nil
}

func summaryTask(c *config.Config, l log.Logger, args *RunArgs) error {

	scan := bufio.NewScanner(os.Stdin)

	// print configuration details
	l.Log(fmt.Sprintf("\nInitializing Colonize using the following configuration:\n%s", args.ConfigBuffer.String()))

	if !args.AcceptDefaults {
		// Confirm initialization
		accept := ""
		for accept == "" {
			l.PrintPretty("Enter 'yes' to accept this configuration: ", color.Bold)
			scan.Scan()
			accept = scan.Text()
		}

		if accept != "yes" {
			return errors.New("Colonize initialization cancelled by user")
		}
	}

	return nil
}

func createArtifactsTask(c *config.Config, l log.Logger, args *RunArgs) error {

	l.Log("Creating configuration file...")
	err := c.ConfigFile.WriteToFile(".colonize.yaml")
	if err != nil {
		return errors.New("Failed to create .colonize.yaml file: " + err.Error())
	}

	l.Log("Creating environments directory...")
	err = os.Mkdir(c.ConfigFile.Environments_Dir, 0755)
	if err != nil {
		os.Remove(".colonize.yaml")
		return errors.New("Failed to create environments directory: " + err.Error())
	}

	l.Log("Creating branch order file...")
	util.Touch(c.ConfigFile.Branch_Order_File)

	if len(args.InitEnvironments) > 0 {
		envs := strings.Split(args.InitEnvironments, ",")
		for _, env := range envs {
			fn := env + ".tfvars"
			l.Log(fmt.Sprintf("Creating %s environment variable file...", env))
			util.Touch(c.ConfigFile.Environments_Dir, fn)
		}
	}

	return nil
}
