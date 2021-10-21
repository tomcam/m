package app

import (
	"github.com/tomcam/m/pkg/default"
	"github.com/tomcam/m/pkg/util"
	//"github.com/tomcam/m/pkg/mark"
	"os"
	"path/filepath"
)

func (app *App) publishFile(filename string) error {

	// Figure out this file's relative position in the output
	// directory true. For example:
	//   /Users/tom/code/m/cmd/mb -> /Users/tom/code/m/cmd/mb/test/test.md
	// Results in:
	//   /test
	rel := relDirFile(app.site.path, filename)
	// Get the fully qualified name of the destination file
	target := replaceExtension(filename, "html")
	target = filepath.Join(app.site.publishPath, rel, filepath.Base(target))

	var err error
	var body []byte
	// Convert Markdown file to a byte slice of HTML
	if body, err = app.MdFileToHTML(filename); err != nil {
		return err
	}

	if err = os.WriteFile(target, body, defaults.PublicFilePermissions); err != nil {
		// TODO: Improve error handling
		return err
	}
	return nil
}

// mdFileToHTML converts the markdown file in filename to HTML.
// It may include optional front matter.
func (app *App) MdFileToHTML(filename string) ([]byte, error) {
	// Read file into a byte slice.
	s := util.FileToBytes(filename)
	// Convert to HTML
	return app.mdToHTML(s)
}

// buildPublishDirs() creates a mirror of the source
// directory in the publish directory.
func (app *App) buildPublishDirs() error {
	for dir := range app.site.dirs {
		// Get the relative path.
		rel := relDirFile(app.site.path, filepath.Join(dir, "a"))
		// Join it with the publish directory.
		full := filepath.Join(app.site.publishPath, rel)
		if err := os.MkdirAll(full, defaults.PublicFilePermissions); err != nil {
			app.Verbose("buildPublishDirs(): Unable to create path %v", full)
			// TODO: Check error handling here
			//return ErrCode("0403", app.site.publishPath,"" )
			return ErrCode("PREVIOUS", err.Error())
		}
	}
	return nil
}
