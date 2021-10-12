package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tomcam/m/pkg/default"
	"github.com/yuin/goldmark"
)

// App contains all runtime options required to convert a markdown
// file or project to an HTML file or site.
// Compound data structure for config example at
// https://gist.github.com/alexedwards/5cd712192b4831058b21
type App struct {
	Site Site
	Cmd  *cobra.Command

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
func NewApp(path string) *App {
	app := App{
		HTML: nil,
		Cmd: &cobra.Command{
			Use:   defaults.ProductShortName,
			Short: "Create static sites",
			Long:  `Headless CMS to create static sites`,
		},
	}
	app.Parser = goldmark.New()
	//app.Parser = mark.GetParser().Parser()
	//app.parser = mark.GetParser().Parser()
	//app.parser = markdown.GetParser().Parser()
	//app.parser = markdown.GetParser().Parser()

	// path is the location for the project.
	// If not specified, use the current directory.
	if path != "" {
		app.Site.Path = path
	} else {
		app.Site.Path = currPath()
	}

	// If there are any configuration files,
	// environment variables, etc. with info
	// that overrides what was just initialized,
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

func (app *App) build() error {
	fmt.Println("FAKE BUILD")
	return nil
}
func (app *App) NewSite() error {
	// Create minimal directory structure: Publish directory
	// .site directory, .themes, etc.
	if err := createDirStructure(&defaults.SitePaths); err != nil {
		//return errs.ErrCode("PREVIOUS", err.Error())
		return ErrCode("PREVIOUS", err.Error())
	}
	return nil
}
