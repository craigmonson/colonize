package prep

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/craigmonson/colonize/config"
	"github.com/craigmonson/colonize/util"
	//"github.com/craigmonson/colonize/variables"
)

func Run(c *config.ColonizeConfig) error {
	err := BuildCombinedValuesFile(c)
	if err != nil {
		return err
	}

	err = BuildCombinedVarsFile(c)
	if err != nil {
		return err
	}

	err = BuildCombinedTfFile(c)
	if err != nil {
		return err
	}

	return nil
}

func BuildCombinedValuesFile(c *config.ColonizeConfig) error {
	combined, err := combineFiles(c.WalkableValPaths)
	if err != nil {
		return err
	}
	return writeCombinedFile(c.CombinedValsFilePath, combined)
}

func BuildCombinedVarsFile(c *config.ColonizeConfig) error {
	combined, err := combineFiles(c.WalkableVarPaths)
	if err != nil {
		return err
	}
	return writeCombinedFile(c.CombinedVarsFilePath, combined)
}

func BuildCombinedTfFile(c *config.ColonizeConfig) error {
	// get list of files to combine  they can be any tf files
	tfFiles := findTfFilesToCombine(c.WalkableTfPaths)

	combined, err := combineFiles(tfFiles)
	if err != nil {
		return err
	}
	return writeCombinedFile(c.CombinedTfFilePath, combined)
}

func findTfFilesToCombine(dirPaths []string) []string {
	combineable := []string{}
	for _, path := range dirPaths {
		fList, _ := ioutil.ReadDir(path)
		for _, fPath := range fList {
			fullPath := util.PathJoin(path, fPath.Name())
			if isValidTfFile(fullPath) {
				fmt.Println("IS VALID: " + fullPath)
				combineable = append(combineable, fullPath)
			}
		}
	}

	return combineable
}

func isValidTfFile(path string) bool {
	// skip variable.tf, and files that don't end in '.tf'
	vFile, _ := regexp.MatchString("/variables\\.tf$", path)
	tfFile, _ := regexp.MatchString("\\.tf$", path)
	return (!vFile) && tfFile
}

func combineFiles(paths []string) ([]byte, error) {
	combined := []byte{}
	for _, path := range paths {
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			if os.IsPermission(err) {
				return nil, err
			}

			continue
		}
		combined = append(combined, contents...)
	}
	return combined, nil
}

func writeCombinedFile(path string, content []byte) error {
	return ioutil.WriteFile(path, content, 0666)
}
