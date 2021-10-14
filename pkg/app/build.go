package app

import (
	"os"
)

func (app *App) build(pathname string) error {
	app.Note("app.build()")
  if pathname == "" {
    app.Note("\tCurrent directory")
  } else {
    // No argument to build, so just do it
    // in the current directory.
    if err := os.Chdir(pathname); err != nil {
      app.QuitError(ErrCode("1101", err.Error(), pathname))
    }
    // Changed directory successfully so update it internally.
    app.site.path = pathname
    app.Note("\tLocation: %s", app.site.path)
  }
	return nil
}
