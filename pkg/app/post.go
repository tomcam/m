package app

import (
	"fmt"
	"github.com/tomcam/m/pkg/default"
	"github.com/gosimple/slug"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// newPost() generates a post in a collection.
// collection  specifies the directory and any permalink variables.
// postname is a human-readable string that will get slugified into
// a filename (lowercase, spaces replaced
// with hyphens, all that jazz).
// That could easily require creating
// a directory. for example, the values  "/site/news/:year/:month:", "Site update"
// in November 2022 would generate a filename something like
// "/site/news/2022/11/site-update.md" and it would create the directory
// /site/news/2022/11/ if it's not already present.

func (app *App) newPost(collection, postname string) error {
	app.Debug("newPost(%v/%v)", collection, postname)
	// TODO: probaly need to normalize collection name with leading and trailing  directory separators
	// Ensure site is initialized properly
	if !app.Site.configLoaded {
		dir := currDir()
		if err := app.changeWorkingDir(dir); err != nil {
			return ErrCode("1114", dir)
		}
		if err := app.readSiteConfig(); err != nil {
			return ErrCode("PREVIOUS", err.Error())
		}
	}

	// Make sure the collection, e.g. "blog", has leading
	// and trailing slashes.
	pathSep := string(os.PathSeparator)
	if !strings.HasPrefix(collection, pathSep) {
		collection = filepath.Join(pathSep, collection)
	}
	if !strings.HasSuffix(collection, pathSep) {
		// TODO: This fails; seems to be by design.
		// filepath.Join() just doesn't like the
		// trailing directory separator.
		// See: https://go.dev/play/p/42x3HiccT6_S
		// collection = filepath.Join(collection, pathSep)
		collection = collection + pathSep
	}
	filename := string(app.Site.Collections[collection].Permalink)
	if filename == "" {
		msg := fmt.Sprintf( /* Can't find the permalink for */ "%v in starter file %v",
			collection, app.Site.starterFile)
		return ErrCode("1038", msg)
	}

	postname = slug.Make(postname)
	//	":year", ":month", ":monthnum", ":daynum", ":day", ":hour", ":minute", ":second", ":postname", ":author":

	//now := time.Now().Date()
	now := time.Now()
	day := now.Weekday().String()
	//year, monthnum, daynum := time.Now().Date()
	year, monthnum, daynum := now.Date()
	month := time.Month(int(monthnum))
	filename = strings.ReplaceAll(filename, ":year", strconv.Itoa(year))
	filename = strings.ReplaceAll(filename, ":month", month.String())
	filename = strings.ReplaceAll(filename, ":monthnum", strconv.Itoa(int(monthnum)))
	filename = strings.ReplaceAll(filename, ":daynum", strconv.Itoa(daynum))
	filename = strings.ReplaceAll(filename, ":day", day)
	filename = strings.ReplaceAll(filename, ":hour", strconv.Itoa(now.Hour()))
	filename = strings.ReplaceAll(filename, ":minute", strconv.Itoa(now.Minute()))
	filename = strings.ReplaceAll(filename, ":second", strconv.Itoa(now.Second()))
	if strings.Contains(collection, ":author") && app.Site.Author.FullName == "" {
		return ErrCode("1039", collection)
	}
	app.Note("Author: %v", app.Site.Author.FullName)
	filename = strings.ReplaceAll(filename, ":author", app.Site.Author.FullName)
	// TODO create directory
	dir := filename[1:strings.IndexRune(filename, ':')]
	dir = filepath.Join(app.Site.path, dir)
	if !dirExists(dir) {
    app.Note("Creating directory %v", dir)
		err := os.MkdirAll(dir, defaults.ProjectFilePermissions)
		if err != nil {
			return ErrCode("0412", dir)
		}
	}
	/// xxxjlk/m
	filename = strings.ReplaceAll(filename, ":postname", postname) + defaults.DefaultMarkdownExtension
	app.Note("\tFinished path %v", filename)
	return nil
}
