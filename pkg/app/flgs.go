package app

import (
	"github.com/tomcam/m/pkg/default"
	"os"
	"path/filepath"
)

func (app *App) addFlags() {
	// TODO: Make sure this comports with reality
	app.RootCmd.PersistentFlags().StringVar(&app.cfgFile, "config", "", "config file (default is "+filepath.Join(homeDir(), ".", defaults.ConfigFilename)+")")
	app.RootCmd.PersistentFlags().BoolVarP(&app.Flags.Verbose, "verbose", "v", false, "verbose output")
	app.RootCmd.PersistentFlags().BoolVarP(&app.Flags.Debug, "debug", "d", false, "show debug output")
	app.RootCmd.PersistentFlags().BoolVarP(&app.Flags.Info, "info", "i", false, "Show info after "+os.Args[0]+" runs")
	app.RootCmd.PersistentFlags().BoolVarP(&app.Flags.InfoVerbose, "info-verbose", "b", false, "Show info after "+os.Args[0]+" runs with full path information")
	app.RootCmd.PersistentFlags().BoolVarP(&app.Flags.InfoFrontMatter, "front", "f", false, "show front matter")
}
