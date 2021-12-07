package app
import(
  "os"
  "fmt"
	"path/filepath"
	"github.com/tomcam/m/pkg/default"
)
func (app *App) createStubIndex() error {
  page := fmt.Sprintf("# Welcome to %s\nThank you for using %s.",app.Site.name, app.Site.name)
  return app.createSimplePage("index.md", "", page)
}
// Creates dir if it doesn't exist
func (app *App) createSimplePage(filename string, dir string, contents string) error {
	app.Debug("simplePage(%v, %v, %v)", filename, dir, contents)
	app.Note("simplePage(%v, %v, %v)", filename, dir, contents)
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
  app.Print("\tabout to create directory %v", dir)
	err := os.MkdirAll(dir, defaults.ProjectFilePermissions)
	if err != nil {
		return ErrCode("0412", dir)
	}

	// Get the fully qualified filename to generate
	filename = filepath.Join(dir, filename)

  app.Print("About to write file %v", filename)
	if err := writeTextFile(filename, contents); err != nil {
    // TODO: Ensure all erroc odes in this function are unique 
		return ErrCode("0413", filename)
	}
	return nil
}
