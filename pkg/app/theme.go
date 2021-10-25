package app

import (
	"fmt"
	//"encoding/json"
	"embed"
	"github.com/tomcam/m/pkg/default"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
      // TODO: Handle error properly & and document error code
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
        // TODO: Handle error properly & and document error code
				app.Note("\tos.MkdirAll() error: %v", err.Error())
				return err
			}
			return nil
		}
		//app.Note("\t\t%v", path)
		target = filepath.Join(app.site.factoryThemesPath, path)
		f, err := factoryThemeFiles.Open(path)
		if err != nil {
      // TODO: Handle error properly & and document error code
			app.Note("\tFS.Open(%v) error: %v", path, err.Error())
			return err
		}
		b, err := io.ReadAll(f)
		if err != nil {
      // TODO: Handle error properly & and document error code
			app.Note("\tio.ReadAll(%v) error: %v", f, err.Error())
			return err
		}
		err = ioutil.WriteFile(target, b, defaults.ProjectFilePermissions)
		if err != nil {
      // TODO: Handle error properly & and document error code
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
// - If no theme is specified, use the default theme
// TODO: Create docs for the Metabuzz file
// - If no theme is named in the site file,
//   look for one named in the Metabuzz file
// - If no theme is named in the Metabuzz file,
//   use the default theme named in defaults.DefaultThemeName
func (app *App) loadTheme() {
	// Theme designation could something like:
	//  debut
	//  debut/gallery
	//  debut/gallery/item

	fullTheme := ""
  // See if anything's in the front matter
  // regarding the theme.
  // TODO: Start accounting for theme in other
  // places, like config files
	if app.page.frontMatterRaw["theme"] == nil {
		fullTheme = ""
	} else {
		fullTheme = fmt.Sprint(app.page.frontMatterRaw["theme"])
	}

  // If no theme specified, use the default theme.
	if fullTheme == "" {
		fullTheme = defaults.DefaultThemeName
	}

  // If it's something like debut/gallery, loop
  // around and load from root to branch.
  // That way styles are overridden the way
  // CSS expects.
	themeDirs := strings.Split(fullTheme, "/")
	theme := ""

  // Get directory from which themes will be copied
  source := filepath.Join(app.site.factoryThemesPath,
    defaults.ThemesDir)
  dest := app.site.siteThemesPath
	for level := 0; level < len(themeDirs); level++ {
		if level == 0 {
      // Build the deepest directory necessary, e.g.
      // .mb/pub/themes/debut/gallery
			err := os.MkdirAll(filepath.Join(app.site.siteThemesPath, fullTheme), defaults.PublicFilePermissions)
			if err != nil {
        // TODO: Handle error properly & and document error code
				app.Note("\tos.MkdirAll() error: %v", err.Error())
				return
			}
		}
    theme = themeDirs[level]
    // Get the next level of directory and append 
    // to the previous directory
		source = filepath.Join(source, theme)
		dest = filepath.Join(dest, theme)
		if err := copyDirAll(source, dest); err != nil {
			//return ErrCode("0401", source)
			//app.QuitError("0401", source)
			// msg := fmt.Errorf("Error attempting to create project file %s: %v", projectFile, err.Error()).Error()
			app.QuitError(err)
		}
	}
}

func (app *App) cfgStr() string {
	return ""
}
