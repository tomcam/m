package app

import (
	"github.com/gosimple/slug"
/*
  "os"
  "fmt"
	"path/filepath"
	"github.com/tomcam/m/pkg/default"
*/
)

func (app *App) newPost(collection, postname string) error {
	app.Print("newPost(%v/%v)", collection, postname)
  // TODO: probaly need to normalize colleciton name with leadin gand trailing  directory separators
  permalink := app.Site.Collections[collection]
  // TODO: Handle internal error where empty string is return. Create test case.

	postname = slug.Make(postname)
  app.Print("\tAbout to create %v/%v", permalink,postname)

  app.Print("\tCollections: %#v", app.Site.Collections)
	return nil
}
