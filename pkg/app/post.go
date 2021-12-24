package app

import (
	"github.com/gosimple/slug"
  "strings"
	"path/filepath"
	"os"
)

func (app *App) newPost(collection, postname string) error {
	app.Print("newPost(%v/%v)", collection, postname)
	// TODO: probaly need to normalize colleciton name with leading and trailing  directory separators
	// Ensure site is initialized properly
	if !app.Site.configLoaded {
		dir := currDir()
		if err := app.changeWorkingDir(dir); err != nil {
      app.Note("cojuldn't changedir")
			return ErrCode("1113", dir)
		}
		if err := app.readSiteConfig(); err != nil {
      app.Note("cojuldn't read site config")
			return ErrCode("PREVIOUS", err.Error())
		}
	}

  app.Print("HI")
  pathSep := string(os.PathSeparator)
  pathSep = "/"
  app.Print("pathSep: %v", pathSep)
  if !strings.HasPrefix(collection, pathSep) {
    collection = filepath.Join(pathSep, collection)
    app.Note("No preix. created %v", collection)
  }
  if !strings.HasSuffix(collection, pathSep) {
    // TODO: This is all failing!
    //collection = filepath.Join(collection, pathSep)
    collection = collection + pathSep
    app.Note("No suffix. created %v", collection)
  }
  filename := string(app.Site.Collections[collection].Permalink)

  app.Print("\tPremalink is now: %v", filename)
	// TODO: Handle internal error where empty string is return. Create test case.

	postname = slug.Make(postname)

  filename = strings.ReplaceAll(filename, ":postname", postname)
	app.Print("\tAbout to create %v", filename)

	return nil
}
