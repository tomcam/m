package app

import (
	"github.com/tomcam/m/pkg/default"
	"github.com/tomcam/m/pkg/util"
	"os"
	"path/filepath"
	"strings"
)

// Called by  getProjectTree()
// Builds a list all files and all directories in the project.
// Excludes the assets directory and the publish directory.
func (app *App) visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		// Find out what directories to exclude
		exclude := util.NewSearchInfo(app.excludeDirs())
		if err != nil {
			// Quietly fail if unable to access path.
			return err
		}
		isDir := info.IsDir()

		// Skip any directory to be excluded, such as
		// the pub directory itself and anything
		// the user specified in the siteConfig's "Exclude"
		name := info.Name()

		// Exclude this directory if it starts with "."
		// UNLESS it's the project directory. This allows
		// us to build a project where the name  starts
		// with a dot, but its subdirectories that
		// start with a dot will be skipped.
		if strings.HasPrefix(name, ".") && isDir {
			return filepath.SkipDir
			if currDir() != app.Site.path {
				return filepath.SkipDir
			}
		}

		// Exclude this directory if found on the, ah, exclusion list.
		if exclude.Contains(name) && isDir {
			app.Verbose("Excluding directory: %s", name)
			return filepath.SkipDir

		}

		if exclude.Contains(name) {
			app.Verbose("Excluding: %s", name)
			return nil
		}

		if isDir {
			app.setMdOption(path, NormalDir)
		}

		*files = append(*files, path)
		return nil
	}

}

// Obtain a list of all files in the specified project tree starting
// at the root.
// Ignore directories starting with a .
// Ignore the assets directory
func (a *App) getProjectTree(path string) (tree []string, err error) {
	// This should only be called once so I imagine the
	// following is unnecessary
	var files []string
	err = filepath.Walk(path, a.visit(&files))
	if err != nil {
		return []string{}, ErrCode("0702", err.Error(), path)
	}
	// fmt.Fprintf(os.Stdout, "Directory tree for %+v\n", files)
	return files, nil
}

// Returns a list of files and directories to be excluded from the source directory when the
// project is built. It's based on internal configuration (for example, it excludes the
// publish directory) and any existing excludes (for example, Exclude=["pub", "node_modules"])
// in the site config file.
func (a *App) excludeDirs() []string {
	// fmt.Println("Excluded in Site.toml:", a.Site.ExcludeDirs)
	// Add the publish directory if it isn't already there.
	return append(a.Site.ExcludeDirs,
		defaults.CommonPath,
		defaults.HeadTagsPath,
		defaults.HeadTagsPath,
		defaults.SCodePath,
		defaults.ThemesDir)
}
