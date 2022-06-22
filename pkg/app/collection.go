package app

import (
	"fmt"
	"github.com/tomcam/m/pkg/default"
	"os"
	"path/filepath"
	"strings"
)

// Repreents
type Collection struct {
	// Format used to generate a path to the post.
	Permalink string `yaml:"Permalink"`

	// Sort directory by filename, date, whatever
	Sort string `yaml:"Sort"`
}

// newCollection() creates a collection from the path and/or
// permalink described in name, for example,
// "/blog/:year/:monthnum/:daynum"
// And adds it to Site.Collections
// filename is the name of the starter file. If interactive,
// just pass ""
func (app *App) newCollection(name string /* starter Starter, */, filename string) error {
	// Can be called from command line so make sure all is initialized
	if !app.Site.configLoaded {
		if err := app.changeWorkingDir(currDir()); err != nil {
			app.Debug("\tUnable to change to directory (%v)", currDir())
			return ErrCode("1109", currDir())
		}
		if err := app.readSiteConfig(); err != nil {
			return ErrCode("PREVIOUS", err.Error())
		}
	}
	app.Debug("newCollection(%v)", name)
	// The name is a path to the file or collection.
	// It may also be a permalink.
	path := name
	// If no folder is given, assume project root.
	// Remember Go uses Unix folder conventions even
	// under Windows
	if path == "" {
		path = "/"
	}

	// Ensure it has a leading "/" unless preceded with a dot
	if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, ".") {
		path = filepath.Join("/", path)
	}

	var c Collection
	var permalink string
	var err error
	// TODO: Need test case
	if permalink, err = fixPermalink(path, filename); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}

	// The key to the collection will be the base directory
	// (everything up to the first permalink variable)
	base := permalinkBase(permalink)

	// Quit if this is a duplicate
	if app.Site.Collections[base].Permalink == permalink {
		msg := ""
		if filename != "" {
			msg = fmt.Sprintf("Starter file %s already has a collection named %s", filename, base)
		} else {
			msg = fmt.Sprintf("Collection named %s already exists", base)
		}
		return ErrCode("0954", msg)
	}
	c.Permalink = permalink
	app.Debug("\tAbout to add %#v to %#v at %v", c, app.Site.Collections, base)
	app.Site.Collections[base] = c
	app.Debug("\tapp.Site.Collections[%v] is now %v", base, app.Site.Collections[base])
	//app.Print("Permalink: %v\n app.Site.Collections.[permalink]: %v\n. FirstDir: %v\n", permalink, app.Site.Collections[permalink], permalinkBase(permalink))

	// Create the specified folder as a subdirectory
	// of the current project.
	// The base directory is everything up to the first
	// colon. Permalink is guaranteed to start
	// with a directory separator.
	// TODO: refactor with permalinkBase()?
	dir := permalink[1:strings.IndexRune(permalink, ':')]
	err = os.MkdirAll(dir, defaults.ProjectFilePermissions)
	if err != nil {
		return ErrCode("0415", dir)
	}

	if err = app.writeSiteConfig(); err != nil {
		return ErrCode("1301", name)
	}
	return nil
}
