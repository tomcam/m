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
	// Location of source theme files computed at
	// runtime
	path string
	// Name is the name of the theme for this page,
	// e.g. "wide"
	Name          string        `yaml:"Name"`
	Branding      string        `yaml:"Branding"`
	Description   string        `yaml:"Description"`
	Stylesheets   []string      `yaml:"Stylesheets"`
	Nav           layoutElement `yaml:"Nav"`
	Header        layoutElement `yaml:"Header"`
	Article       layoutElement `yaml:"Article"`
	Footer        layoutElement `yaml:"Footer"`
	Sidebar       layoutElement `yaml:"Sidebar"`
	Language      string        `yaml:"Language"` // 'en', 'fr', etc.
	HTMLStartFile htmlFragment  `yaml:"HTML-start-file"`
	HTMLEndFile   htmlFragment  `yaml:"HTML-end-file"`
}

type layoutElement struct {
	// Inline HTML
	HTML string `yaml:"HTML"`

	// Filename specifying HTML or Markdown
	File string `yaml:"File"`
}

type htmlFragment struct {
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
		target = filepath.Join(app.Site.factoryThemesPath, path)
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

// themeName() determines the theme name in
// proper order, from most to least proximate.
// Correct operation is:
// - Look for theme named in front matter
// - If no theme is named in front matter,
//   look for one named in the site file
// - If none is there, use Viper
// - If no theme is specified, use the default theme
func (app *App) themeName() string {
	app.Debug("\t\tthemeName(): Checking front matter")
	theme := app.Page.frontMatterMust("theme")
	// See if anything's in the front matter
	// regarding the theme.
	// TODO: Start accounting for theme in other
	// places, like config files
	// TODO: Getting lazy. Remember to marshal front matter appropriately
	// If no theme specified, use the default theme.
	if theme == "" {
		theme = defaults.DefaultThemeName
	}
	return theme
}

// loadTheme() finds the theme specified for this page.
// TODO:
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

	fullTheme := app.themeName()
	app.Page.FrontMatter.Theme = fullTheme
	// If it's something like debut/gallery, loop
	// around and load from root to branch.
	// That way styles are overridden the way
	// CSS expects.
	themeDirs := strings.Split(fullTheme, "/")
	theme := ""

	// Get directory from which themes will be copied
	source := filepath.Join(app.Site.factoryThemesPath,
		defaults.ThemesDir)
	dest := app.Site.siteThemesPath
	for level := 0; level < len(themeDirs); level++ {
		if level == 0 {
			// Build the deepest directory necessary, e.g.
			// .mb/pub/themes/debut/gallery
			err := os.MkdirAll(filepath.Join(app.Site.siteThemesPath, fullTheme), defaults.PublicFilePermissions)
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
		app.Page.Theme.path = source
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
		// Load theme info into app.Page.theme
		if err := app.loadThemeConfig(dest); err != nil {
			// TODO: Handle error properly & and document error code
			app.QuitError(err)
		}
		app.loadStylesheets()

	}
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
func (app *App) loadStylesheets() {
	// If no style sheets don't waste time here
	if len(app.Page.Theme.Stylesheets) <= 0 {
		return
	}
	// Create the published style sheet directory
	// TODO: Track this to make sure it's not repeated unnecessarily
	err := os.MkdirAll(app.Site.cssPublishPath, defaults.PublicFilePermissions)
	if err != nil {
		return
	}
	app.publishStylesheets()

}

func (app *App) publishStylesheets() {
	// Go through the list of stylesheets for this theme.
	// Copy stylesheets for this theme from the local
	// theme directory to the publish
	// CSS directory for stylesheets.
	// This doesn't handle everything. Some stylesheets,
	// such as "theme-dark.css" and "theme-light.css",
	// don't get copied until publish time because
	// they depend on configuration options.
	for _, stylesheet := range app.Page.Theme.Stylesheets {
		// Check every stylesheet to see if it's
		// a dark theme vs a light Theme. If it
		// is, change to dark if requested.
		file := app.getMode(stylesheet)
		source := filepath.Join(app.Page.Theme.path, file)
		dest := filepath.Join(app.Site.cssPublishPath, file)
		// Keep list of stylesheets that got published
		copied := app.copyMust(source, dest)
		if copied != "" {
			app.Page.stylesheets = append(app.Page.stylesheets, dest)
		}
	}
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
	if stylesheet == "theme-light.css" && strings.ToLower(app.Page.frontMatterMust("Mode")) == "dark" {
		stylesheet = "theme-dark.css"
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
	//app.Note("\tloadThemeConfig(%v)", filename)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}

	// Save the current theme. Force to lowercase because
	// it's  filename
	app.Page.Theme.Name = strings.ToLower(filepath.Base(path))

	err = yaml.Unmarshal(b, &app.Page.Theme)
	if err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}

	return nil

}

func (app *App) cfgStr() string {
	return ""
}
