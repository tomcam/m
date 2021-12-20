package app

import (
	"github.com/tomcam/m/pkg/default"
	"github.com/tomcam/m/pkg/util"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Theme struct {
	// Determines what Metabuzz Theme Framework features
	// are supported by this theme.
	Supports supports `yaml:"Supports"`

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
	publishStylesheets []string

	// Name is the slug name of the theme for this page,
	// e.g. "wide"
	Name string `yaml:"Name"`

	// Branding is the pretty name of the theme,
	// say, "Wide" or "Metabuzz Wide"
	Branding string `yaml:"Branding"`

	// A sentence or two describing why to use this theme
	Description string `yaml:"Description"`

	// Raw list of all possible stylesheets required to
	// publish, but some actually may be omitted
	// on publish. For example, only sidebar-right.css
	// or sidebar-left.css, but not both.
	Stylesheets []string `yaml:"Stylesheets"`

	// Copyright/licensing terms, e.g. GPL 3.0, Apache,
	// whatever
	License string `yaml:"License"`

	// Designer of the theme
	Author string `yaml:"Author"`

	// Page layout elements. Sidebar becomes the
	// "aside" tag
	Nav     layoutElement `yaml:"Nav"`
	Header  layoutElement `yaml:"Header"`
	Article layoutElement `yaml:"Article"`
	Footer  layoutElement `yaml:"Footer"`
	Sidebar layoutElement `yaml:"Sidebar"`
} // type Theme

// Describes what elements of the
// Metabuzz Theme Framework (MTF) this theme supports.
type supports struct {
	// Supports the theme framework itself
	// If this is false then all the rest are false
	MTF bool `yaml:"MTF"`

	// Theme supports dark mode/light mode
	Mode bool `yaml:"Mode"`

	// Theme supports header
	Header bool `yaml:"Header"`

	// Theme suports navbar
	Nav bool `yaml:"Nav"`

	// Theme supports sidebar
	Sidebar bool `yaml:"Sidebar"`

	// Theme supports footer.
	Footer bool `yaml:"Footer"`
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
// gets distributed with Metabuzz.
// So have the subdirectory available from within this
// subdirectory at compile time. Then you can run the finished
// executable anywhere and it will display the
// list of files even though the themes directory
// doesn't exist at runtime.

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
	app.Note("\t\t\t\tcopyTheme(%v, %v)", source, dest)
	err := os.MkdirAll(dest, defaults.PublicFilePermissions)
	if err != nil {
		// TODO: Handle error properly & and document error code
		app.Debug("\t\t\t\tos.MkdirAll(%v) error: %v", dest, err.Error())
		return ErrCode("PREVIOUS", err.Error())
	}
	// TODO: See if there's a faster Go lib for this
	if err := CopyDirectory(source, dest, false); err != nil {
		app.Note("CopyDirectory(%s, %s, false) failed", source, dest)
		// TODO: Problbably want a an original error code
		return ErrCode("PREVIOUS", err.Error())
	}
	/*
		if err := copyDirAll(source, dest); err != nil {
			app.Debug("\t\t\t\tFAILED copyDirAll(%v, %v)", source, dest)
			return ErrCode("PREVIOUS", err.Error())
		}
	*/
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

// xxx

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

	// Add YAML or whatever to the base directory because that's the
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
//func (app *App) newTheme(from, to string, factory bool) error {
func (app *App) newTheme(from, to string) error {
	dir := currDir()
	if err := os.Chdir(dir); err != nil { // TODO: Handle error properly & and document error code
		// TODO: Change this if it works
		return ErrCode("1105", dir)
	}
	app.Site.path = dir
	app.setSiteDefaults()
	if err := app.readSiteConfig(); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	app.Debug("newTheme(%v, %v)", from, to)
	// Get directory from which themes will be copied
	if from == "" {
		return ErrCode("1035", "")
	}

	// Get directory to which the theme will be copied for this site
	dest := filepath.Join(app.Site.siteThemesPath, to)

	if dirExists(dest) {
		// TODO: create test case
		// Target theme directory already there
		return ErrCode("0952", dest)
	}

	// Derive theme filename from theme path to theme directory
	source := filepath.Join(app.Site.siteThemesPath, from)
	sourceCfgFile := themeCfgName(source)
	var theme Theme
	// TODO: create test case
	if err := readYAMLFile(sourceCfgFile, &theme); err != nil {
		return ErrCode("0132", sourceCfgFile)
	}

	// Create destination directory
	if err := os.MkdirAll(dest, defaults.PublicFilePermissions); err != nil {
		return ErrCode("0413", dest)
	}
	app.Debug("About to copy %v to %v using %v", source, dest)
	if err := app.copyTheme(source, dest); err != nil {
		return ErrCode("0931", "from '"+source+"' to '"+dest+"'")
	}

	// The destination theme has been copied but it still has the old
	// names inside. Maket those corrections.
	if err := app.copyThemeUpdate(source, dest); err != nil {
		return ErrCode("0930", "from '"+source+"' to '"+dest+"'")
	}
	return nil
}

// copyThemeUpdate() takes place after a source theme has been copied
// to the destination directory.  Pass it the fully qualified source and
// destination directories for the themes.
// The destination theme config filename and other
// details now need to be updated, because they still have the source
// theme names. So if you were copying "debut" to "newdeb", you'd
// want to change debut.yaml to newdeb.yaml, debut.css to newdeb.css,
// etc.
// Currently works only with root themes, not children.
func (app *App) copyThemeUpdate(source, dest string) error {
	// if "themes/foo" had been copied to "themes/bar",
	// then "themes/bar/foo.yaml" has to be renamed
	// "themes/bar/bar.yaml"
	// (Simplifying pathnames for purposes of illustration)
	// So build up the bad filename "themes/bar/foo.yaml"
	destDir := filepath.Dir(dest)
	sourceThemeName := filepath.Base(source)
	destThemeName := filepath.Base(dest)
	badCfgFile := filepath.Join(destDir, destThemeName, sourceThemeName+"."+defaults.ConfigFileDefaultExt)
	// Get the filename for the corrected theme config file.
	// In this example, it's "bar"
	// Derive the correct filename for it, in this example
	// "themes/bar/bar.yaml"
	destCfgFile := filepath.Join(destDir, destThemeName, destThemeName+"."+defaults.ConfigFileDefaultExt)

	// And rename the bad file to the new one.
	if err := os.Rename(badCfgFile, destCfgFile); err != nil {
		return ErrCode("0223", "from '"+source+"' to '"+dest+"'")
	}
	// Now go through the stylesheets. See if there's one by old name,
	// in this example, "foo.css". If so, change it to "bar.css".
	var theme Theme
	// TODO: create test case
	b, err := ioutil.ReadFile(destCfgFile)
	if err != nil {
		app.Debug("\t\t\t\tfailed to read %v", destCfgFile)
		// TODO: Handle error properly & and document error code
		return ErrCode("0132", destCfgFile)
	}
	if err := yaml.Unmarshal(b, &theme); err != nil {
		// TODO: Handle error properly & and document error code
		return ErrCode("0134", destCfgFile)
	}

	newStylesheets := []string{}
	// See if there's a CSS file by the same name as
	// the theme. If so, rename it to the new theme name.
	// This means reading all the themes in the config
	// file, writing them to a new array (renaming
	// the theme-named stylesheet if found), creating
	// the new list of stylesheets, and writing out
	// a new config file with the updated list.
	for _, stylesheet := range theme.Stylesheets {
		if stylesheet == sourceThemeName+".css" {
			oldStyleFile := filepath.Join(destDir, destThemeName, stylesheet)
			stylesheet = destThemeName + ".css"
			newStyleFile := filepath.Join(destDir, destThemeName, stylesheet)
			// Rename the file
			if err := os.Rename(oldStyleFile, newStyleFile); err != nil {
				return ErrCode("0223", "from '"+oldStyleFile+"' to '"+newStyleFile+"'")
			}

		}
		newStylesheets = append(newStylesheets, stylesheet)
	}
	theme.Stylesheets = newStylesheets
	// Write out the new cfg file, with the search/replaced stylesheet name in the
	// Stylesheets list.
	if err := writeYamlFile(destCfgFile, &theme); err != nil {
		// TODO: This may not be the correct move. Add a code?
		return ErrCode("0225", destCfgFile)
	}

	return nil

}

// themeCfgName() generatestheme name from fully qualified filename.
// Does not check to see if such a file exists.
func themeCfgName(pathname string) string {
	//cfgName := filepath.Base(pathname, defaults.ConfigFileDefaultExt)
	theme := filepath.Base(pathname)
	return filepath.Join(pathname, theme+"."+defaults.ConfigFileDefaultExt)
}

// updateThemes() replaces all factory themes for the
// project with the latest ones
func (app *App) updateThemes() error {
	app.Note("updateThemes() NOT IMPLEMENTED")
	return nil
}
