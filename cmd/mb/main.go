package main

import (
	"github.com/tomcam/m/pkg/app"
)

func main() {
	app := app.NewApp()
	app.Note("main: call info() after app.NewApp()")
	if app.Flags.Info || app.Flags.InfoVerbose {
		app.ShowInfo()
	}
	if app.Flags.InfoFrontMatter {
		app.ShowFrontMatter()
	}

	app.Execute()
}
