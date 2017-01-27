package cmd

import (
	"os"

	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/initialize"
	"github.com/spf13/cobra"
)

const Header string = `
  ______             __                      __                     
 /      \           |  \                    |  \                    
|  $$$$$$\  ______  | $$  ______   _______   \$$ ________   ______  
| $$   \$$ /      \ | $$ /      \ |       \ |  \|        \ /      \ 
| $$      |  $$$$$$\| $$|  $$$$$$\| $$$$$$$\| $$ \$$$$$$$$|  $$$$$$\
| $$   __ | $$  | $$| $$| $$  | $$| $$  | $$| $$  /    $$ | $$    $$
| $$__/  \| $$__/ $$| $$| $$__/ $$| $$  | $$| $$ /  $$$$_ | $$$$$$$$
 \$$    $$ \$$    $$| $$ \$$    $$| $$  | $$| $$|  $$    \ \$$     \
  \$$$$$$   \$$$$$$  \$$  \$$$$$$  \$$   \$$ \$$ \$$$$$$$$  \$$$$$$$
                                                                    
--------------------------------------------------------------------

Colonize is a configurable, albeit opinionated way to organize and 
manage your terraform templates. It revolves around the idea of 
environments, and allows you to organize templates, and template 
data around that common organizational structure.

`

const Footer string = "\n\nColoinze project initialization complete!"

var acceptDefaults bool
var initEnvironments string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a Colonize project",
	Long:  `This command is used to aid in the generation the .colonize.yaml configuration file and project directory structure.`,
	Run: func(cmd *cobra.Command, args []string) {

		Log.Log(Header)

		_, err := GetConfig(false)
		if err == nil {
			Log.Log("Colonize project already initialized. Exiting.")
			os.Exit(0)
		}

		var blankConfig config.Config
		err = initialize.Run(&blankConfig, Log, initialize.RunArgs{
			AcceptDefaults:   acceptDefaults,
			InitEnvironments: initEnvironments,
		})
		if err != nil {
			Log.Log("ERROR: " + err.Error())
			os.Exit(-1)
		}

		Log.Log(Footer)
	},
}

func init() {
	initCmd.Flags().BoolVarP(
		&acceptDefaults,
		"accept-defaults",
		"",
		false,
		"automatically accepts default values, skipping manual setup",
	)
	initCmd.Flags().StringVarP(
		&initEnvironments,
		"environments",
		"",
		"",
		"specify environment files to initialize the project with",
	)
	RootCmd.AddCommand(initCmd)
}
