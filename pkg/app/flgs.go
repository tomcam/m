package app

import (
	"github.com/tomcam/m/pkg/default"
	"path/filepath"
	//"github.com/spf13/cobra"
)

func (app *App) addFlags() {
	app.RootCmd.PersistentFlags().StringVar(&app.cfgFile, "config", "", "config file (default is "+filepath.Join(homeDir(), "."+defaults.ProductName+".yaml)"))

	// Local flags which will only run when this command
	// is called directly, e.g.:
	app.RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
