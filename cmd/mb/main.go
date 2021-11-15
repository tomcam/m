package main

import (
	"github.com/tomcam/m/pkg/app"
)

func main() {
	//fmt.Printf("main()\n")
	// I may need to move NewApp invidivually to createNewSite, etc.
	//app := app.NewApp(pathname)
	//fmt.Printf("main() about to call NewApp()\n")
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
