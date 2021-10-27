package app

import (
	"fmt"
	//"encoding/json"
	"embed"
	"github.com/tomcam/m/pkg/default"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Theme struct {
	Branding    string        `yaml:"Branding"`
	Description string        `yaml:"Description"`
	Stylesheets []string      `yaml:"Stylesheets"`
	Nav         layoutElement `yaml:"Nav"`
	Header      layoutElement `yaml:"Header"`
	Article     layoutElement `yaml:"Article"`
	Footer      layoutElement `yaml:"Footer"`
	Sidebar     layoutElement `yaml:"Sidebar"`
}

type layoutElement struct {
	// Inline HTML
  HTML string `yaml:"HTML"`

	// Filename specifying HTML or Markdown
  File string `yaml:"File"`
}

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
	// TODO: Can this whole thing be replaced with a copyDirAll()?
	// Is there a perf benefit either way?
	var target string
	fs.WalkDir(factoryThemeFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// TODO: Handle error properly & and document error code
			return err
		}
		// Handle subdirectory.
		// path is the relative path of the file, for example,
		// it might be /en/products or something like that
		if d.IsDir() {
			if path == "." {
				return nil
			}
			// Get name of destination directory.
			target = filepath.Join(app.cfgPath, path)
			// Create the destination directory.
			err := os.MkdirAll(target, defaults.PublicFilePermissions)
			if err != nil {
				// TODO: Handle error properly & and document error code
				app.Note("\tos.MkdirAll() error: %v", err.Error())
				return err
			}
			return nil
		}
		//app.Note("\t\t%v", path)
		// Handle individual file
		target = filepath.Join(app.site.factoryThemesPath, path)
		f, err := factoryThemeFiles.Open(path)
		if err != nil {
			// TODO: Handle error properly & and document error code
			app.Note("\tFS.Open(%v) error: %v", path, err.Error())
			return err
		}
		// Read the file into a byte array.
		b, err := io.ReadAll(f)
		if err != nil {
			// TODO: Handle error properly & and document error code
			app.Note("\tio.ReadAll(%v) error: %v", f, err.Error())
			return err
		}
		// Copy the recently read file to its destination
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
	// TODO: Getting lazy. Remember to marshal front matter appropriately
	if app.page.frontMatterRaw["theme"] == nil {
		// If no theme specified, use the default theme.
		fullTheme = defaults.DefaultThemeName
	} else {
		fullTheme = fmt.Sprint(app.page.frontMatterRaw["theme"])
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
		app.page.themePath = source
		dest = filepath.Join(dest, theme)
		if err := copyDirAll(source, dest); err != nil {
			//return ErrCode("0401", source)
			//app.QuitError("0401", source)
			// msg := fmt.Errorf("Error attempting to create project file %s: %v", projectFile, err.Error()).Error()
			// TODO: Handle error properly & and document error code
			app.QuitError(err)
		}
		// Theme directory is known. Use it to load the .yaml file
		// for this theme.
		// Load theme info into app.page.theme
		if err := app.loadThemeConfig(dest); err != nil {
			// TODO: Handle error properly & and document error code
			app.QuitError(err)
		}
    app.Note("THEME: %v\n%v\n\n", theme, app.page.theme)
		app.loadStylesheets()

	}
}

// loadStylesheets() finds the stylesheets
// named in the theme file, which is in
// app.page.theme after the call to
// loadThemeConfig()
// In the spirit of browsers, which don't stop
// loading the page when a stylesheet can't
// be found, this function doesn't return errors
// TODO: Track when these things are copied
// to avoid redoing this work unnecessarily
func (app *App) loadStylesheets() {
	// If no style sheets don't waste time here
	if len(app.page.theme.Stylesheets) <= 0 {
		return
	}
	// Create the published style sheet directory
	// TODO: Track this to make sure it's not repeated unnecessarily
	err := os.MkdirAll(app.site.cssPublishPath, defaults.PublicFilePermissions)

	if err != nil {
		return
	}

	// Go through the list of stylesheets for this theme.
	// Copy stylesheets for this theme from the local
	// theme directory to the publish
	// CSS directory for stylesheets
	for _, stylesheet := range app.page.theme.Stylesheets {
		// Check every stylesheet to see if it's
		// a dark theme vs a light theme. If it
		// is, change to dark if requested.
		file := app.getMode(stylesheet)
		source := filepath.Join(app.page.themePath, file)
		dest := filepath.Join(app.site.cssPublishPath, file)
		// Keep list of stylesheets that got published
		copied := app.copyMust(source, dest)
		if copied != "" {
			app.page.stylesheets = append(app.page.stylesheets, dest)
		}
	}
}

// getMode() checks if the stylesheet is dark or light
// and adjusts as needed.
func (app *App) getMode(stylesheet string) string {
	// Check every stylesheet to see if it's named
	// "theme-light.css". If it is, and if Dark mode
	// has been specified, alter its name to
	// theme-dark.css.
	if stylesheet == "theme-light.css" && app.page.frontMatter.Mode == "dark" {
		stylesheet = "theme-dark.css"
	}
	return stylesheet
}

// loadThemeConfig reads the theme's config file, so
// if the theme is named "debut" that file would be
// named debut.yaml. Write to app.page.theme.
// The path passed in is something
// like /Users/tom/code/m/cmd/mb/foo/.mb/pub/themes/wide
// so all that's needed now is to create the fully
// qualified filename, for example, you pass in
//   Users/tom/mb/foo/.mb/pub/themes/wide/
// and convert it to
//   Users/tom/mb/foo/.mb/pub/themes/wide/wide.yaml
func (app *App) loadThemeConfig(path string) error {
	// Strip off everything but the theme name itself.
	// Add that to the base directory because that's the
	// base of the filename. Then add the rest to the
	// filename, which is the theme name + ".yaml"
	filename := filepath.Join(path, filepath.Base(path)+"."+defaults.ConfigFileDefaultExt)
	//app.Note("\tloadThemeConfig(%v)", filename)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}

	err = yaml.Unmarshal(b, &app.page.theme)
	if err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}

	return nil

}

func (app *App) cfgStr() string {
	return ""
}
