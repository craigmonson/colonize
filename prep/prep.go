package prep

import (
	"io/ioutil"
	"os"

	"github.com/craigmonson/colonize/config"
	//"github.com/craigmonson/colonize/variables"
)

func Run(c *config.ColonizeConfig) error {
	// build first pass variables map.
	err := BuildCombinedValuesFile(c)
	if err != nil {
		return err
	}

	// build.
	return nil
}

func BuildCombinedValuesFile(c *config.ColonizeConfig) error {
	combined := []byte{}
	for _, path := range c.WalkableValPaths {
		contents, err := ioutil.ReadFile(path)
		if err != nil && os.IsPermission(err) {
			if os.IsPermission(err) {
				return err
			}

			continue
		}
		combined = append(combined, contents...)
	}

	return ioutil.WriteFile(c.CombinedValsFilePath, combined, 0666)
}
