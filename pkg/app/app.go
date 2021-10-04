package app

import (
  //"github.com/tomcam/m/pkg/util"

  /*
	"github.com/spf13/viper"
	"github.com/spf13/cobra"
	"github.com/tomcam/mb/pkg/errs"
	"github.com/yuin/goldmark"
	"path/filepath"
  */
)

// App contains all runtime options required to convert a markdown
// file or project to an HTML file or site.
// Compound data structure for config example at
// https://gist.github.com/alexedwards/5cd712192b4831058b21
type App struct {

  Site Site
}

// NewApp allocates, and initializes to default
// values, an App struct (which contains the runtime
// environment for a Metabuzz process). Everything
// necessary to create a new project must be set
// by the time App.updateConfig() is called.
// 
// path is the location for the project.
// 
func NewApp(path string) *App {
  app := App{}

  // path is the location for the project.
  // If not specified, use the current directory.
  if path != "" {
    app.Site.Path = path
  } else {
    app.Site.Path = currPath()
  }

  // If there are any configuration files,
  // environment variables, etc. with info
  // that overrides what was just done,
  // readt them in.
  app.updateConfig()

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

