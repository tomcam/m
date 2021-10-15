package app

import (
	"github.com/spf13/cobra"
	"github.com/tomcam/m/pkg/default"
	"github.com/yuin/goldmark"
)

// App contains all runtime options required to convert a markdown
// file or project to an HTML file or site.
// Compound data structure for config example at
// https://gist.github.com/alexedwards/5cd712192b4831058b21
type App struct {
	site *Site

	// Cobra Command Processes command lin options
	Cmd *cobra.Command

	Parser goldmark.Markdown

	// Contents of HTML file after being converted from Markdown
	HTML []byte
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
		HTML:   nil,
		Parser: goldmark.New(),
		Cmd: &cobra.Command{
			Use:   defaults.ProductShortName,
			Short: "Create static sites",
			Long:  `Headless CMS to create static sites`,
		},
	}

	// Process command line
	app.addCommands()

	// If there are any configuration files,
	// environment variables, etc. with info
	// that overrides what was just initialized,
	// read them in.
	app.updateConfig()

  var err error
  if app.site, err = app.site.New(); err != nil {
    app.QuitError(ErrCode("PREVIOUS", err.Error()))
  }
	return &app
}

// updateConfig() determines where configuration file (and other
// forms of configuration info, such as
// environment variables) can be found, then reads in
// all that info. It overrides defaults established
// in NewApp(). It isn't necessary. That us, NewApp()
// will have initialized the App data structure sufficiently
// to create a new project in the absence of any
// overriding config information.
func (app *App) updateConfig() {
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
}
