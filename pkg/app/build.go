package app

import (
	"github.com/tomcam/m/pkg/default"
	"os"
  //"errors"
)

func (app *App) build(pathname string) error {
	if pathname != "" {
		// Change to the specified directory.
		if err := os.Chdir(pathname); err != nil {
			return ErrCode("0901", err.Error())
		}
	}

	// Determine current fully qualified directory location.
	// Can't use relative paths internally.
	pathname = currPath()

	// Changed directory successfully so
	// pass it to initialize the site and update internally.
	app.site.defaults(pathname)

	app.Note("app.build(%s)", app.site.path)
	//app.Note("site.siteFilePath: %s", app.site.siteFilePath)

  // Create minimal directory structure: Publish directory,
	// site directory, .themes, etc.
  var err error
	if err = createDirStructure(&defaults.SitePaths); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}

	// Delete any existing publish dir
	if err := os.RemoveAll(app.site.publishPath); err != nil {
		return ErrCode("0302", app.site.publishPath)
	}

	// Create an empty publish dir
	if err := os.MkdirAll(app.site.publishPath, defaults.PublicFilePermissions); err != nil {
		return ErrCode("0403", app.site.publishPath)
	}

	// Return with success code.
	return nil
}
