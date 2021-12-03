package main

import (
	"github.com/tomcam/m/pkg/app"
)

func main() {
	app := app.NewApp()
	app.Execute()
	app.Debug("\tfinished calling app.Execute()")
	if app.Flags.Info || app.Flags.InfoVerbose {
		app.ShowInfo("")
	}
	if app.Flags.InfoFrontMatter {
		app.ShowFrontMatter()
	}
}
