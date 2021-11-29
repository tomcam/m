package app

import (
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

	// List of all stylesheets needed to publish
	// this page after being massaged,
	// for example, to ensure responive.css comes
	// last among other things, and that only one of
	// theme-light.css or theme-dark.css are chosen.
	// This is different from Stylesheets, which is
	// list of all stylesheets listed in the theme config file.
	stylesheetList []string

	// Themes can be nested, e.g. debut/gallery/item.
	// Each level get its own entry here.
	// So in the case of debut/gallery/item,
	// levels[0] is 'debut', levels[1] is 'galery'
	// and levels[2] is 'item'
	levels []string

	// List of all levels of nested stylesheets as read in from
	// config file. So if the stylesheet is 'debut/gallery', then
	// you might have something like this:
	//   stylesheetsAllLevels['debut'] might be ['reset.css', 'fonts.css', 'bind.css', 'sizes.css', 'theme-light.css', 'theme-dark.css', 'layout.css', 'debut.css', 'responsive.css'
	//   stylesheetsAllLevels['gallery'] could be only  ['gallery.css'] if that's all you want changed
	stylesheetsAllLevels map[string][]string

	// List of all stylesheets mentioned in the theme config
	// file, whether they are needed to publish this page or
	// not. So it would contain things like theme-light.css
	// and theme-dark.css, even though only one is needed.
	// It doesn't assume stylesheets are in the corrected order
	// to be published. For example, in the published page
	// responsive.css needs to come list. It may not
	// appear list in this list. stylesheetList contains
	// only the subset of stylsheets needed to publish the page.
	Stylesheets []string `yaml:"Stylesheets"`

	// Name is the name of the theme for this page,
	// e.g. "wide"
	Name        string `yaml:"Name"`
	Branding    string `yaml:"Branding"`
	Description string `yaml:"Description"`

	Nav     layoutElement `yaml:"Nav"`
	Header  layoutElement `yaml:"Header"`
	Article layoutElement `yaml:"Article"`
	Footer  layoutElement `yaml:"Footer"`
	Sidebar layoutElement `yaml:"Sidebar"`
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
// gets distributed with Metabuz.
// So have the subdirectory available from within this
// subdirectory at compile time. Then you can run the finished
// executable anywhere and it will display the
// list of files even though the themes directory
// doesn't exist at runtime.

//go:embed themes/*
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
				app.Debug("\tos.MkdirAll() error: %v", err.Error())
				return ErrCode("0409", target)
			}
			app.Debug("\t\tcreated directory %v", target)
			return nil
		}
		// It's a file, not a directory
		app.Debug("\t\tCreated theme directory %v", target)
		// Handle individual file
		target = filepath.Join(app.Site.factoryThemesPath, path)
		f, err := factoryThemeFiles.Open(path)
		if err != nil {
			// TODO: Handle error properly & and document error code
			app.Debug("\tFS.Open(%v) error: %v", path, err.Error())
			return err
		}
		// Read the file into a byte array.
		b, err := io.ReadAll(f)
		if err != nil {
			// TODO: Handle error properly & and document error code
			app.Debug("\tio.ReadAll(%v) error: %v", f, err.Error())
			return err
		}
		// Copy the recently read file to its destination
		err = ioutil.WriteFile(target, b, defaults.ProjectFilePermissions)
		if err != nil {
			app.Debug("\t\tcopyFactoryThemes(): err after WriteFile:  %#v", err)
			// TODO: Handle error properly & and document error code
			return ErrCode("0216", err.Error(), target)
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

// loadNestedTheme() finds the theme specified for this page.
// It then copies the required files to the theme publihs
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
func (app *App) loadNestedTheme(source string, dest string, name string) error {
	app.Debug("\t\t\tloadNestedThemeLevel(%v, %v, %v)", source, dest, name)

	// See if this theme has already been published.
	_, ok := app.Site.publishedThemes[dest]
	if !ok {
		err := app.copyTheme(source, dest)
		// TODO: May want to improve error handling
		if err != nil {
			return ErrCode("PREVIOUS", err.Error())
		}
	} else {
		// Theme already loaded
		app.Print("Good news! %v theme already loaded", dest)
		return nil
	}

	app.Page.Theme.publishPath = dest

	// Theme directory is known. Load its config
	// (e.g. .yaml) file
	if err := app.loadThemeConfig(source); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	//app.Print("\t\t\t\t\t%v", app.Page.Theme.Stylesheets)
	app.Page.Theme.stylesheetsAllLevels[name] = app.Page.Theme.Stylesheets

	return nil
} // loadNestedTheme()

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
	app.Page.Theme.Name = fullTheme

	// Start building up the nested theme name, if any. So if it's
	// debut/gallery/item, it starts as debut, then
	// it's debut/gallery, then it's debut/gallery/item
	name := ""
	for level := 0; level < len(app.Page.Theme.levels); level++ {
		name = filepath.Join(name, app.Page.Theme.levels[level])
		// Get fully qualified directory from which themes will be copied
		// TODO: Isn't there an established way to do this?
		source := filepath.Join(app.Site.factoryThemesPath, defaults.SiteThemesDir, name)
		dest := filepath.Join(app.themePublishDir(name))
		// xxx
		// Get directory to which the theme will be copied for this site
		app.Page.Theme.sourcePath = source

		// Finds the theme specified for this page.
		// Copy the required files to the theme publish directory.
		if err := app.loadNestedTheme(source, dest, name); err != nil {
			return ErrCode("PREVIOUS", err.Error())
		}
		app.Page.Themes[name] = app.Page.Theme
	}

	return nil

} //loadTheme()

// themePublishDir() returns the name of the directory
// needed to publish style sheets for the theme
func (app *App) themePublishDir(theme string) string {
	//return filepath.Join(app.Site.cssPublishPath, defaults.ThemesDir, app.Page.FrontMatter.Theme)
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

	app.Page.Theme.publishPath = path
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
	app.Note("newTheme(%v, %v, %v)", from, to, factory)
	// Get directory from which themes will be copied
	//source := filepath.Join(app.Site.siteThemesPath, defaults.SiteThemesDir, from)
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
	app.Note("updateThemes()")
	return nil
}
