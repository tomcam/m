package app
import(
	//"github.com/tomcam/m/pkg/app"
)

func (app *App) publishFile (filename string) error {
  app.Note("\t%v", filename)
  return nil
}
