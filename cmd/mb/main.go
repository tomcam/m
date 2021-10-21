package main

import (
	"github.com/tomcam/m/pkg/app"
	//"os"
)

func main() {
	//app := app.NewApp(pathname)
	app := app.NewApp()
	app.Execute()
	if app.Flags.Info || app.Flags.InfoVerbose {
		app.ShowInfo()
	}
	if app.Flags.InfoFrontMatter {
		app.ShowFrontMatter()
	}

}
