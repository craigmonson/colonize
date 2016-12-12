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
	err := BuildCombinedValsFile(c)
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

	err = BuildCombinedDerivedFiles(c)
	if err != nil {
		return err
	}

	err = BuildRemoteFile(c)
	if err != nil {
		return err
	}

	return nil
}

func BuildCombinedValsFile(c *config.ColonizeConfig) error {
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

func BuildCombinedDerivedFiles(c *config.ColonizeConfig) error {
	combined, err := combineFiles(c.WalkableDerivedPaths)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(c.CombinedValsFilePath)
	if err != nil {
		return err
	}
	substituted := subDerivedWithVariables(content, combined)
	err = writeCombinedFile(c.CombinedDerivedValsFilePath, substituted)
	if err != nil {
		return err
	}
	derVars := getDerivedAsVariables(c)
	return writeCombinedFile(c.CombinedDerivedVarsFilePath, derVars)
}

func BuildRemoteFile(c *config.ColonizeConfig) error {
	valFile, err := getOneValsFile(c)
	if err != nil {
		return err
	}
	remote, err := ioutil.ReadFile(c.RemoteFilePath)
	if err != nil {
		return err
	}
	substituted := subDerivedWithVariables(valFile, remote)
	return writeCombinedFile(c.CombinedRemoteFilePath, substituted)
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
	return getListAsVariables(getConfDerivedVarList(c))
}

func getConfigAsValues(c *config.ColonizeConfig) []byte {
	output := ""
	for _, ary := range getConfDerivedVarList(c) {
		// foo = "bar"
		output = output + ary[0] + " = \"" + ary[1] + "\"\n"
	}

	return []byte(output)
}

func getDerivedAsVariables(c *config.ColonizeConfig) []byte {
	combined, err := combineFiles(c.WalkableDerivedPaths)
	if err != nil {
		panic(err)
	}
	derList := getDerivedVarList(combined)
	return getListAsVariables(derList)
}

func getListAsVariables(varList [][2]string) []byte {
	output := ""
	for _, ary := range varList {
		// variable "foo" {}
		output = output + "variable \"" + ary[0] + "\" {}\n"
	}

	return []byte(output)
}

func getDerivedVarList(content []byte) [][2]string {
	output := [][2]string{}
	for k, v := range getVariableMap(content) {
		output = append(output, [2]string{k, v})
	}
	return output
}

// we're returning a slice because we want the lists to stay in order when
// we print it out... with a map, there's no guarantee of order.  There's no
// functional need to do this for terraform et al... just trying to keep it
// nice for the user.
func getConfDerivedVarList(c *config.ColonizeConfig) [][2]string {
	return [][2]string{
		[2]string{"environment", c.Environment},
		[2]string{"origin_path", c.OriginPath},
		[2]string{"tmpl_name", c.TmplName},
		[2]string{"tmpl_path_dashed", strings.Replace(c.TmplName, "/", "-", -1)},
		[2]string{"tmpl_path_underscored", strings.Replace(c.TmplName, "/", "_", -1)},
		[2]string{"root_path", c.RootPath},
	}
}

func subDerivedWithVariables(content, derived []byte) []byte {
	for k, v := range getVariableMap(content) {
		derived = bytes.Replace(derived, []byte("${var."+k+"}"), []byte(v), -1)
	}

	return derived
}

func getOneValsFile(c *config.ColonizeConfig) ([]byte, error) {
	combined, err := ioutil.ReadFile(c.CombinedValsFilePath)
	if err != nil {
		return nil, err
	}
	derived, err := ioutil.ReadFile(c.CombinedDerivedValsFilePath)
	if err != nil {
		return nil, err
	}

	return append(combined, derived...), nil
}

func getVariableMap(content []byte) map[string]string {
	varMap := map[string]string{}
	for _, line := range strings.Split(string(content), "\n") {
		// skip if the line doesn't match blah = "blah"
		if matched, _ := regexp.MatchString("^.*=.*\".*\"$", line); !matched {
			continue
		}
		KV := strings.Split(line, "=")
		// if, for some reason it split to more or less than 2, skip it
		if len(KV) != 2 {
			continue
		}
		// clean it up...
		key := strings.TrimSpace(KV[0])
		val := strings.TrimSpace(KV[1])
		// remove first and last double quotes from the value
		val = strings.TrimSuffix(strings.TrimPrefix(val, `"`), `"`)
		varMap[key] = val
	}

	return varMap
}
