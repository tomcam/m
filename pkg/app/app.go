package app

import (
	// "fmt"
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
// environment for a Metabuzz process).
// path is the location for the project.
//func (app *App) NewApp(path string) *App {
func NewApp(path string) *App {
  a := App{
    //app.init(path)
    //init(path)
  }
  a.init(path)
  return &a
}


// init determines where configuration file (and other
// forms of configuration info) can be found, then reads in
// all that info. If there's no config file found, generates
// sensible defaults.
// path is the location for the project.
// If "", use the current directory
func (app *App) init(path string) {
  if path != "" {
    app.Site.Path = path
  }
}

