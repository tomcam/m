package app

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tomcam/m/pkg/default"
	"github.com/tomcam/m/pkg/util"
	"os"
	"path/filepath"
	"strings"
)

// cfgBool() obtains a value set from a config file, environment
// variable, whatever. Simple abstraction over viper
func cfgBool(option string) bool {
	return viper.GetBool(option)
}

// curDir() returns the current directory name.
func currDir() string {
	if path, err := os.Getwd(); err != nil {
		return "unknown directory"
	} else {
		return path
	}
}

// cfgString() obtains a value set from a config file, environment
// variable, whatever. Simple abstraction over viper
func cfgString(option string) string {
	return viper.GetString(option)
}

// createDirStructure() creates the specified site structure
// in the current directory.
// TODO: pass in the current directory to save a little time?
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

// hasExtension() returns true if the string ends in the specified extension
// (case insensitive). Need to supply the period too:
// if hasExtension(filename, ".aside") {
func hasExtension(filename, extension string) bool {
	return filepath.Ext(filename) == extension
}

// hasExtensionFrom() Returns true if the fully qualified filename
// ends in any of the extensions listed in extensions.
func hasExtensionFrom(path string, extensions *util.SearchInfo) bool {
	return extensions.Contains(filepath.Ext(path))
}

// homeDir() returns the user's home directory, or just "." for
// the current directory if it can't be determined through system
// calls.
func homeDir() string {
	var home string
	var err error
	if home, err = os.UserHomeDir(); err != nil {
		return "."
	}
	return home
}

// inputString() gets a string from the keyboard and returns it
// See also promptString()
func inputString() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// isMarkdownFile() returns true of the specified filename has one of
// the extensions used for Markdown files.
func isMarkdownFile(filename string) bool {
	return hasExtensionFrom(filename, defaults.MarkdownExtensions)
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
	return isSiteFilePath(path)

}

// isSiteFilePath() (formerly isCfgPath()) looks for the special name used for the subdirectory
// used to hold site config file & info
// formerly isSitePath
func isSiteFilePath(path string) bool {
	return dirExists(siteFilePath(path))
}

// promptString() displays a prompt, then awaits for keyboard
// input and returns it on completion.
// See also inputString(), promptYes()
func promptString(prompt string) string {
	fmt.Print(prompt + " ")
	return inputString()
}

// promptStringDefault() displays a prompt, then awaits for keyboard
// input and returns it on completion. It precedes the end of the
// prompt with a default value in brackets.
// See also inputString(), promptYes()
func promptStringDefault(prompt string, defaultValue string) string {
	fmt.Print(prompt + " [" + defaultValue + "] ")
	answer := inputString()
	if answer == "" {
		return defaultValue
	} else {
		return answer
	}
}

// promptYes() displays a prompt, then awaits
// keyboard input. If the answer starts with Y,
// returns true. Otherwise, returns false.
// See also inputString(), promptString()
func promptYes(prompt string) bool {
	// See also inputString(), promptYes()
	answer := promptString(prompt)
	return strings.HasPrefix(strings.ToLower(answer), "y")
}

// relDirFile() takes a base directory,
// for example, /users/tom/mysite, and a filename, for
// example, /users/tom/mysite/articles/announce.md,
// and returns the relative directory, which would be
// the directory named /articles in this case.
func relDirFile(baseDir, filename string) string {
	// Begin at the end of the base directory
	// xxx
	start := len(baseDir)
	// Extract the target directory from the
	// input filename
	l := len(filepath.Dir(filename))
	// End at the beginning of the filename
	stop := l - start
	// TODO: Playing with fire?
	if stop < 0 {
		stop = start
	}
	return string(filename[start : start+stop])
}

// replaceExtension() is passed a filename and returns a filename
// with the specified extension.
func replaceExtension(filename string, newExtension string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename)) + "." + newExtension

}

// siteFilePath() returns the expected pathname where
// the site file lives. Example: "/Users/tom/html/foo/.mb"
// formerly SitePath
func siteFilePath(path string) string {
	return filepath.Join(path, defaults.CfgPath)
}

// WriteTextFile creates a file called filename without checking to see if it
// exists, then writes contents to it.
func writeTextFile(filename, contents string) error {
	var out *os.File
	var err error
	if out, err = os.Create(filename); err != nil {
		// TODO: Renumber error code?
		return ErrCode("0204", "Problem creating file %v: %v\n", filename, err.Error())
	}
	if _, err = out.WriteString(contents); err != nil {
		// TODO: Renumber error code?
		return ErrCode("0903", "Problem writing to file %v: %v\n", filename, err.Error())
	}
	return nil
}
