package app

import (
//"fmt"
//"github.com/tomcam/m/pkg/default"
//"github.com/tomcam/m/pkg/app"
//"os"
//"path/filepath"
)

// Site contains configuration specific to each site, such as
// its directory location, title, publish directory,
// branding information, etc.
type Site struct {
	// Home directory for the project. All other
	// paths, such as location of publish directory,
	// graphics & javascript assets, etc., are based on
	// this location.
	// If you run:
	//   mb new site /Users/tom/html/foo
	// It would be /Users/tom/html/foo
	// If you just run:
	//   mb new site
	// It's initialized to the name of the current directory.
	Path string
}
