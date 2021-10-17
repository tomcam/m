package app

import (
	//"flag"
	//"fmt"
	//"github.com/tomcam/m/pkg/default"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yuin/goldmark"
	"os"
	//"sync"
	//"context"
	"github.com/yuin/goldmark/parser"
)

// App contains all runtime options required to convert a markdown
// file or project to an HTML file or site.
// Compound data structure for config example at
// https://gist.github.com/alexedwards/5cd712192b4831058b21
type App struct {
	site Site

	// Cobra Command Processes command lin options
	//Cmd *cobra.Command
	RootCmd cobra.Command
	// For viper
	cfgFile string

	// Global options such as Verbose
	flags Flags

	page Page

	parser    goldmark.Markdown
	parserCtx parser.Context

	// Contents of HTML file after being converted from Markdown
	HTML []byte
}
type Flags struct {
	// DontCopy means don't copy theme directory to the site directory.
	// Use the global theme set (which means if you change it, it
	// will affect all new sites created using that theme)
	DontCopy bool

	// Global verbose mode
	Verbose bool

	// Display debug info
	Info bool

	QTest bool
	RTest bool
}

// NewApp allocates, and initializes to default
// values, an App struct (which contains the runtime
// environment for a Metabuzz process). Everything
// necessary to create a new project must be set
// by the time App.updateConfig() is called.
//
// path is the location for the project.
//
func NewApp() *App {
	app := App{
		page:   Page{},
		site:   Site{},
		parser: goldmark.New(),
		//parser: parserWithOptions(),
		parserCtx: parser.NewContext(),
		RootCmd:   cobra.Command{},
	}

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
	app.Verbose("app.Execute()")
	app.initCobra()
	cobra.CheckErr(app.RootCmd.Execute())
}

func (app *App) initCobra() {
	app.addCommands()
	app.Note("initCobra()")
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
	if app.cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(app.cfgFile)
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

func (app *App) qTest() {
	app.Note("qTest()")
}
