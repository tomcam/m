package app

import (
	"github.com/tomcam/m/pkg/default"
	"github.com/tomcam/m/pkg/util"
	"gopkg.in/yaml.v3"
	//"io"
	//"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Theme struct {
	// Themes can be nested, e.g. debut/gallery/item.
	// Each level get its own entry here.
	levels []string

	// TODO: Rename either this or levels
	level string

	// Tracks level of nesting for this theme. So if
	// the theme is specified as debut/gallery/item,
	// debut is 0, gallery is 1, and item is 2.
	nestingLevel int

	// Location of theme files after they have been
	// copied to the publish directory for themes
	// used by this site.
	publishPath string

	// Location of source theme files computed at
	// runtime
	sourcePath string

	// List of all stylesheet after being massaged,
	// for example, to ensure responive.css comes
	// last among other things.
	stylesheetList []string

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
} // type Theme

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
// gets distributed with Metabuzz.
// So have the subdirectory available from within this
// subdirectory at compile time. Then you can run the finished
// executable anywhere and it will display the
// list of files even though the themes directory
// doesn't exist at runtime.

// TODOgo:don'tembed themes/*
//var factoryThemeFiles embed.FS

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
	//app.Debug("\t\tthemeName(): Checking front matter")
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

// copyTheme() Directly a copies all files in a theme from the
// fully qualified directory source to the fully qualified
// directory dest.
func (app *App) copyTheme(source string, dest string) error {
	app.Debug("\t\t\t\tcopyTheme(%v, %v)", source, dest)
	err := os.MkdirAll(dest, defaults.PublicFilePermissions)
	if err != nil {
		// TODO: Handle error properly & and document error code
		app.Debug("\t\t\t\tos.MkdirAll(%v) error: %v", dest, err.Error())
		return ErrCode("PREVIOUS", err.Error())
	}
	// TODO: See if there's a faster Go lib for this
	if err := copyDirAll(source, dest); err != nil {
		app.Debug("\t\t\t\tFAILED copyDirAll(%v, %v)", source, dest)
		return ErrCode("PREVIOUS", err.Error())
	}
	return nil
}

// publishThemeAssets() Reads the source theme directory
// and chooses which non-stylesheet file get copied
// to the published theme directory.
// from is a fully qualified directory name for the theme to be copied.
// to is a fully qualified directory name for it to be copied to
func (app *App) publishThemeAssets(from string, to string) error {
	app.Debug("\t\t\t\tpublishTheme(%v, %v)", from, to)
	// Get the directory listing.
	candidates, err := ioutil.ReadDir(from)
	if err != nil {
		return (ErrCode("1033", from, err.Error()))
	}
	excludeFromDir := util.NewSearchInfo(app.Page.FrontMatter.ExcludeFiles)
	//jlkapp.Note("\t\t\t\tpublishThemeAssets: exclude files %v",app.Page.FrontMatter.ExcludeFiles)
	//app.Note("\t\t\t\tpublishThemeAssets: FrontMatter: %#v",app.Page.FrontMatter)
	for _, file := range candidates {
		if file.IsDir() {
			break
		}
		filename := file.Name()
		if !hasExtension(filename, ".css") &&
			!hasExtensionFrom(filename, defaults.ExcludedAssetExtensions) &&
			!excludeFromDir.Contains(filename) &&
			!hasExtensionFrom(filename, defaults.MarkdownExtensions) {
			copyFrom := filepath.Join(from, filename)
			// TODO: should probably go to the page publish directory
			copyTo := filepath.Join(to, filename)
			app.Debug("\t\t\t\t\tCopy(%v,%v)", copyFrom, copyTo)
			if err := Copy(copyFrom, copyTo); err != nil {
				// TODO. This should actually lok something like
				//return ErrCode("1234", "PREVIOUS', copyFrom, err.Error())
				return ErrCode("PREVIOUS", err.Error())
			}
		}
	}
	return nil

} // publishThemeAssets()

// loadThemeLevel() finds the theme specified for this page.
// It then copies the required files to the theme publish
// directory.
//  A Theme designation can something like:
//  debut
//  debut/gallery
//  debut/gallery/item
// This get gets called once for each level. In the case
// of "debut/gallery" the stylesheets for "debut.yaml" are
// loaded, then those in "debut/gallery.yaml" are loaded
// after it, and so on.
// This gives a form of inheritance because anyting in
// gallery.yaml that's identical to something in debut.yaml
// overrides whatever is in debut.yaml.
// Called from loadTheme() once per level.
func (app *App) loadThemeLevel(source string, dest string, level int) error {
	app.Debug("\t\t\tloadThemeLevel(%v, %v, %v)", source, dest, level)
	// See if this theme has already been published.
	// TODO: cache themes in app.Site.publishedThemes[dest]?
	if !dirExists(source) {
		return ErrCode("1028", source)
	}

	//app.Page.Theme.publishPath = dest
	if err := app.loadThemeConfig(source); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	// Theme directory is known. Load its config
	// (e.g. .yaml) file
	err := os.MkdirAll(dest, defaults.PublicFilePermissions)
	if err != nil {
		// TODO: Handle error properly & and document error code
		app.Debug("\t\t\t\tos.MkdirAll(%v) error: %v", dest, err.Error())
		return ErrCode("PREVIOUS", err.Error())
	}

	//err = app.publishThemeAssets(source, dest)
	err = app.publishThemeAssets(source, app.Site.publishPath)
	// TODO: May want to improve error handling
	if err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	app.Page.themes = append(app.Page.themes, app.Page.Theme)
	return nil
} // loadThemeLevel()

// loadTheme() finds the theme specified for this page.
// Load the theme and all its descendants, because
// a theme could be as sample as "debut" or it could be
// gallery and its descendants, e.g. "debut/gallery/item".
//
// TODO:
// - Create docs for the Metabuzz file
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
	fullTheme := strings.ToLower(app.Page.FrontMatter.Theme)
	app.Debug("\t\tloadTheme %v", fullTheme)
	// If it's something like debut/gallery, loop around and load from root to branch.
	// That way styles are overridden the way
	// CSS expects.
	// They've all been forced to lowercase, so "Debut/gallery/image"
	// becomes "debut/gallery/image"
	app.Page.Theme.levels = strings.Split(fullTheme, "/")
	//theme := ""

	// Get directory from which themes will be copied
	from := filepath.Join(app.Site.factoryThemesPath, defaults.SiteThemesDir)

	// Get directory to which the theme will be copied for this site
	//to := app.Site.publishPath
	themeName := ""
	//app.Print("\t\t\tfrom: %v. to: %v", from, to)
	for level := 0; level < len(app.Page.Theme.levels); level++ {
		// Build up each level of nested them: "debut",
		// "debut/gallery", "debut/gallery/item"
		themeName = filepath.Join(themeName, app.Page.Theme.levels[level])
		// Get the next level of directory and append
		// to the previous directory
		source := filepath.Join(from, themeName)
		dest := app.themePublishDir(themeName)
		app.Page.Theme.sourcePath = source
		// TODO: I can probably remove nestingLevel entirely, unless i need it
		// to detect latest sidebar or mode
		app.Page.Theme.nestingLevel = level
		app.Page.Theme.level = themeName
		// xxx
		// Finds the theme specified for this page.
		// Copy the required files to the theme publish directory.
		app.Debug("\t\t\tloading theme(%v,%v,%v", source, dest, level)
		if err := app.loadThemeLevel(source, dest, level); err != nil {
			return ErrCode("PREVIOUS", err.Error())
		}
	}
	return nil
} //loadTheme()

// themePublishDir() returns the name of the directory
// needed to publish style sheets for the theme
func (app *App) themePublishDir(theme string) string {
	//return filepath.Join(app.Site.cssPublishPath, defaults.ThemesDir, app.Page.FrontMatter.Theme)
	//return filepath.Join(app.Site.publishPath, defaults.ThemesDir, app.Page.FrontMatter.Theme)
	return filepath.Join(app.Site.publishPath, defaults.ThemesDir, theme)
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
// named debut.yaml. Write to app.Page.Theme
// after forcing to lower case.
// The path passed in is something
// like /foo/.mb/pub/themes/wide
// so all that's needed now is to create the fully
// qualified filename. For example, you pass in:
//   /foo/.mb/pub/themes/wide/
// and it gets converted to:
//   Users/tom/mb/foo/.mb/pub/themes/wide/wide.yaml
func (app *App) loadThemeConfig(path string) error {

	// Add YAML or whatver to the base directory because that's the
	// base of the filename. Then add the rest to the
	// filename, which is the theme name + ".yaml"
	filename := filepath.Join(path, filepath.Base(path)+"."+defaults.ConfigFileDefaultExt)
	app.Debug("\t\t\t\tloadThemeConfig(%v)", filename)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		app.Debug("\t\t\t\tloadThemeConfig() failed to read %v", filename)
		// TODO: Handle error properly & and document error code
		return err
	}

	err = yaml.Unmarshal(b, &app.Page.Theme)
	if err != nil {
		app.Debug("\t\t\t\tloadThemeConfig() Unable to marshal YAML from %v", filename)
		// TODO: Handle error properly & and document error code
		return err
	}

	// TODO: Start using this to prevent multiple copies of the theme
	// Save the current theme. Force to lowercase because
	// it's  filename
	theme := strings.ToLower(app.Page.FrontMatter.Theme)
	app.Page.FrontMatter.Theme = theme
	app.Page.Theme.Name = theme
	app.Page.Theme.publishPath = path
	// TODO: This doesn't seem to be used
	app.Site.publishedThemes[path] = true

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

// newTheme() copies an existing root theme in the
// site themes directory and creates a new theme
// based on it, placing the new theme in the
// site themes directory.
// from is the name of the theme, not its path,
// e.g. "debut" or "pillar".
// TODO: Validate new name so it works as a slug/directory name
func (app *App) newTheme(from, to string, factory bool) error {
	if err := app.readSiteConfig(); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	app.Debug("newTheme(%v, %v, %v)", from, to, factory)
	app.Print("newTheme(from %v, to %v, factor? %v)", from, to, factory)
	// Get directory from which themes will be copied
	//source := filepath.Join(app.Site.siteThemesPath, defaults.SiteThemesDir, from)
	if from == "" {
		return ErrCode("1035", "")
	}
	source := filepath.Join(app.Site.siteThemesPath, from)

	// Get directory to which the theme will be copied for this site
	dest := filepath.Join(app.Site.siteThemesPath, to)
	app.Note("About to copy %v to %v", source, dest)
	app.Note("siteThemesPath: %v", filepath.Join(app.Site.siteThemesPath))

	//app.ShowInfo(".")
	return nil
}

// updateThemes() replaces all factory themes for the
// project with the latest ones
func (app *App) updateThemes() error {
	app.Note("updateThemes() NOT IMPLEMENTED")
	return nil
}
