package app

import (
	"github.com/tomcam/m/pkg/default"
	"github.com/tomcam/m/pkg/util"
	//"strings"
	//"github.com/tomcam/m/pkg/mark"
	//"json"
	"fmt"
	"os"
	"path/filepath"
)

func (app *App) publishFile(filename string) error {

	// Figure out this file's relative position in the output
	// directory true. For example:
	//   /Users/tom/code/m/cmd/mb -> /Users/tom/code/m/cmd/mb/test/test.md
	// Results in:
	//   /test
	rel := relDirFile(app.site.path, filename)
	// Get the fully qualified name of the destination file
	target := replaceExtension(filename, "html")
	target = filepath.Join(app.site.publishPath, rel, filepath.Base(target))

	var err error
	var body []byte
	// Convert Markdown file to a byte slice of HTML
	// Return with YAML front matter in app.page.frontMatter
	if body, err = app.MdFileToHTML(filename); err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}

	// Read the theme configuration file (usally called
	// themename.yaml, where themename is replaced
	// by a theme directory name, such as wide or pillar
	// TODO: make sure loadTheme() looks in all correct
	// places for the theme, such as config files, not
	// just the page front matter
	app.loadTheme()

	// Write HTML text of the body
	fullPage := app.page.theme.HTMLStartFile.HTML +
		app.page.theme.Language + ">" + "\n" +
		"<meta charset=\"utf-8\">" + "\n" +
		"<head>" +
		metatag("description", app.descriptionTag()) +
		metatag("viewport", "width=device-width,initial-scale=1") +
		app.stylesheetTags() +
		"</head>" + "\n" + "<body>" + "\n" +
		app.header() +
		string(body) +
		// xxx"<article id=\"article\">"
		"</body>" + "\n" + "</html>" + "\n"

	if err = os.WriteFile(target, []byte(fullPage), defaults.PublicFilePermissions); err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}
	// app.Note("Description: %v", app.descriptionTag())

	return nil
}

// stylesheetTags generates all stylesheet tags at once
func (app *App) stylesheetTags() string {
	stylesheets := ""
	for _, stylesheet := range app.page.theme.Stylesheets {
		stylesheets = stylesheets +
			stylesheetTag(filepath.Join(app.site.cssPublishPath, stylesheet))
	}
	return stylesheets
}

// descriptionTag() reads Description from front matter
// and returns as the full
// TODO: Get as []byte and also soee FullDescriptionTag
func (app *App) descriptionTag() string {
	if app.page.frontMatterRaw["Description"] != nil {
		return fmt.Sprint(app.page.frontMatterRaw["Description"])
	}
	return ""
}

// mdFileToHTML converts the markdown file in filename to HTML.
// It may include optional front matter.
// TODO: Do I jactually use this?
func (app *App) MdFileToHTML(filename string) ([]byte, error) {
	// Read file into a byte slice.
	s := util.FileToBytes(filename)
	// Convert to HTML
	return app.mdToHTML(s)
}

// buildPublishDirs() creates a mirror of the source
// directory in the publish directory.
func (app *App) buildPublishDirs() error {
	for dir := range app.site.dirs {
		// Get the relative path.
		rel := relDirFile(app.site.path, filepath.Join(dir, "a"))
		// Join it with the publish directory.
		full := filepath.Join(app.site.publishPath, rel)
		if err := os.MkdirAll(full, defaults.PublicFilePermissions); err != nil {
			app.Verbose("buildPublishDirs(): Unable to create path %v", full)
			// TODO: Check error handling here
			//return ErrCode("0403", app.site.publishPath,"" )
			return ErrCode("PREVIOUS", err.Error())
		}
	}
	return nil
}

// descriptionTag() does everything it can to
// generate a Description metatag for the file.
// TODO: Add this my ghetto description function
func (a *App) FulldescriptionTag() {
	/*
		// Best case: user supplied the description in the front matter.
		if a.FrontMatter.Description != "" {
			a.Page.descriptionTag = a.FrontMatter.Description
		} else if a.Site.Branding != "" {
			a.Page.descriptionTag = a.Site.Branding
		} else if a.Site.Name != "" {
			a.Page.descriptionTag = a.Site.Name
		} else {
			a.Page.descriptionTag = "Powered by " + defaults.ProductName
		}
	*/
}

// stylesheetTag() produces just that.
// Given the name of a stylesheet, like say "markdown.css",
// return it in a link tag.
func stylesheetTag(stylesheet string) string {
	// If no stylesheet specified just return empty string
	if stylesheet == "" {
		return ""
	}
	return `<link rel="stylesheet" href="` + stylesheet + `">` + "\n"
}

// layoutElementToHTML() takes an page region (header, nav, article, sidebar, or footer)
// and converts it to HTML. All we know is that it's been specified
// but we don't know whether's a Markdown file, inline HTML, whatever.
func (app *App) layoutElementToHTML(tag string) string {
	html := ""
	switch tag {
	default:
		html = ""
		// TODO: Consider logging this error or something.
	case "header":
		html = app.layoutElement(tag)
	}
	return html
}

// layoutElement() takes one of the layout elements (which
// are "header", "nav", "article", "sidebar", and "footer")
// and figures out its HTML representation. According
// to the theme config it maybe in inline HTML from the
// config file, or it may be a file containing either
// markdown or HTML.
func (app *App) layoutElement(tag string) string {
	var l layoutElement
	switch tag {
	case "header":
		l = app.page.theme.Header
		return app.layoutEl(l)
	}
	return ""
}

// siteThemesPath() determines the directory a
// theme file is found it.
//func (app *App) siteThemesPath(theme string) string() {
func (app *App) siteThemesPath() string {
  return filepath.Join(app.site.siteThemesPath, app.page.theme.Name)
}

// TODO: Probably want to return a byte slice
func (app *App) layoutEl(l layoutElement) string {
	// Inline HTML is top priority
	if l.HTML != "" {
		return l.HTML
	}

  // No inline HTML. Get filename.
  filename := l.File

  // Locate it in the theme directory
  //filename = filepath.Join(app.site.siteThemesPath, filepath.Base(filename))
  //filename = filepath.Join(app.siteThemesPath(), filepath.Base(filename))
  filename = filepath.Join(app.siteThemesPath(), filename)

  // TODO: Should probably return an error
  // Quit silently if file can't be found
  if !fileExists(filename) {
    app.Debug("Can't find theme file %v", filename)
    return ""
  }


  var err error
	var html []byte
	// Convert file contents to a byte slice of HTML
  if isMarkdownFile(filename) {
    if html, err = app.MdFileToHTML(filename); err != nil {
      // TODO: Handle error properly & and document error code
      return ""
    }
  } else {
    // html = fileToBuf(filename)
    html = fileToBuf(filename)
  }
	return string(html)
}

func (app *App) header() string {
	return app.layoutElementToHTML("header")
}

//layoutElement struct {
//	HTML string `yaml:"HTML"`
// File string `yaml:"File"`
