package app

import (
	"bufio"
	"fmt"
	"github.com/plus3it/gorecurcopy"
	"github.com/spf13/viper"
	"github.com/tomcam/m/pkg/default"
	"github.com/tomcam/m/pkg/util"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// cfgBool() obtains a value set from a config file, environment
// variable, whatever. Simple abstraction over viper
func cfgBool(option string) bool {
	return viper.GetBool(option)
}

// Copy() does just that. It copies a single file named source to
// the file named in dest.
func Copy(src, dest string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		// TODO: document error code and add to errors.go
		return ErrCode("0112", src)
	}

	if !sourceFileStat.Mode().IsRegular() {
		// TODO: document error code and add to errors.go
		return ErrCode("0113", src)
	}
	source, err := os.Open(src)
	if err != nil {
		// TODO: document error code and add to errors.go
		return ErrCode("0114", src)
	}
	destination, err := os.Create(dest)
	if err != nil {
		// TODO: document error code and add to errors.go
		return ErrCode("0209", dest)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		// TODO: document error code and add to errors.go
		return ErrCode("0251", dest)
	}
	// Success
	return nil
}

func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}

// curDir() returns the current directory name.
func currDir() string {
	//if path, err := os.Executable(); err != nil {
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

// copyDirOnly() copies a directory nonrecursively.
// Doesn't other directories
// Thanks to https://github.com/plus3it/gorecurcopy/blob/master/gorecurcopy.go
func copyDirOnly(source, dest string) error {
	err := os.MkdirAll(dest, defaults.PublicFilePermissions)
	if err != nil {
		return ErrCode("0410", err.Error(), dest)
	}
	entries, err := ioutil.ReadDir(source)
	if err != nil {
		// TODO: document error code and add to errors.go
		return ErrCode("0708", err.Error(), source)
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		destPath := filepath.Join(dest, entry.Name())
		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			// TODO: document error code and add to errors.go
			return ErrCode("0129", err.Error(), sourcePath)
		}
		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			// Do nothing. What's the syntax for that?
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
				// TODO: document error code and add to errors.go
				return ErrCode("0130", err.Error(), sourcePath)
			}
		default:
			if err := Copy(sourcePath, destPath); err != nil {
				// TODO: document error code and add to errors.go
				return ErrCode("0214", err.Error(), sourcePath)
			}
		}
	}
	return nil
}

// copyDirAll() does a recursive copy of the directory and its subdirectories
func copyDirAll(source, dest string) error {
	if source == "" {
		// TODO: document error code and add to errors.go
		return ErrCode("0704", source)
	}
	if dest == "" {
		// TODO: document error code and add to errors.go
		return ErrCode("0705", dest)
	}

	if dest == source {
		// TODO: document error code and add to errors.go
		return ErrCode("0707", "from '"+source+"' to '"+dest+"'")
	}

	err := gorecurcopy.CopyDirectory(source, dest)
	if err != nil {
		// TODO: document error code and add to errors.go
		return ErrCode("0406", "from '"+source+"' to '"+dest+"'", "")
	}
	return nil
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

// fileToBuf() reads the named file into a byte slice and returns
// that byte slice. In the spirit of HTML it simply returns an empty
// slice on failure.
func fileToBuf(filename string) []byte {
	var input []byte
	var err error
	// Read the whole file into memory as a byte slice.
	input, err = ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}
	}
	return input
}



// fileToString() sucks up a file and returns its contents as a string.
// Fails quietly  if unable to open the file, since
// we're just generating HTML.
func fileToString(infile string) string {
	input, err := ioutil.ReadFile(infile)
	if err != nil {
		return ""
	}
	return string(input)
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
	// The above code could just be replaced with this I guess.
	return isSiteFilePath(path)

}

// isSiteFilePath() (formerly isCfgPath()) looks for the special name used for the subdirectory
// used to hold site config file & info
// formerly isSitePath
func isSiteFilePath(path string) bool {
	siteFile := filepath.Join(path, defaults.CfgDir, defaults.SiteConfigFilename)
	return fileExists(siteFile)
}

// metatag() generates a meta tag. It's complicated.
func metatag(tag string, content string) string {
	const quote = `"`
	return ("\n<meta name=" + quote + tag + quote + " content=" + quote + content + quote + ">\n")
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
	return filepath.Join(path, defaults.CfgDir)
}

// userConfigPath() returns the location used to store
// application data for this user. For example, this is
// includes the subdirectory where themes get copied when
// when Metabuzz is installed. To distinghuish it from
// other products in the same directory, append the
// short name of the product with a dot in front of
// it, e.g. ".mb"
// See:
// https://pkg.go.dev/os#UserConfigDir
func userConfigPath() string {
	// Try to get the official path for application config.
	var path string
	var err error
	if path, err = os.UserConfigDir(); err != nil {
		// On failure, just return the user's home directory.
		path = homeDir()
	}
	return filepath.Join(path, defaults.CfgDir)
}

// wrapTag() returns the HTML code contents surrounded
// by the tag specified, which might look like
// '<nav>' or it might look like '<article id="article">'
// If block is true then it simply adds a newline for
// clarity.
func wrapTag(tag string, contents string, block bool) string {
	var newline string
	if block {
		newline = "\n"
	}
	var endTag, output string
	if len(tag) > 3 {
		output = newline + tag + contents + tag[:1] + "/"
		if strings.Contains(tag, "id=") {
			endTag = strings.Fields(tag[1:])[0]
			output += endTag + ">"
		} else {
			output += tag[1:]
		}
		return output + newline
	}
	return ""
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
