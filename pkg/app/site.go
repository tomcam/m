package app

import (
	//"fmt"
	"github.com/tomcam/m/pkg/default"
	"github.com/tomcam/m/pkg/errs"
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


func (s *Site) New() error {
  // Create minimal directory structure: Publish directory
	// .site directory, .themes, etc.
  if err := createDirStructure(&defaults.SitePaths); err != nil {
		//return errs.ErrCode("PREVIOUS", err.Error())
		return errs.ErrCode("PREVIOUS", err.Error())
	}
  return nil
}


