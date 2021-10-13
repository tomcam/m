package main

import (
	"github.com/tomcam/m/pkg/app"
)

func main() {

	a := app.NewApp()
	if err := a.Cmd.Execute(); err != nil {
		a.Note("Execution error?")
		a.QuitError(err)
	} else {
		a.Note("Executed without error, I guess")
	}
}
