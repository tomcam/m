package app

import (
	"fmt"
	"github.com/rodaine/table"
	"github.com/spf13/viper"
	//"github.com/tomcam/m/pkg/default"
	//"os"
	"strings"
)

// frontMatter() displays the raw contents of the front matter
func (app *App) ShowFrontMatter() {
	fmt.Println("FRONT MATTER")
	fmt.Printf("%#v\n", app.Page.frontMatterRaw)
}

// ShowInfo() displays debug information about the app and site.
func (app *App) ShowInfo(pathname string) error {
	app.Print("Default theme: %v", viper.GetString("Theme"))
	// Change to specified directory.
	// Update app.Site.path and build all related directories
	if pathname == "" || pathname == "." {
		pathname = currDir()
	}
	if !isProject(pathname) {
		app.QuitError(ErrCode("0922", pathname))

	}

	if err := app.changeWorkingDir(pathname); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	if err := app.readSiteConfig(); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	}

	tbl := table.New("Site Directories", "")
	tbl.AddRow("Project name", app.Site.name)
	tbl.AddRow("Project directory", exists("", app.Site.path))
	tbl.AddRow("Collections", app.Site.Collections)
	tbl.AddRow("Config file directory", exists("", app.cfgPath))
	tbl.AddRow("Site file", exists("", app.Site.Filename))
	//tbl.AddRow("Asset path", app.Site.assetPath)
	//tbl.AddRow("Image path", app.Site.imagePath)
	tbl.AddRow("Common path", exists("", app.Site.commonPath))
	tbl.AddRow("Head tags path", exists("", app.Site.headTagsPath))
	tbl.AddRow("Publish path", exists("", app.Site.publishPath))
	//tbl.AddRow("CSS publish path", app.Site.cssPublishPath)
	tbl.AddRow("Factory themes path", exists("", app.Site.factoryThemesPath))
	tbl.AddRow("Site themes path", exists("", app.Site.siteThemesPath))
	tbl.AddRow("", "")
	tbl.AddRow("APPLICATION DATA", "")
	tbl.AddRow("User application data", exists("", app.applicationDataPath))
	tbl.Print()
	//tbl = table.New("Application Directories", "")
	return nil
}

// App.Verbose() displays a message followed
// by a newline to stdout
// if the verbose flag was used. Formats it like Fprintf.
func (a *App) Verbose(format string, ss ...interface{}) {
	if a.Flags.Verbose {
		fmt.Println(fmtMsg(format, ss...))
	}
}

// App.Debug() displays a message followed by a newline
// to stdout.
// Formats it like Fprintf.
func (a *App) Debug(format string, ss ...interface{}) {
	if a.Flags.Debug {
		fmt.Println(fmtMsg(format, ss...))
	}
}

// App.Note() displays a message followed by a newline
// to stdout, preceded by the text "NOTE: "
// For temporary use
// Overrides the verbose flag. Formats it like Fprintf.
func (a *App) Note(format string, ss ...interface{}) {
	fmt.Println("NOTE: " + fmtMsg(format, ss...))
}

// App.Print() displays a message followed by a newline
// to stdout.
// Formats it like Fprintf.
func (a *App) Print(format string, ss ...interface{}) {
	fmt.Println(fmtMsg(format, ss...))
}

// App.Warning() displays a message followed by a newline
// to stdout, preceded by the text "Warning: "
// Overrides the verbose flag. Formats it like Fprintf.
func (a *App) Warning(format string, ss ...interface{}) {
	fmt.Println("Warning: " + fmtMsg(format, ss...))
}

// fmtMsg() formats string like Fprintf and writes to a string
func fmtMsg(format string, ss ...interface{}) string {
	return fmt.Sprintf(format, ss...)
}
