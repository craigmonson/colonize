package cmd

import (
  "os"

  "github.com/spf13/cobra"
  "github.com/craigmonson/colonize/config"
  "github.com/craigmonson/colonize/initialize"
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
data around that common idiom.

Starting interactive setup:
`

const Footer string = "\n\nColoinze project initialization complete!"
var acceptDefaults bool

var initCmd = &cobra.Command{
  Use: "init",
  Short: "Initialize a Colonize project",
  Long: `This command is used to aid in the generation the .colonize.yaml configuration file and project directory structure.`,
  Run: func(cmd *cobra.Command, args []string) {

    Log.Log(Header)

    _, err := GetConfig(false)
    if err == nil {
      Log.Log("Colonize project already initialized. Exiting.")
      os.Exit(0)
    }

    var blankConfig config.ColonizeConfig
    err = initialize.Run(&blankConfig, Log, acceptDefaults)
    if err != nil {
      Log.Log("ERROR: " + err.Error())
      os.Exit(-1)
    }

    err = blankConfig.ConfigFile.WriteToFile(".colonize.yaml")
    if err != nil {
      Log.Log("ERROR: Failed to create configuration file: " + err.Error())
      os.Exit(-1)
    }
    Log.Log("\nConfiguration file saved...")

    err = os.Mkdir(blankConfig.ConfigFile.Environments_Dir,0755)
    if err != nil {
      Log.Log("ERROR: Failed to crete environments directory: " + err.Error())
      os.Exit(-1)
    }
    Log.Log("Environments directory created...")

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
  RootCmd.AddCommand(initCmd)
}
