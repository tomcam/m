package app

import (
	"github.com/tomcam/m/pkg/default"
	"os"
	"path/filepath"
)

type description struct {
	filename string
	dir      string
	embedded embedded
	mdtext   string
}

type embedded struct {
	filename string
	contents string
}

/* SVG file of an exciting 100x100px gray box */
var svgFile = `<?xml version="1.0" encoding="utf-8"?>
<svg version="1.1"
     baseProfile="full"
     width="100" height="100"
     xmlns="http://www.w3.org/2000/svg">
  <rect x="0" y="0" width="100" height="100" style="fill: rgb(216, 216, 216);"/>
</svg>
`

var (
	siteTest = []description{
		{"index.md",
			"",
			embedded{filename: "box-100x100.svg", contents: svgFile},
			`# Home
Go [one level deep](one/index.html), [two levels deep](two/three/index.html)
* Host: {{ hostname }}
* Time: {{ ftime }}
* Location of this file: {{ path }}
* Theme {{ .Page.FrontMatter.Theme }}
#### Code fence
` +
				"```" +
				`
print("hello, world.")
` +
				"```" +
				`
**SVG Box**
![100x100 SVG box](box-100x100.svg)
`},

		{"index.md",
			"one",
			embedded{filename: "", contents: ""},
			`
---
Theme: pillar 
---
# Page 1
This page is 1 level deep.
The time is {{ ftime }}
`,
		},

		{"pagetype.md",
			"one",
			embedded{filename: "", contents: ""},
			`
---
Theme: debut/gallery
---
# PageType test
This uses the theme named {{ .Page.FrontMatter.Theme }}
The time is {{ ftime }}
`,
		},

		{"index.md",
			"two/three",
			embedded{filename: "box-100x100.svg", contents: svgFile},
			`
# Page 2
This page is 2 levels deep.
Location of this file: {{ path }}
**Box**
![100x100 SVG box](box-100x100.svg)
Go [home 1](/index.html)
Go [home 2](\/index.html)
Go [home 3](/)
Go [home 4](/./index.html)
`},
	}

	/* Directory structure for the test site */
	testDirs = [][]string{
		// Don't forget the svg file is copied into directories 0 and 2,
		// in other words, those named  "one" and "three"
		{"one"},
		{"two", "three"},
	}
)

// writeSiteFromArray() takes an array of
// structures containing a filename,
// a path to that filename, and the markdown
// text itself, and writes them out to
// a test site. It may also contain nonzero
// values for embedded, which is the contents
// of a single embedded value written out as
// its own file, for example, an SVG image.
func writeSiteFromArray(sitename string, site []description) error {
	for _, f := range site {
		// Write out the Markdown page in its directory.
		path := filepath.Join(f.dir, f.filename)
		err := writeTextFile(path, f.mdtext)
		if err != nil {
			return ErrCode("PREVIOUS", err.Error(), path)
		}
		// If there's an embedded file, write it out
		if f.embedded.filename != "" {
			path := filepath.Join(f.dir, f.embedded.filename)
			err := writeTextFile(path, f.embedded.contents)
			if err != nil {
				return ErrCode("PREVIOUS", err.Error(), path)
			}
		}
	}
	return nil
}

// kitchenSink() Generates a test site from an
// array of structures containing a filename,
// a path to that filename, and the markdown
// text itself.
// TODO: Why not just us newSite()? I must have had some
// reason but this is asking for trouble.
// Consider removing altogether now that I have
// starters working, though that would mean embeddding
// an image file whereas now the .SVG file is generated
// on the fly.
func (app *App) kitchenSink(pathname string) error {

	var err error
	// Create a project at the specified path
	err = os.MkdirAll(pathname, defaults.ProjectFilePermissions)
	if err != nil {
		return ErrCode("0401", pathname)
	}
	// Update app.Site.path and build all related directories
	if err := app.setWorkingDir(pathname); err != nil {
		return err
	}

	// Change to specified directory.
	pathname = app.Site.path

	// Create directory structure for test site
	if err := createDirStructure(&testDirs); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}

	// Create minimal directory structure: Publish directory
	// .site directory, .themes, etc.
	if err := createDirStructure(&defaults.SitePaths); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	// Build the site from the array of data structures
	if err := writeSiteFromArray(pathname, siteTest); err != nil {
		// TODO: Improve error handling
		app.QuitError(err)
		return err
	}

	// Get factory themes and copy to project. They will then
	// be copied on demand to the publish directory as needed.
	// This makes it easy to find themes and modify theme.
	if err := app.copyMbDir(); err != nil {

		return ErrCode("PREVIOUS", err.Error())
	}

	if err := app.writeSiteConfig(""); err != nil {
		app.Note("Error writing site file %v", app.Site.siteFilePath)
		return ErrCode("PREVIOUS", err.Error())
	}
	app.Print("Created site %v", app.Site.path)
	return nil

}
