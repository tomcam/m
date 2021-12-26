package app

import (
	"fmt"
  "reflect"
	"github.com/tomcam/m/pkg/default"
	"os"
	"path/filepath"
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
	err := os.MkdirAll(dir, defaults.ProjectFilePermissions)
	if err != nil {
		return ErrCode("0412", dir)
	}

	// Get the fully qualified filename to generate
	filename = filepath.Join(dir, filename)

	app.Debug("\tAbout to write file %v", filename)
	if err := writeTextFile(filename, contents); err != nil {
		// TODO: Ensure all erroc odes in this function are unique
		return ErrCode("0413", filename)
	}
	return nil
}

func isZero(v reflect.Value) bool {
	return !v.IsValid() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}


//func (app *App) createPageFrontMatter(filename string, dir string, frontMatter interface{}) {
// Pre: pathname may contain a directory but it's already been created
func (app *App) createPageFrontMatter(pathname string, article string, frontMatter FrontMatter) error {
  //app.Print("FrontMatter:\n%s", frontMatterToString(app.Page.FrontMatter))
  f := frontMatterToString(frontMatter)
	if err := writeTextFile(pathname, f + article); err != nil {
		return ErrCode("0410", pathname)
	}
  return nil
}

// frontMatterToString  generates the front matter 
// section of a page in "sparse" format, that is,
// without a bunch of empty fields.
// Extract only the string fields with contents
// and include those, for example, 
// FrontMatter.Theme or FrontMatter.Mode
// If nothing in the front matter is set, returns
// an empty string.
// Hmm... see https://stackoverflow.com/a/66511341
func frontMatterToString(f FrontMatter) string {
  fields := reflect.ValueOf(f)
  frontMatter := ""
  for i:= 0; i < fields.NumField(); i++ {
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
