package app

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tomcam/m/pkg/default"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// App contains all runtime options required to convert a markdown
// file or project to an HTML file or site.
// Compound data structure for config example at
// https://gist.github.com/alexedwards/5cd712192b4831058b21
type App struct {
	// Location (on startup) of user application data directory
	applicationDataPath string
	Site                Site

	// Cobra Command Processes command line options
	RootCmd cobra.Command
	// For viper
	cfgFile string

	// Location of directory containing themes, publish
	// directory, etc.
	cfgPath string

	// Global options such as Verbose
	Flags Flags

	// All built-in functions must appear here to be publicly available
	funcs map[string]interface{}
	// Copy of funcs but without "scode"
	fewerFuncs map[string]interface{}

	Page Page

	parser    goldmark.Markdown
	parserCtx parser.Context
} // type Application

type Flags struct {
	// Global Debug mode. When true, App.Debug()
	// output gets displayed
	Debug bool

	// DontCopy means don't copy theme directory to the site directory.
	// Use the global theme set (which means if you change it, it
	// will affect all new sites created using that theme)
	DontCopy bool

	// Factory means when you do a new theme command instead
	// of copying from  site theme, copy from a factory theme
	// instead.
	Factory bool

	// Name of starters file to generate pages when
	// site is created (or later)
	Starters string

	// Name of site config file to use when
	// site is created
	Site string

	// Global verbose mode
	Verbose bool

	// Display debug info with short pathnames
	Info bool
	// Display debug info
	InfoVerbose bool

	// Display front matter
	InfoFrontMatter bool
}

// NewApp allocates, and initializes to default
// values, an App struct (which contains the runtime
// environment for a Metabuzz process). Everything
// necessary to create a new project must be set
// by the time App.updateConfig() is called.
//
func NewApp() *App {
	app := App{
		Page: Page{},
		// Missing here: initializing the parser.
		// Can't set parser options until command
		// line has been processed.
		// So that happens at App.initConfig()
		RootCmd:             cobra.Command{},
		applicationDataPath: userConfigPath(),
	}
	// TODO: This didn't help newTheme()
	//app.Site.Starters = make(map[string]Starter)
	app.Site.publishedThemes = make(map[string]bool)
	return &app
}

// loadConfigs() looks for the many possible sources of
// configuration info (environment, local files, user
// document directory files, etc.)
// Call it after command line has been processed because
// the command line is our final, highest priority place
// to look for config info.
// Based on old initConfigs()
// https://github.com/tomcam/mb/blob/master/pkg/app/application.go#L57
func (app *App) loadConfigs() {
	app.Note("loadConfigs()")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func (app *App) Execute() {
	app = NewApp()
	app.Debug("app.Execute()")
	app.initCobra()
	cobra.CheckErr(app.RootCmd.Execute())
}

func (app *App) initCobra() {
	app.Verbose("initCobra()")
	app.addCommands()
	// RootCmd represents the base command when called without any subcommands
	//var RootCmd = &cobra.Command{

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	app.addFlags()

	cobra.OnInitialize(app.initConfig)

}

// initConfig reads in config file and ENV variables if set.
func (app *App) initConfig() {
	if app.cfgPath != "" {
		// Use config file from the flag.
		// XXX
		viper.SetConfigFile(app.cfgPath)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory
		viper.AddConfigPath(home)

		// Default config file extension, probably "yaml"
		viper.SetConfigType(defaults.ConfigFileDefaultExt)

		// Make the config file name obvious but "." to
		// hide it.
		viper.SetConfigName("." + defaults.ProductName)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		app.Note("Using config file:", viper.ConfigFileUsed())
	}

	// Parser couldn't be initialized until command line and
	// other options were processed
	app.parser = app.parserWithOptions()
	app.parserCtx = parser.NewContext()

	// Add snazzy Go template functions like ftime() etc.
	app.addTemplateFunctions()

}

// setPaths computes values for location of site
// theme files, publish directory, etc.
// Most of them are relative to the site directory.
// It must be called after command line flags, env
// variables, and other application configuration has been done.
// home is the fully qualified directory name
// the project lives in.
// Must be in the working directory at app.Site.path.
// This is based on App.SiteDefaults() in the previous
// version of Metabuzz.
func (app *App) setPaths() {
	app.Debug("setPaths")
	app.Site.name = filepath.Base(app.Site.path)
	// Compute location of base directory used for all
	// config info, which includes directories for
	// CSS files, graphic assets, HTML partials, etc.
	app.cfgPath = filepath.Join(app.Site.path, defaults.CfgDir)

	// Compute full pathname of the site file.
	app.Site.siteFilePath = filepath.Join(app.cfgPath,
		defaults.SiteConfigFilename)
	// Make a read-only public version for use in
	// templates: {{ Site.Filename }}
	app.Site.Filename = app.Site.siteFilePath

	// Compute the publish directory  (aka WWW directory)
	app.Site.publishPath = filepath.Join(app.cfgPath,
		defaults.DefaultPublishPath)

	// Compute the directory location for asset files
	app.Site.assetPath = filepath.Join(app.Site.publishPath,
		defaults.DefaultAssetPath)

	// Compute the directory location for CSS files
	// to be published for this theme.
	// TODO: I actually create the directory for this
	// in buildPublishDirs(). May need to revisit this...
	// maybe generate more dirs, maybe refactor it.
	app.Site.cssPublishPath = filepath.Join(app.Site.assetPath,
		defaults.DefaultPublishCssPath)

	// Compute the directory location for image files
	app.Site.imagePath = filepath.Join(app.Site.assetPath,
		defaults.DefaultPublishImgPath)

	// Compute the directory location for common files
	app.Site.commonPath = filepath.Join(app.cfgPath,
		defaults.CommonPath)

	// Compute the directory location for the
	// complete set of factory themes.
	// The entire /themes directory gets copied here.
	app.Site.factoryThemesPath = app.cfgPath

	// Compute the directory location of themes
	// that get copied over selectively for a
	// particular site.
	app.Site.siteThemesPath = filepath.Join(app.Site.publishPath,
		defaults.SiteThemesDir)

	// Compute the directory location for tags
	// that live in the HTML <head>
	app.Site.headTagsPath = filepath.Join(app.cfgPath,
		defaults.HeadTagsPath)

	// Create a new, empty map to hold the
	// source directory tree.
	app.Site.dirs = make(map[string]dirInfo)

} // setPaths()

// setWorkingDir() changes to the specified
// directory and sets app.Site.path accordingly.
func (app *App) setWorkingDir(dir string) error {
	app.Debug("setWorkingDir()")
	if dir == "." || dir == "" {
	} else {
		if err := os.Chdir(dir); err != nil {
			return ErrCode("PREVIOUS", dir)
		}
	}
	app.Site.path = currDir()
	app.setSiteDefaults()
	return nil
}

// cfgLower() reads a string value specified in
// the key parameter from configuration. Returns
// the value of that key forced to lowercase.
// A cfg value is one that can come from several
// places. For example, the theme name might normally
// come from the individual Page.FrontMatter setting.
// Or you might prefer to use Site.FrontMatter to
// set a default theme for the entire site, then change it
// only for specific pages in Page.FrontMatter.
func (app *App) cfgLower(key string) string {
	value := ""
	switch key {
	// Handle values we know could be in either or
	// both the front matter or site config file.
	case "Sidebar", "Theme":
		// Check front matter first.
		value = app.frontMatterMustLower(key)
		if value == "" {
			//App.Note("\tcfgLower(): sound find %v.[%v]",
			// Next fall back to site config
		}

		return strings.ToLower(value)
	}
	return value
}

// SiteMustLower() obtains the value of a
// requested key from the site configuration
// and returns it as a string forced to lowercase.
// The Must means it doesn't return an error
// if the key doesn't exist. It simply returns
// an empty string.
func (app *App) siteMustLower(key string) string {
	return structFieldByNameStrMust(app.Site, key)
}

// frontMatterMust() obtains the value of a
// requested key from the front matter.
// It's  called frontMatterMust() because it doesn't
// return an error if, for example, the requested
// doesn't exist, or doesn't have a definition.
// TODO: Perf? Get as []byte?
func (app *App) frontMatterMust(key string) string {
	// If the key exists, return its value.
	// This also works
	return structFieldByNameStrMust(app.Page.FrontMatter, key)
}

// frontMatterMustLower() obtains the value of a
// requested key from the front matter, then
// forces the return value to lowercase.
// The Must means it doesn't return an error
// if the key doesn't exist. It simply returns
// an empty string.
func (app *App) frontMatterMustLower(key string) string {
	// If the key exists, return its value.
	return strings.ToLower(app.frontMatterMust(key))
}

// siteThemesPath() determines the directory the
// theme file is found in.
func (app *App) siteThemesPath() string {
	return filepath.Join(app.Site.siteThemesPath, app.Page.Theme.Name)
}

// frontMatterRawToStruct() takes the generic map of front
// matter produced by Goldmark's YAML parser
// and copies it to the FrontMatter struct.
// The *Must functions are used because its structure
// is known so why check for errors.
// I'll probably regret this but what the hay.
func (app *App) frontMatterRawToStruct() {
	for k, v := range app.Page.frontMatterRaw {
		setFieldMust(&app.Page.FrontMatter, k, v)
	}
}

// Based on https://stackoverflow.com/a/26746461
func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

// setFieldMust() is identical to setField but strips the
// error checking.
// Use in cases like frontMatterRawtoStruct()
// where the structure type is known in advance.
func setFieldMust(obj interface{}, name string, value interface{}) {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return
	}

	if !structFieldValue.CanSet() {
		return
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return
	}

	structFieldValue.Set(val)
	return
}
