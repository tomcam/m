package app

import (
	"os"
)

func (app *App) build(pathname string) error {
	if pathname != "" {
		// Change to the specified directory.
		if err := os.Chdir(pathname); err != nil {
			//app.QuitError(ErrCode("1101", err.Error(), pathname))
			return ErrCode("1101", err.Error(), pathname)
		}
	}

	// Determine current fully qualified directory location.
	// Can't use relative paths internally.
	pathname = currPath()
	// Changed directory successfully so
	// pass it to initialize the site and update internally.
	app.site.defaults(pathname)
	app.Note("app.build(%s)", app.site.path)
  app.Note("site.siteFilePath: %s", app.site.siteFilePath)
	// Return with success code.
	return nil
}
