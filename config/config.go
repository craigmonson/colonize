package config

import (
	"io/ioutil"

	"github.com/craigmonson/colonize/util"
	"gopkg.in/yaml.v2"
)

// Find the main config file: .colonize.yaml will ALWAYS be in the root of the
// project.

type ColonizerConfig struct {
	Autogenerate_Comment      string
	Combined_Vals_File        string
	Combined_Vars_File        string
	Combined_Tf_File          string
	Vars_File_Env_Post_String string
	Vals_File_Env_Post_String string
	Templates_Dir             string
	Environments_Dir          string
}

func LoadConfig(path string) (*ColonizerConfig, error) {
	conf := ColonizerConfig{}

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return &conf, err
	}

	err = yaml.Unmarshal(contents, &conf)

	if err != nil {
		return &conf, err
	}

	return &conf, nil
}

func LoadConfigInTree(path string) (*ColonizerConfig, error) {
	configPath, err := util.FindRoot(path)
	if err != nil {
		return &ColonizerConfig{}, err
	}

	return LoadConfig(configPath)
}
