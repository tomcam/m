package app

import (
	"fmt"
	"github.com/tomcam/m/pkg/default"
	"io/ioutil"
	"os"
	"path/filepath"
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
	app.Note("app.site.defaults(%v)", pathname)
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
	// TODO: I think buildPublishPath() obsoletes this
	if err := os.MkdirAll(app.site.publishPath, defaults.PublicFilePermissions); err != nil {
		app.Note("Unable to create path %v", app.site.publishPath)
		//return ErrCode("0403", app.site.publishPath,"" )
		return ErrCode("PREVIOUS", err.Error())
	}

	// Get a list of all files & directories in the site.
	if _, err = app.getProjectTree(app.site.path); err != nil {
		return ErrCode("0913", app.site.path)
	}

 // Build the target directory true
  app.buildPublishDirs()

	// Loop through the list of permitted directories for this site.
	for dir := range app.site.dirs {
		// Change to each directory
		if err := os.Chdir(dir); err != nil {
			// TODO: Document this error code
			return ErrCode("1101", dir)
		}
		// Get the files in just this directory
		files, err := ioutil.ReadDir(".")
		if err != nil {
			// TODO: Document this error code
			return ErrCode("0703", dir)
		}

 		// Go through all the Markdown files and convert.
		// Start search index JSON file with opening '['
		// TODO: Add this back
		//app.DelimitIndexJSON(a.Site.SearchJSONFilePath, true)
		commaNeeded := false
		for _, file := range files {
			if !file.IsDir() && isMarkdownFile(file.Name()) {
				app.site.fileCount++
				// It's a Markdown file, not a dir or anything else.
				if commaNeeded {

					// TODO: Add error checking
					// TODO: Add this back
					// app.AddCommaToSearchIndex(app.site.SearchJSONFilePath)
					commaNeeded = false
				}
				if err = app.publishFile(filepath.Join(dir, file.Name())); err != nil {
					return ErrCode("PREVIOUS", err.Error())
				}
				commaNeeded = true
			}
		}

		// Close search index JSON file with ']'
		// TODO: Add this back
		// DelimitIndexJSON(a.Site.SearchJSONFilePath, false)

	}
	fmt.Printf("%v ", app.site.fileCount)
	if app.site.fileCount != 1 {
		fmt.Println("files")
	} else {
		fmt.Println("file")
	}

	app.Note("Project tree:\n%v", app.site.dirs)
	if app.flags.Info {
		app.info()
	}
	// Return with success code.
	return nil
}
