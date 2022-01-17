package app

import (
	"fmt"
	"github.com/tomcam/m/pkg/default"
	"os"
	"path/filepath"
	"reflect"
)

// createSubIndex() generates a simple index.md in the root
// directory.
func (app *App) createStubIndex() error {
	page := fmt.Sprintf("# Welcome to %s\nhello, world.", app.Site.name)
	if !fileExists(filepath.Join(app.Site.path, "index"+defaults.DefaultMarkdownExtension)) {
		return app.createSimplePage("index.md", "", page)
	}
	// index already exists
	return nil
}

// createSimplePage generates a page of text.
// Creates dir if it doesn't exist
// Try to use use createPageFrontMatter() but this is
// perfect for generating a simple index page.
func (app *App) createSimplePage(filename string, dir string, contents string) error {
	app.Debug("simplePage(%v, %v, %v)", filename, dir, contents)
	if filename == "" {
		return ErrCode("1037", "")
	}
	// If no folder is given, assume project root.
	// Remember Go uses Unix folder conventions even
	// under Windows
	if dir == "" {
		dir = "."
	}
	dir = filepath.Join(app.Site.path, dir)
	// Create the specified folder as a subdirectory
	// of the current project.
	// TODO: Could probably remove this
	err := os.MkdirAll(dir, defaults.ProjectFilePermissions)
	if err != nil {
		return ErrCode("0412", dir)
	}

	// Get the fully qualified filename to generate
	filename = filepath.Join(dir, filename)

	app.Debug("\tAbout to write file %v", filename)
	if err := writeTextFile(filename, contents); err != nil {
		return ErrCode("1302", filename)
	}
	return nil
}

// createPageFrontMatter generates a page located at
// fully qualified pathname (assumes the directory has been created
// before calling), with text of article and a filled-in FrontMatter
// Pre: pathname may contain a directory but it's already been created
func (app *App) createPageFrontMatter(pathname string, article string, frontMatter FrontMatter) error {
	f := frontMatterToString(frontMatter)
	// TODO: Is removing the first char like this brittle? Is there a case
	// in which there would be no path separator?
	if pathname[0:1] == string(os.PathSeparator) {
		pathname = trimFirstChar(pathname)
	}
	if err := writeTextFile(pathname, f+article); err != nil {
		// xxx
		return ErrCode("0229", pathname)
	}
	return nil
}

// frontMatterToString generates the front matter
// section of a page in "sparse" format, that is,
// without a bunch of empty fields.
// So it might create something like this if called
// from a starter. Could have even fewer
// fields; simply depends on what nonempty values
// are in the FrontMatter struct.
//
//   ---
//   Theme: hero
//   Title: Assemble
//   Sidebar: left
//   ---
//
// Extract only the string fields with contents
// and include those, for example,
// FrontMatter.Theme or FrontMatter.Mode
// If nothing in the front matter is set, returns
// an empty string.
// Hmm... see https://stackoverflow.com/a/66511341
func frontMatterToString(f FrontMatter) string {
	fields := reflect.ValueOf(f)
	frontMatter := ""
	for i := 0; i < fields.NumField(); i++ {
		k := fields.Type().Field(i).Name
		contents := structFieldByNameStrMust(f, k)
		if contents != "" {
			// TODO: stringbuilder
			frontMatter += k + ": " + contents + "\n"
		}
	}
	if frontMatter != "" {
		frontMatter = "---" + "\n" + frontMatter + "---" + "\n"
	}
	return frontMatter
}
