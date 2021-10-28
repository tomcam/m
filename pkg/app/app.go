package app

import (
	//"flag"
	//"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tomcam/m/pkg/default"
	"github.com/yuin/goldmark"
	"io"
	"os"
	"path/filepath"
	"strings"
	//"sync"
	//"context"
	"github.com/yuin/goldmark/parser"
)

// App contains all runtime options required to convert a markdown
// file or project to an HTML file or site.
// Compound data structure for config example at
// https://gist.github.com/alexedwards/5cd712192b4831058b21
type App struct {
	// Location (on startup) of user application data directory
	applicationDataPath string
	Site                Site

	// Cobra Command Processes command lin options
	//Cmd *cobra.Command
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

	// This is the global themes path. It's the source directory
	// containing the themes as they came from the factory.
	// It's where the site's themes come from when a new site
	// is created.
	// TODO: renamed from themesPath
	factoryThemesPath string
} // type Application

type Flags struct {
	// Global Debug mode. When true, App.Debug()
	// output gets displayed
	Debug bool

	// DontCopy means don't copy theme directory to the site directory.
	// Use the global theme set (which means if you change it, it
	// will affect all new sites created using that theme)
	DontCopy bool

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
		//deleteme: make([]byte)
		Page: Page{},
		Site: Site{},
		// Missing here: initializing the parser.
		// Can't set parser options until command
		// line has been processed.
		// So that happens at App.initConfig()
		RootCmd:             cobra.Command{},
		applicationDataPath: userConfigPath(),
		factoryThemesPath:   filepath.Join(userConfigPath(), defaults.ThemesDir),
	}
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
	app.Verbose("app.Execute()")
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

// setSiteDefaults() intializes the Site object
// It's on app instead of app.site so I can use
// global flags and debugging features
// like App.Note().
// Must be in the working directory at app.Site.path.
//func (app *App) setSiteDefaults(home string) {
func (app *App) setSiteDefaults() {
	app.Verbose("\tsetSiteDefaults()")
	app.setPaths()
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
	app.Site.name = filepath.Base(app.Site.path)
	// Compute location of base directory used for all
	// config info, which includes directories for
	// CSS files, graphic assets, HTML partials, etc.
	app.cfgPath = filepath.Join(app.Site.path, defaults.CfgDir)

	// Compute full pathname of the site file.
	app.Site.siteFilePath = filepath.Join(app.cfgPath,
		defaults.SiteConfigFilename)

	// Compute the publish directory  (aka WWW directory)
	app.Site.publishPath = filepath.Join(app.cfgPath,
		defaults.DefaultPublishPath)

	// Compute the directory location for asset files
	app.Site.assetPath = filepath.Join(app.Site.publishPath,
		defaults.DefaultAssetPath)

	// Compute the directory location for CSS files
	// to be published for this theme.
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
	// Don't want that name configurable.
	// Therefore cfgPath is enough.
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

	// Compute the directory location for theme files
	app.factoryThemesPath = filepath.Join(app.cfgPath,
		defaults.ThemesDir)

	// Create a new, empty map to hold the
	// source directory tree.
	app.Site.dirs = make(map[string]dirInfo)

} // setPaths()

// setWorkingDir() changes to the specified
// directory and sets app.site.path accordingly.
func (app *App) setWorkingDir(dir string) error {
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

// CopyMust() copies a single file named source to
// the file named in dest but doesn't return an
// error if something goes wrong.
// If a file got copied, returns its
// destination name
func (app *App) copyMust(src, dest string) string {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return ""
	}

	if !sourceFileStat.Mode().IsRegular() {
		return ""
	}
	source, err := os.Open(src)
	if err != nil {
		return ""
	}
	destination, err := os.Create(dest)
	if err != nil {
		return ""
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		return ""
	}
	// Success
	return dest
}

// cfgLower() reads a string value specified in
// the key parameter from configuration.
// A cfg value is one that can come from several
// places. For example, the theme name might normally
// come from the individual Page.FrontMatter setting.
// Or you might prefer to use Site.FrontMatter to
// set a default theme for the entire site, then change it
// only for specific pages in Page.FrontMatter.
//
// NOTE: The return value is forced to lowercase
func (app *App) cfgLower(key string) string {
	value := ""
	return strings.ToLower(value)
}
