package app

import (
	"fmt"
	"github.com/rodaine/table"
	//"os"
	"strings"
)

// frontMatter() displays the raw contents of the front matter
func (app *App) ShowFrontMatter() {
	fmt.Println("FRONT MATTER")
	fmt.Printf("%#v\n", app.page.frontMatter)
}

// ShowInfo() displays debug information about the app and site.
func (app *App) ShowInfo(pathname string) error {
  // Change to specified directory.
  // Update app.site.path and build all related directories
   if err := app.setWorkingDir(pathname); err != nil {
    return ErrCode("PREVIOUS", err.Error())
  }

	table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
		return strings.ToUpper(fmt.Sprintf(format, vals...))
	}

	tbl := table.New("Site Directories", "")
	tbl.AddRow("Project name", app.site.name)
	tbl.AddRow("Project directory", app.site.path)
	tbl.AddRow("Config file directory", app.cfgPath)
	tbl.AddRow("Site file", app.site.siteFilePath)
	tbl.AddRow("Asset path", app.site.assetPath)
	tbl.AddRow("CSS path", app.site.cssPath)
	tbl.AddRow("Image path", app.site.imagePath)
	tbl.AddRow("Common path", app.site.commonPath)
	tbl.AddRow("Head tags path", app.site.headTagsPath)
	tbl.AddRow("Publish path", app.site.publishPath)
	tbl.AddRow("Factory themes path", app.site.factoryThemesPath)
	tbl.AddRow("APPLICATION DATA", "")
	tbl.AddRow("User application data", app.applicationDataPath)
	tbl.Print()
	//tbl = table.New("Application Directories", "")
	return nil
}

// App.Verbose() displays a message followed
// by a newline to stdout
// if the verbose flag was used. Formats it like Fprintf.
func (a *App) Verbose(format string, ss ...interface{}) {
	if a.Flags.Verbose {
		fmt.Println(a.fmtMsg(format, ss...))
	}
}

// App.Note() displays a message followed by a newline
// to stdout, preceded by the text "NOTE: "
// For temporary use
// Overrides the verbose flag. Formats it like Fprintf.
func (a *App) Note(format string, ss ...interface{}) {
	fmt.Println("NOTE: " + a.fmtMsg(format, ss...))
}

// App.Warning() displays a message followed by a newline
// to stdout, preceded by the text "Warning: "
// Overrides the verbose flag. Formats it like Fprintf.
func (a *App) Warning(format string, ss ...interface{}) {
	fmt.Println("Warning: " + a.fmtMsg(format, ss...))
}

// fmtMsg() formats string like Fprintf and writes to a string
func (a *App) fmtMsg(format string, ss ...interface{}) string {
	return fmt.Sprintf(format, ss...)
}
