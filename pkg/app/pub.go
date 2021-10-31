package app

import (
	"github.com/tomcam/m/pkg/default"
	"github.com/tomcam/m/pkg/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func (t *Theme) themeFromYaml(filename string) *Theme {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	err = yaml.Unmarshal(b, t)
	if err != nil {
		return nil
	}
	return t
}

func (app *App) publishFile(filename string) error {

	// Figure out this file's relative position in the output
	// directory true. For example:
	//   /Users/tom/code/m/cmd/mb -> /Users/tom/code/m/cmd/mb/test/test.md
	// Results in:
	//   /test
  app.Debug("publishFile(%#v)", filename)
	rel := relDirFile(app.Site.path, filename)
	app.Page.filePath = filename
	var err error
	app.readSiteFileConfig()
	if err != nil {
		app.Note("\trSiteFileConfig() failed: %v", err.Error())
		// TODO: Handle error properly & and document error code
		return err
	}
	app.Page.dir = currDir()
	// Get the fully qualified name of the destination file
	target := replaceExtension(filename, "html")
	target = filepath.Join(app.Site.publishPath, rel, filepath.Base(target))

	var body []byte
	// Convert Markdown file to a byte slice of HTML
	// Return with YAML front matter in app.Page.frontMatter
	if body, err = app.MdFileToHTML(filename); err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}

	// Read the theme configuration file (usally called
	// themename.yaml), where themename is replaced
	// by a theme directory name, such as wide or pillar
	// TODO: make sure loadTheme() looks in all correct
	// places for the theme, such as config files, not
	// just the page front matter
	app.Note("Front matter is %#v", app.Page.FrontMatter)
	app.loadTheme()
  app.Note("After loadTheme() front matter Struct is:\n%#v", app.Page.FrontMatter)
  app.Note("After loadTheme() raw front matter is:\n%#v", app.Page.frontMatterRaw)

	// Write HTML text of the body
	fullPage := app.Site.HTMLStartFile +
		app.Page.Theme.Language + ">" + "\n" +
		"<meta charset=\"utf-8\">" + "\n" +
		"<head>" +
		metatag("description", app.descriptionTag()) +
		metatag("viewport", "width=device-width,initial-scale=1") +
		app.stylesheetTags() +
		"</head>" + "\n" + "<body>" + "\n" +
		app.header() +
		app.article(body, "article") +
		app.footer() +
		app.Site.HTMLEndFile

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
	for _, stylesheet := range app.Page.Theme.Stylesheets {
		stylesheets = stylesheets +
			stylesheetTag(filepath.Join(app.Site.cssPublishPath, stylesheet))
	}
	return stylesheets
}

// descriptionTag() reads Description from front matter
// and returns as the full
// TODO: Get as []byte and also soee FullDescriptionTag
func (app *App) descriptionTag() string {
	description := app.frontMatterMust("Description")
	// TODO: Incorporiate logic from FullDescriptionTag
	return description
}

// mdFileToHTML converts the markdown file in filename to HTML.
// It may include optional front matter.
// TODO: Document that it also runs interps, which perhaps
// should be a separate step
func (app *App) MdFileToHTML(filename string) ([]byte, error) {
	// Read file into a byte slice.
	b := util.FileToBytes(filename)
	s := app.interps(filename, string(b))
	// Convert to HTML
	return app.mdToHTML([]byte(s))
}

// buildPublishDirs() creates a mirror of the source
// directory in the publish directory.
func (app *App) buildPublishDirs() error {
	for dir := range app.Site.dirs {
		// Get the relative path.
		rel := relDirFile(app.Site.path, filepath.Join(dir, "a"))
		// Join it with the publish directory.
		full := filepath.Join(app.Site.publishPath, rel)
		if err := os.MkdirAll(full, defaults.PublicFilePermissions); err != nil {
			app.Verbose("buildPublishDirs(): Unable to create path %v", full)
			// TODO: Check error handling here
			//return ErrCode("0403", app.Site.publishPath,"" )
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
	case "header", "nav", "footer":
		html = app.layoutElement(tag)
		if html != "" {
			return wrapTag("<"+tag+">", html, true)
		}
	case "sidebar":
		html = app.layoutElement(tag)
		if html != "" {
			return wrapTag("<"+tag+">", html, true)
		}
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
		l = app.Page.Theme.Header
	case "nav":
		l = app.Page.Theme.Nav
	case "sidebar":
		l = app.Page.Theme.Sidebar
	case "footer":
		l = app.Page.Theme.Footer
	}
	return app.layoutEl(l)
}

// TODO: Probably want to return a byte slice
func (app *App) layoutEl(l layoutElement) string {
	var err error
	var html []byte
	// Inline HTML is top priority
	if l.HTML != "" {
		html = []byte(l.HTML)
	}

	// No inline HTML. Get filename.
	filename := l.File

	// Locate it in the theme directory
	filename = filepath.Join(app.siteThemesPath(), filename)

	// TODO: Should probably return an error
	// Quit silently if file can't be found
	if !fileExists(filename) {
		// TODO: Handle error properly & and document error code
    // T
		app.Debug("Can't find theme file %v", filename)
		return ""
	}
	app.Debug("\t\tlayoutEl(%v", filename)

	// Convert file contents to a byte slice of HTML
	if isMarkdownFile(filename) {
		if html, err = app.MdFileToHTML(filename); err != nil {
			app.Debug("\t\tlayoutEl() ERROR converting Markdown file %v", filename)
			return ""
		} else {
			return string(html)
		}
	} else {
		html = fileToBuf(filename)
	}
	// Handle the case where pure HTML was specified. First
	// handle any Go template values.
	return app.interps(filename, string(html))
}

// article() takes the already-generated HTML and returns
// the text wrapped in an "<article>" tag.
// You can optionally include an id tag with it.
func (app *App) article(body []byte, params ...string) string {
	// interps runs custom Go template functions like ftime
	html := app.interps(app.Page.filePath, string(body))
	if len(params) < 1 {
		// Optional ID tag was not supplied
		html = wrapTag("<"+"article"+">", html, true)
	} else {
		// Use the supplied ID tag
		id := params[0]
		html = wrapTag("<"+"article"+" id=\""+id+"\""+">", html, true)
	}
	return html
}

func (app *App) header() string {
	return app.layoutElementToHTML("header")
}

func (app *App) nav() string {
	return app.layoutElementToHTML("nav")
}

func (app *App) sidebar() string {
	return app.layoutElementToHTML("sidebar")
}

func (app *App) footer() string {
	return app.layoutElementToHTML("footer")
}

// sidebarType() determines what sidebar to use,
// if any.
// If no value has been set for this page,
// it assigns the sidebar value set in
// Site.DefaultSidebar.
func (app *App) sidebarType() {
	// TODO: Maket this a cfg value
	sidebar := app.frontMatterMustLower("Sidebar")
	if sidebar == "" {
		sidebar = strings.ToLower(app.Site.Sidebar)
	}
	if sidebar == "left" || sidebar == "right" {
		app.Page.FrontMatter.Sidebar = sidebar
	}
}

// sidebar() adds a sidebar as appropriate.
// Formerly calle just sidebar()
func (a *App) addSidebar() {
	a.Note("\taddSideBar() not implemented")
	/*
		if a.FrontMatter.Sidebar == "left" || a.FrontMatter.Sidebar == "right" {
			a.appendStr(a.layoutElementToHTML(&a.Page.Theme.PageType.Sidebar, "<aside id=\"sidebar\">"))
		}
	*/
}
