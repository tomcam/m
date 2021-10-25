package app

import (
	//"flag"
	//"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tomcam/m/pkg/default"
	"github.com/yuin/goldmark"
	"os"
	"path/filepath"
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
	site                Site

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

	page Page

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
		page: Page{},
		site: Site{},
		//parser: goldmark.Markdown,
		//parser: parserWithOptions(),
		parserCtx:           parser.NewContext(),
		RootCmd:             cobra.Command{},
		applicationDataPath: userConfigPath(),
		factoryThemesPath:   filepath.Join(userConfigPath(), defaults.ThemesDir),
	}
	app.setSiteDefaults("")
	// TODO: Get values from viper here I think.
	// If there are any configuration files,
	// environment variables, etc. with info
	// that overrides what was just initialized,
	// read them in.
	//app.updateConfig()
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
	app.Note("initConfig()")
	if app.cfgPath != "" {
		// Use config file from the flag.
		// XXX
		viper.SetConfigFile(app.cfgPath)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".mb" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mb")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		app.Note("Using config file:", viper.ConfigFileUsed())
	}
}

// setSiteDefaults() intializes the Site object
// It's on app instead of app.site so I can use
// read global flags and use debugging features
// like App.Note(). home is the project directory
// or, if left as "", is set to the current directory
func (app *App) setSiteDefaults(home string) {
	// By now you should already be in the directory
	// specified by home
	if home == "" {
		home = currDir()
	}
	app.Verbose("\tsetSiteDefaults(%v)\n", home)
	app.setPaths(home)
}

// setPaths computes values for location of site
// theme files, publish directory, etc.
// Most of them are relative to the site directory.
// It must be called after command line flags, env
// variables, and other application configuration has been done.
// home is the fully qualified directory name
// the project lives in.
// By now it is the current directory.
// This is based on App.SiteDefaults() in the previous
// version of Metabuzz.
func (app *App) setPaths(home string) {
	// This is the fully qualified path of the current
	// directory, which is also guaranteed to be the
	// root directory of the project.
	app.Verbose("\tsetPaths(%v)\n", home)
	app.site.path = home
	app.site.name = filepath.Base(app.site.path)
	// Compute location of base directory used for all
	// config info, which includes directories for
	// CSS files, graphic assets, HTML partials, etc.
	app.cfgPath = filepath.Join(app.site.path, defaults.CfgDir)

	// Compute full pathname of the site file.
	app.site.siteFilePath = filepath.Join(app.cfgPath,
		defaults.SiteConfigFilename)

	// Compute the publish directory  (aka WWW directory)
	app.site.publishPath = filepath.Join(app.cfgPath,
		defaults.DefaultPublishPath)

	// Compute the directory location for asset files
	app.site.assetPath = filepath.Join(app.cfgPath,
		defaults.DefaultAssetPath)

	// Compute the directory location for CSS files
	app.site.cssPath = filepath.Join(app.site.assetPath,
		defaults.DefaultPublishCssPath)

	// Compute the directory location for image files
	app.site.imagePath = filepath.Join(app.site.assetPath,
		defaults.DefaultPublishImgPath)

	// Compute the directory location for common files
	app.site.commonPath = filepath.Join(app.cfgPath,
		defaults.CommonPath)

	// Compute the directory location for factory
	// themes for this project.
	// The entire /themes directory gets copied here.
	// Don't want that name configurable.
	// Therefore cfgPath is enough.
	app.site.factoryThemesPath = app.cfgPath

	// Compute the directory location for tags
	// that live in the HTML <head>
	app.site.headTagsPath = filepath.Join(app.cfgPath,
		defaults.HeadTagsPath)

	// Compute the directory location for theme files
	app.factoryThemesPath = filepath.Join(app.cfgPath,
		defaults.ThemesDir)

	// Create a new, empty map to hold the
	// source directory tree.
	app.site.dirs = make(map[string]dirInfo)

} // setPaths()
