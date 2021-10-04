package app

import (
	"github.com/tomcam/m/pkg/default"
	"os"
	"path/filepath"
)

// createDirStructuure() creates the specified site structure
// in the current directory.
func createDirStructure(dirs *[][]string) (err error) {
	// Obtain current directory in a portable way.
	basedir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Build up a directory tree for each row
	// in dirs
	for _, row := range *dirs {
		path := basedir
		for _, subdir := range row {
			// Append the next subdirectory in the path
			// in a portable way
			path = filepath.Join(path, subdir)
		}
		err := os.MkdirAll(path, defaults.PublicFilePermissions)
		if err != nil {
			return err
		}
	}
	return nil
}

// currPath) returns the current directory name.
func currPath() string {
	if path, err := os.Getwd(); err != nil {
		return "unknown directory"
	} else {
		return path
	}
}

// dirExists() returns true if the name passed to it is a directory.
func dirExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	} else {
		return false
	}
}

// fileExists() returns true, well, if the named file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// isProject() looks at the structure of the specified directory
// and tries to determine if there's already a project here.
// It does so by looking for site config subdirectory.
func isProject(path string) bool {
	// If the directory doesn't exist, that's easy.
	if !dirExists(path) {
		return false
	}

	// The directory exists. Does it contain a site directory?
	return isCfgPath(path)

}

// isCfgPath() looks for the special name used for the subdirectory
// used to hold site config file & info
// formerly isSitePath
func isCfgPath(path string) bool {
	return dirExists(cfgPath(path))
}

// cfgPath() returns the expected pathname where
// the site file lives. Example: "/Users/tom/html/foo/.mb"
// formerly SitePath
func cfgPath(path string) string {
	return filepath.Join(path, defaults.CfgDir)
}
