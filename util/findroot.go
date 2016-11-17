package util

import (
	"errors"
	"os"
	"path"
)

const RootFile = ".colonize.yaml"

func FindRoot(myPath string) (string, error) {
	// cleanup and combine rootfile and search path
	cleanedP := path.Clean(myPath)
	fSearch := path.Join(cleanedP, RootFile)

	// if the file exists, return the full path to the root file.
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
	return FindRoot(newP)
}

// Search for a file.
func fileExists(p string) bool {
	if _, err := os.Stat(p); err == nil {
		return true
	}

	return false
}
