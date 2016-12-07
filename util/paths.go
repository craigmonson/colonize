package util

import (
	"errors"
	"os"
	"path"
	"strings"
)

const RootFile = ".colonize.yaml"

// search through myPath till you find the config file, and return the path to
// it.
func FindCfgPath(myPath string) (string, error) {
	// cleanup and combine config file and search path
	cleanedP := path.Clean(myPath)
	fSearch := path.Join(cleanedP, RootFile)

	// if the file exists, return the full path to the config file.
	if fileExists(fSearch) {
		return fSearch, nil
	}

	// split it for recurisve search.
	newP, _ := path.Split(cleanedP)

	// if this is the end of the line for the path, return an error
	if newP == "/" || newP == "" || newP == "." || newP == ".." {
		return "", errors.New(RootFile + " not found in the directory tree.")
	}

	// recurse
	return FindCfgPath(newP)
}

// return the basename of a given path
func GetBasename(myPath string) string {
	return path.Base(myPath)
}

func GetTmplRelPath(full string, rootPath string) string {
	return strings.Replace(strings.Replace(full, rootPath, "", 1), "/", "", 1)
}

func GetDir(myPath string) string {
	return path.Dir(myPath)
}

func GetTreePaths(myPath string) []string {
	elems := strings.Split(myPath, "/")
	paths := make([]string, len(elems))
	for i, _ := range elems {
		paths[i] = strings.Join(elems[0:i+1], "/")
	}

	return paths
}

func AppendPathToPaths(ps []string, pName string) []string {
	res := []string{}
	for _, p := range ps {
		res = append(res, path.Join(p, pName))
	}

	return res
}

func PrependPathToPaths(ps []string, pName string) []string {
	res := []string{}
	for _, p := range ps {
		res = append(res, path.Join(pName, p))
	}

	return res
}

func PathJoin(p string, p2 string) string {
	return path.Join(p, p2)
}

func AddFileToWalkablePath(myPath string, fName string) []string {
	return AppendPathToPaths(GetTreePaths(myPath), fName)
}

// Search for a file.
func fileExists(p string) bool {
	if _, err := os.Stat(p); err == nil {
		return true
	}

	return false
}
