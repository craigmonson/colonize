package prep

import (
	"bytes"
	//"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

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

	err = BuildCombinedDerivedFile(c)
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
	complete := append(getConfigAsValues(c), combined...)
	return writeCombinedFile(c.CombinedValsFilePath, complete)
}

func BuildCombinedVarsFile(c *config.ColonizeConfig) error {
	combined, err := combineFiles(c.WalkableVarPaths)
	if err != nil {
		return err
	}
	complete := append(getConfigAsVariables(c), combined...)
	return writeCombinedFile(c.CombinedVarsFilePath, complete)
}

func BuildCombinedTfFile(c *config.ColonizeConfig) error {
	// get list of files to combine  they can be any tf files
	tfFiles := findTfFilesToCombine(
		c.WalkableTfPaths,
		c.Variable_Tf_File,
		c.Derived_File,
	)

	combined, err := combineFiles(tfFiles)
	if err != nil {
		return err
	}
	return writeCombinedFile(c.CombinedTfFilePath, combined)
}

func BuildCombinedDerivedFile(c *config.ColonizeConfig) error {
	combined, err := combineFiles(c.WalkableDerivedPaths)
	if err != nil {
		return err
	}
	//complete := append(getDerivedFromConfig(c), combined...)
	substituted := subDerivedWithVariables(c, combined)
	return writeCombinedFile(c.CombinedDerivedFilePath, substituted)
}

func findTfFilesToCombine(dirPaths []string, vFile, dFile string) []string {
	combineable := []string{}
	for _, path := range dirPaths {
		fList, _ := ioutil.ReadDir(path)
		for _, fPath := range fList {
			fullPath := util.PathJoin(path, fPath.Name())
			if isValidTfFile(fullPath, vFile, dFile) {
				combineable = append(combineable, fullPath)
			}
		}
	}

	return combineable
}

func isValidTfFile(path, vFile, dFile string) bool {
	// skip variable.tf, and files that don't end in '.tf'
	isVarFile, _ := regexp.MatchString("/"+vFile+"$", path)
	isDerivedFile, _ := regexp.MatchString("/"+dFile+"$", path)
	isTfFile, _ := regexp.MatchString("\\.tf$", path)
	return !isVarFile && !isDerivedFile && isTfFile
}

func combineFiles(paths []string) ([]byte, error) {
	combined := []byte{}
	for _, path := range paths {
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			if os.IsPermission(err) {
				return nil, err
			}

			// it's ok if it doesn't exist
			continue
		}
		combined = append(combined, contents...)
	}
	return combined, nil
}

func writeCombinedFile(path string, content []byte) error {
	return ioutil.WriteFile(path, content, 0666)
}

func getConfigAsVariables(c *config.ColonizeConfig) []byte {
	output := ""
	for _, ary := range getDerivedVarList(c) {
		// variable "foo" {}"
		output = output + "variable \"" + ary[0] + "\" {}\n"
	}

	return []byte(output)
}

func getConfigAsValues(c *config.ColonizeConfig) []byte {
	output := ""
	for _, ary := range getDerivedVarList(c) {
		// foo = "bar"
		output = output + ary[0] + " = \"" + ary[1] + "\"\n"
	}

	return []byte(output)
}

// we're returning a slice because we want the lists to stay in order when
// we print it out... with a map, there's no guarantee of order.
func getDerivedVarList(c *config.ColonizeConfig) [][2]string {
	return [][2]string{
		[2]string{"environment", c.Environment},
		[2]string{"origin_path", c.OriginPath},
		[2]string{"tmpl_name", c.TmplName},
		[2]string{"tmpl_path_dashed", strings.Replace(c.TmplName, "/", "-", -1)},
		[2]string{"tmpl_path_underscored", strings.Replace(c.TmplName, "/", "_", -1)},
		[2]string{"root_path", c.RootPath},
	}
}

func subDerivedWithVariables(c *config.ColonizeConfig, derived []byte) []byte {
	for k, v := range getVariableMap(c) {
		derived = bytes.Replace(derived, []byte("${var."+k+"}"), []byte(v), -1)
	}

	return derived
}

func getVariableMap(c *config.ColonizeConfig) map[string]string {
	varMap := map[string]string{}
	content, _ := ioutil.ReadFile(c.CombinedValsFilePath)
	for _, line := range strings.Split(string(content), "\n") {
		// skip if the line doesn't match
		if matched, _ := regexp.MatchString("^.*=.*\".*\"$", line); !matched {
			continue
		}
		KV := strings.Split(line, "=")
		// if, for some reason it split to more or less than 2, skip it
		if len(KV) != 2 {
			continue
		}
		key := strings.TrimSpace(KV[0])
		val := strings.TrimSpace(KV[1])
		// remove first and last double quotes
		val = strings.TrimSuffix(strings.TrimPrefix(val, `"`), `"`)
		varMap[key] = val
	}

	return varMap
}
