package app

import (
	//"fmt"
	//"github.com/tomcam/m/pkg/default"
	//"github.com/tomcam/m/pkg/errs"
	//"os"
	//"path/filepath"
)

// Site contains configuration specific to each site, such as
// its title, publish directory, and branding string.
type Site struct {
  // Home directory for the project. 
  // If you run: 
  //   mb new site /Users/tom/html/foo
  // It would be /Users/tom/html/foo
  // If you just run:
  //   mb new site
  // It's initialized to the name of the current directory
  Path string
}

func (a *App) setSiteDefaults() {
}

