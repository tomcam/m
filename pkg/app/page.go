package app

import (
	"fmt"
	"github.com/tomcam/m/pkg/default"
	"os"
	"path/filepath"
)

// Page contains information about the Markdown page and
// front matter currently being processed.
type Page struct {
	// Directory location of this page
	dir string

	// TODO: Marshal in front matter as a real struct
	frontMatterRaw map[string]interface{}
	FrontMatter    FrontMatter

	// Theme used by this page.
	Theme Theme

	// In case of a nested theme, such as "debut/gallery"
	// or "debut/gallery/item", include them all here.
	themes []Theme

	// Fully qualified filename of this source file
	filePath string

	// List of stylesheets actually published
	// (for example, only sidebar-left.css
	// or sidebar-right.css will be published)
	stylesheets []string

	// List of stylesheets with full path designations and
	// enclosed in HTML stylesheet tags.
	stylesheetTags string
}

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
