package app

import (
	"bufio"
	"fmt"
	"github.com/plus3it/gorecurcopy"
	"github.com/spf13/viper"
	"github.com/tomcam/m/pkg/default"
	"github.com/tomcam/m/pkg/util"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// cfgBool() obtains a value set from a config file, environment
// variable, whatever. Simple abstraction over viper
func cfgBool(option string) bool {
	return viper.GetBool(option)
}

// Copy() does just that. It copies a single file named source to
// the file named in dest. There's a lot of erro checking,
// for example, i won't copy a file onto itself.
// Not that that happened to me, of course.
func Copy(src, dest string) error {
	if src == dest {
		// Someone done screwed up
		return ErrCode("0221", src)
	}
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

// exists() is a helper utility that simply displays a filename and
// shows if it's actually present
func exists(description, filename string) string {
	found := false
	if isDirectory(filename) {
		found = true
	}
	r := fmt.Sprint(description, " ", filename)
	if fileExists(filename) {
		found = true
	}

	if found {
		r = r + ": (present)"
	} else {
		r = r + ": (Not present)"
	}
	return r
}

// firstN() returns first part of string.
// Gratefully stolen from https://stackoverflow.com/a/41604514/478311
func firstN(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}
	return s
}

func Oldexists(description, filename string) {
	found := false
	if isDirectory(filename) {
		found = true
	}
	fmt.Print(description, " ", filename)
	if fileExists(filename) {
		found = true
	}

	if found {
		fmt.Println(": (present)")
	} else {
		fmt.Println(": (Not present)")
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

// isDirectory() returns true if the specified
// path name is a directory.
func isDirectory(pathName string) bool {
	f, err := os.Stat(pathName)
	if err != nil {
		return false
	}
	return f.Mode().IsDir()
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
	return ("<meta name=" + quote + tag + quote + " content=" + quote + content + quote + ">\n")
}

// promptString() displays a prompt, then awaits for keyboard
// input and returns it on completion.
// See also inputString(), promptYes()
func promptString(format string, ss ...interface{}) string {
	fmt.Print(fmtMsg(format, ss...))
	fmt.Print(" ")
	return inputString()
}

//	fmt.Println("Warning: " + fmtMsg(format, ss...))

// promptStringDefault() displays a prompt, then awaits for keyboard
// input and returns it on completion. It precedes the end of the
// prompt with a default value in brackets.
// See also inputString(), promptYes()
func promptStringDefault(prompt string, defaultValue string) string {
	answer := promptString(prompt + " [" + defaultValue + "] ")
	if answer == "" {
		return defaultValue
	} else {
		return answer
	}
	// TODO Remove if not used
	fmt.Print(prompt + " [" + defaultValue + "] ")
	answer = inputString()
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
	for {
		answer := promptString(prompt)
		if strings.HasPrefix(strings.ToLower(answer), "y") ||
			strings.HasPrefix(strings.ToLower(answer), "n") {
			return strings.HasPrefix(strings.ToLower(answer), "y")
		}
	}
	///return strings.HasPrefix(strings.ToLower(answer), "y")
}

/*
// TODO: Kill, it's a failure
func readStarterConfig(filename string, s *StarterConfig) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	return nil
}
*/

func readYAMLFile(filename string, target interface{}) (err error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, &target)
	if err != nil {
		return err
	}
	//fmt.Printf("readYAMLFile(): %#v", string(b))
	return nil
}

// replaceExtension() is passed a filename and returns a filename
// with the specified extension.
func replaceExtension(filename string, newExtension string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename)) + "." + newExtension

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

// TODO: may be outdated
// writeYamlFile() creates a YAML file based on the filename and
// data structure passed in.
func writeYamlFile(filename string, target interface{}) error {
	theYaml, err := yaml.Marshal(&target)
	if err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	// TODO: TRY TO REUSE ERROR CODES
	return ioutil.WriteFile(filename, theYaml, defaults.ProjectFilePermissions)
}

func WriteStructToYAML(filename string, i interface{}) error {
	b, err := yaml.Marshal(i)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func readYAMLToStruct(filename string, i interface{}) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	return nil
}

// CopyDirectory() does a recursive deep copy 
// unless keepDirs is true,
// in which case it copies only 1 level of directory.
// Slightly modified from:
// https://stackoverflow.com/questions/51779243/copy-a-folder-in-go
func CopyDirectory(scrDir, dest string, keepDirs bool) error {
	entries, err := ioutil.ReadDir(scrDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if !keepDirs {
				return nil
			}
			if keepDirs {
				if err := CreateIfNotExists(destPath, 0755); err != nil {
					return err
				}
				if err := CopyDirectory(sourcePath, destPath, true); err != nil {
					return err
				}
			}
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := Copy(sourcePath, destPath); err != nil {
				return err
			}
		}

		if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}

		isSymlink := entry.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, entry.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}
