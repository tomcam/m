package app

import (
	//"fmt"
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
	// Location of theme files after they have been
	// copied to the publish directory for themes
	// used by this site.
	publishPath string
	// Location of source theme files computed at
	// runtime
	sourcePath string
	// Name is the name of the theme for this page,
	// e.g. "wide"
	Name        string        `yaml:"Name"`
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
	app.Debug("\tcopyFactoryThemes")
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
				return ErrCode("0409", target)
			}
			app.Debug("\t\tcreated directory %v", target)
			return nil
		}
		// It's a file, not a directory
		app.Debug("\t\tCreated theme directory %v", target)
		// Handle individual file
		//app.Page.Theme.sourcePath
		target = filepath.Join(app.Site.factoryThemesPath, path)
		app.Note("app.Site.factoryThemesPath: %v. path: %v", app.Site.factoryThemesPath, path)
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
			app.Note("\t\tcopyFactoryThemes(): err after WriteFile:  %#v", err)
			// TODO: Handle error properly & and document error code
			//app.Debug("\t\tHEY io.WriteFile(%v) error: %v happened during copyFactoryThemes() %v", f, err.Error(), target)
			return ErrCode("0216", err.Error(), target)
			//app.Note("\t\tcopyFactoryThemes(): supposed to return with error value of %#v. err == nil: %v", err, err == nil)
			//return ErrCode("0216", target, err.Error())
			return err
		}
		return nil
	})
	return nil
}

// themeNameToLower() determines the theme name in
// proper order, from most to least proximate.
// TODO: document the following
// It forces the theme name to lowercase.
// Correct operation is:
// - Look for theme named in front matter
// - If no theme is named in front matter,
//   look for one named in the site file
// - If none is there, use Viper
// - If no theme is specified, use the default theme
func (app *App) themeNameToLower() string {
	app.Debug("\t\tthemeName(): Checking front matter")
	// See if anything's in the front matter
	// regarding the theme.
	theme := app.Page.FrontMatter.Theme
	// TODO: Start accounting for theme in other
	// places, like config files
	if theme == "" {
		if app.Site.Theme != "" {
			theme = app.Site.Theme
			app.Debug("\t\tthemeName(): No theme named in front matter. Trying theme name from Site file: %v", app.Site.Theme)
		}
		if theme == "" {
			app.Debug("\t\tthemeName(): No theme named in front matter. Trying default theme name %v", defaults.DefaultThemeName)
			theme = defaults.DefaultThemeName
		}
	}
	return strings.ToLower(theme)
}

// loadTheme() finds the theme specified for this page.
// TODO:
// TODO: Create docs for the Metabuzz file
// - If no theme is named in the site file,
//   look for one named in the Metabuzz file
// - If no theme is named in the Metabuzz file,
//   use the default theme named in defaults.DefaultThemeName
func (app *App) loadTheme() error {
	// Theme designation could something like:
	//  debut
	//  debut/gallery
	//  debut/gallery/item

	app.Page.FrontMatter.Theme = app.themeNameToLower()
	// This is called fullTheme because the theme designation
	// can be something "debut" or it can go down deper,
	// for example, "debut/gallery/item"
	fullTheme := app.Page.FrontMatter.Theme
	app.Debug("\tloadTheme %v", fullTheme)
	// If it's something like debut/gallery, loop around and load from root to branch.
	// That way styles are overridden the way
	// CSS expects.
	themeDirs := strings.Split(fullTheme, "/")
	theme := ""

	// Get directory from which themes will be copied
	source := filepath.Join(app.Site.factoryThemesPath, defaults.SiteThemesDir)
	dest := app.Site.siteThemesPath
	for level := 0; level < len(themeDirs); level++ {
		theme = themeDirs[level]
		// Get the next level of directory and append
		// to the previous directory
		source = filepath.Join(source, theme)
		app.Page.Theme.sourcePath = source
		dest = filepath.Join(dest, theme)
		// See if this theme has already been
		// published.
		_, ok := app.Site.publishedThemes[dest]
		if !ok {
			// Cant find theme in map so apparently it's
			// not yet published.
			// Create a directory for it.
			err := os.MkdirAll(dest, defaults.PublicFilePermissions)
			if err != nil {
				// TODO: Handle error properly & and document error code
				app.Note("\tos.MkdirAll() error: %v", err.Error())
				return ErrCode("PREVIOUS", err.Error())
			}
			// Theme directory created successfuly
		}
		app.Debug("\tloadTheme: copy from %v to %v", source, dest)

		if err := copyDirAll(source, dest); err != nil {
			return ErrCode("PREVIOUS", err.Error())
		}
		app.Page.Theme.publishPath = dest
		app.Site.publishedThemes[dest] = true
		// Theme directory is known. Use it to load the .yaml file
		// for this theme.
		// Load theme info into app.Page.theme
		if err := app.loadThemeConfig(dest); err != nil {
			return ErrCode("PREVIOUS", err.Error())
		}
		app.loadStylesheets()
	}
	return nil
}

// loadStylesheets() finds the stylesheets
// named in the theme file, which is in
// app.Page.theme after the call to
// loadThemeConfig()
// In the spirit of browsers, which don't stop
// loading the page when a stylesheet can't
// be found, this function doesn't return errors
// TODO: Track when these things are copied
// to avoid redoing this work unnecessarily
func (app *App) loadStylesheets() error {
	// If no style sheets don't waste time here
	if len(app.Page.Theme.Stylesheets) <= 0 {
		return nil
	}
	// Create the published style sheet directory
	// TODO: Track this to make sure it's not repeated unnecessarily
	err := os.MkdirAll(app.Site.cssPublishPath, defaults.PublicFilePermissions)
	if err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	/*
		err = os.MkdirAll(app.Page.Theme.sourcePath, defaults.PublicFilePermissions)
		if err != nil {
			return ErrCode("PREVIOUS", err.Error())
		}
	*/

	return nil
}

// getMode() checks if the stylesheet is dark or light
// and adjusts as needed.
// TODO: It should probably call App.cfg() to
// search other places, like an app config file
func (app *App) getMode(stylesheet string) string {
	// Check every stylesheet to see if it's named
	// "theme-light.css" (light theme is the default).
	// If it is, and if Dark mode
	// has been specified, publish theme-dark.css instead.
	mode := strings.ToLower(app.Page.FrontMatter.Mode)
	if stylesheet == "theme-light.css" && mode == "dark" {
		return "theme-dark.css"
	}
	return stylesheet
}

// loadThemeConfig reads the theme's config file, so
// if the theme is named "debut" that file would be
// named debut.yaml. Write to app.Page.Theme.
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
	app.Debug("\tloadThemeConfig(%v)", filename)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		app.Debug("\tloadThemeConfig() failed to read %v", filename)
		// TODO: Handle error properly & and document error code
		return err
	}

	// Save the current theme. Force to lowercase because
	// it's  filename
	app.Page.Theme.Name = strings.ToLower(app.Page.FrontMatter.Theme)
	app.Page.Theme.publishPath = filepath.Join(app.Site.siteThemesPath, app.Page.Theme.Name)

	err = yaml.Unmarshal(b, &app.Page.Theme)
	if err != nil {
		app.Debug("\tloadThemeConfig() Unable to marshal YAML from %v", filename)
		// TODO: Handle error properly & and document error code
		return err
	}
	return nil

}

func (app *App) cfgStr() string {
	return ""
}

func readThemeConfig(filename string) (*Theme, error) {
	var theme Theme
	err2 := readYAMLFile(filename, &theme)
	if err2 != nil {
		panic(err2)
	}
	// This works.
	// fmt.Printf("theme: %v", theme)
	return &theme, nil
}
