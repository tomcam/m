package app

import (
	"embed"
	//"fmt"
	"github.com/tomcam/m/pkg/default"
	"io/fs"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// The following embeds all files and subdirectories
// from the themes subdirectory of this package into
// the executable. That subdirectory contains all
// the factory themes--the base set that always
// gets distributed with Metabuz.
// So have the subdirectory available from within this
// subdirectory at compile time. Then you can run the finished
// executable anywhere and it will display the
// list of files even though the themes directory
// doesn't exist at runtime.

//go:embed themes/*
// TODO: renamed from themeFiles to factoryThemeFiles
var factoryThemeFiles embed.FS

// Todo: changed name from embedDirCpy() to copyFactoryThemes
// copyFactoryThemes() copies the theme files embedded in
// this subdirectory to the project's themes directory.
// In turn, when the site is published only the themes
// it needs will be copied over.
func (app *App) copyFactoryThemes() error {
	var target string
	fs.WalkDir(factoryThemeFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// TODO: Improve error handling
			return err
		}
		// path is the relative path of the file, for example,
		// it might be /en/products or something like that
		if d.IsDir() {
			if path == "." {
				return nil
			}
			// Todo: compute this using config info
			//target = filepath.Join(app.site.path, path)
			target = filepath.Join(app.cfgPath, path)
			err := os.MkdirAll(target, defaults.PublicFilePermissions)
			if err != nil {
				// TODO: Improve error handling
				app.Note("\tos.MkdirAll() error: %v", err.Error())
				return err
			}
			return nil
		}
		//app.Note("\t\t%v", path)
		target = filepath.Join(app.site.factoryThemesPath, path)
		f, err := factoryThemeFiles.Open(path)
		if err != nil {
			// TODO: Improve error handling
			app.Note("\tFS.Open(%v) error: %v", path, err.Error())
			return err
		}
		b, err := io.ReadAll(f)
		if err != nil {
			// TODO: Improve error handling
			app.Note("\tio.ReadAll(%v) error: %v", f, err.Error())
			return err
		}
    err = ioutil.WriteFile(target, b,defaults.ProjectFilePermissions)
    if err != nil {
			// TODO: Improve error handling
			app.Note("\tio.WriteFile(%v) error: %v", f, err.Error())
			return err
    }
		return nil
	})
	return nil
}

// loadTheme() finds the theme specified for this page.
// TODO:
// Correct operation is:
// - Look for theme named in front matter
// - If no theme is named in front matter,
//   look for one named in the site file
// TODO: Create docs for the Metabuzz file
// - If no theme is named in the site file,
//   look for one named in the Metabuzz file
// - If no theme is named in the Metabuzz file,
//   use the default theme named in defaults.DefaultThemeName
func (app *App) loadTheme() {
}
