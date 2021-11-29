package app

import (
	"github.com/tomcam/m/pkg/default"
	"github.com/tomcam/m/pkg/util"
	"os"
	"path/filepath"
	"strings"
)

// publish() copies the specified file to the publish directory,
// as long as it's not excluded.
func (app *App) publish(filename string) error {
	rel, err := filepath.Rel(app.Site.path, filepath.Dir(filename))
	if err != nil {
		// TODO: Perhaps better error context
		return ErrCode("PREVIOUS", err.Error())
	}
	app.Page.filePath = filename
	dest := filepath.Join(app.Site.publishPath, rel, filepath.Base(filename))
	//app.Debug("\tpublish(%v) to %v", filename, dest)
	app.Debug("\tpublish(%v) to %v", filename, rel)
	err = Copy(filename, dest)
	if err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	return nil
}

func (app *App) publishMarkdownFile(filename string) error {
	app.Debug("\tpublishMarkdownFile(%#v)", filename)
	// Figure out this file's relative position in the output
	// directory true. For example:
	//   /Users/tom/code/m/cmd/mb -> /Users/tom/code/m/cmd/mb/test/test.md
	// Results in:
	//   /test
	rel, err := filepath.Rel(app.Site.path, filepath.Dir(filename))
	if err != nil {
		// TODO: Perhaps better error context
		return ErrCode("PREVIOUS", err.Error())
	}
	app.Page.filePath = filename
	//var err error
	// Obtain site configuration from site.yaml
	//app.readSiteConfig()
	if err != nil {
		app.Debug("\trSiteFileConfig() failed: %v", err.Error())
		// TODO: Handle error properly & and document error code
		return err
	}
	app.Page.dir = currDir()

	// Take the input file name, e.g. myarticle.md or whatever,
	// and use it to form the fully qualified name of the
	// output file, e.g.
	// /Users/tom/mb/mysite/.pub/myarticle.html
	target := replaceExtension(filename, "html")
	target = filepath.Join(app.Site.publishPath, rel, filepath.Base(target))

	var body []byte
	// TODO: this looks like the right place for these initializations
	app.Page.Theme.stylesheetsAllLevels = make(map[string][]string)
	app.Page.stylesheets = make(map[string][]string)
	app.Page.Themes = make(map[string]Theme)

	// Convert Markdown file to a byte slice of HTML
	// Return with YAML front matter in app.Page.frontMatter
	if body, err = app.MdFileToHTML(filename); err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}

	app.Page.FrontMatter = FrontMatter{}
	// Convert the FrontMatter map produced by Goldmark into
	// the Page.FrontMatter struct.
	app.frontMatterRawToStruct()

	// Theme has been named in Page.FrontMatter so load it.
	if err = app.loadTheme(); err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}

	// Copy out stylesheets, graphics, and other assets
	// required for this page.
	if err := app.publishPageAssets(); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}

	// Write HTML text of the body
	fullPage := app.Site.HTMLStartFile +
		"\"" + app.Site.Language + "\"" + ">" + "\n" +
		"<meta charset=\"utf-8\">" + "\n" +
		"<head>" +
		metatag("description", app.descriptionTag()) +
		metatag("viewport", "width=device-width,initial-scale=1") +
		metatag("generator", defaults.ProductBranding) +
		app.stylesheetTags() +
		app.header() +
		app.nav() +
		app.article(body, "article") +
		app.sidebar() +
		app.footer() +
		app.Site.HTMLEndFile

	if err = os.WriteFile(target, []byte(fullPage), defaults.PublicFilePermissions); err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}

	// TODO: Write out assets in same dir as page
	// TODO: If  you have a lot of source files in the directory,
	// you'd be copying them all mulitple times. That's a perf issue.

	//app.Page.stylesheets = []string{}
	// TODO: May these reinitializations may be unnecessary
	app.Page.Theme = Theme{}
	//app.Page.Theme.stylesheetsAllLevels= map[string][]string{}
	app.Site.publishedThemes = map[string]bool{}
	return nil
}

func (app *App) normalizeStylesheet(stylesheet string, responsive *bool) {
	darkMode := app.darkMode()
	switch stylesheet {
	case "sidebar-right.css":
	case "sidebar-left.css":
		stylesheet = ""
	case "theme-dark.css":
		if !darkMode {
			stylesheet = "theme-light.css"
		}
	case "theme-light.css":
		if darkMode {
			stylesheet = "theme-dark.css"
		}
	case "responsive.css":
		*responsive = true
		stylesheet = ""
	}
	if stylesheet != "" {
		app.Page.Theme.stylesheetList =
			append(app.Page.Theme.stylesheetList, stylesheet)
	}
}

// TODO: I think I need to change "level" to something else. It begins
// life as an array index but later in the program, as here, it's a
// map key.
func (app *App) addPublishedStylesheet(level string, stylesheet string, responsive *bool) {
	darkMode := app.darkMode()
	switch stylesheet {
	case "sidebar-right.css":
	case "sidebar-left.css":
		stylesheet = ""
	case "theme-dark.css":
		if !darkMode {
			stylesheet = "theme-light.css"
		}
	case "theme-light.css":
		if darkMode {
			stylesheet = "theme-dark.css"
		}
	case "responsive.css":
		*responsive = true
		stylesheet = ""
	}
	if stylesheet != "" {
		app.Page.stylesheets[level] =
			append(app.Page.stylesheets[level], stylesheet)
	}
}

// normalizeStylesheetList() builds the list of stylesheets
// required to publish this theme in the correct order
// and with the right filenames. It transforms the raw list of
// stylesheets needed by this theme from the raw
// collection of filesheets from the theme config file
// into app.Page.stylesheets.
func (app *App) normalizeStylesheetList() {
	app.Note("\t\t\tnormalizeStylesheetList(): %v", app.Page.Theme.stylesheetsAllLevels)
	//for _, level := range app.Page.Theme.levels {
  for level, _ := range app.Page.Theme.stylesheetsAllLevels {
		responsive := false
		app.Debug("\t\t\t\t%v", level)
		for _, stylesheet := range app.Page.Theme.stylesheetsAllLevels[level] {
			app.addPublishedStylesheet(level, stylesheet, &responsive)
			app.Debug("\t\t\t\t\t%v", stylesheet)
		}
		// sidebar-right.css or sidebar-left.css must be
		// penultimate, followed by responsive.css
		sidebar := app.Page.FrontMatter.Sidebar
		var stylesheet string
		switch sidebar {
		case "left", "right":
			stylesheet = "sidebar-" + sidebar + ".css"
			app.Page.stylesheets[level] =
				append(app.Page.stylesheets[level], stylesheet)
		}
		// responsive.css is the final stylesheet to add
		if responsive == true {
			app.Page.stylesheets[level] =
				append(app.Page.stylesheets[level], "responsive.css")
		}
	}
	app.Note("\t\t\t\tPage stylesheets: %v", app.Page.stylesheets)
}

// stylesheetTags() returns the normalized list of
// stylesheets as a string.
func (app *App) stylesheetTags() string {
	app.Print("\t\t\tstylesheetTags()")
	var stylesheets strings.Builder
  // xxx
  for name, _ := range app.Page.stylesheets {
		for _, stylesheet := range app.Page.stylesheets[name] {
		  stylesheet = stylesheetTag(filepath.Join(app.themePublishDir(name), stylesheet))
		  stylesheets.WriteString(stylesheet)
      //app.Note("\t\t\t\t\t%v", stylesheet)
		}
  }
 	return stylesheets.String()
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
// directory in the publish directory and also adds
// paths defined at at startup.
// TODO: Not using this yet?
func (app *App) buildPublishDirs() error {
  return nil
	// Some directories are determined at startup, so add those now

	for dir := range app.Site.dirs {
		// Get the relative path.
		rel, err := filepath.Rel(app.Site.path, dir)
		if err != nil {
			return ErrCode("0501", err.Error())
		}

		// Join it with the publish directory.
		full := filepath.Join(app.Site.publishPath, rel)
		if err := os.MkdirAll(full, defaults.PublicFilePermissions); err != nil {
			app.Verbose("buildPublishDirs(): Unable to create path %v", full)
			// TODO: Check error handling here
			return ErrCode("PREVIOUS", err.Error())
		}
	}
  app.Note("buildPublishDirs(): About to create %v",app.Site.cssPublishPath)
	if err := os.MkdirAll(app.Site.cssPublishPath, defaults.PublicFilePermissions); err != nil {
		app.QuitError(err)
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
	case "header", "nav", "footer":
		html = app.layoutElement(tag)
		if html != "" {
			return wrapTag("<"+tag+">", html, true)
		}
	case "sidebar":
		html = app.layoutElement(tag)
		if html != "" {
			return wrapTag("<aside id='sidebar'>", html, true)
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
		sidebar := app.Page.FrontMatter.Sidebar
		if sidebar != "left" && sidebar != "right" {
			return ""
		}
		l = app.Page.Theme.Sidebar
	case "footer":
		l = app.Page.Theme.Footer
	}
	return app.layoutEl(l)
}

// TODO: Probably want to return a byte slice
func (app *App) layoutEl(l layoutElement) string {
	app.Debug("\t\tlayoutEl(%#v)", l)
	var err error
	var html []byte
	// Inline HTML is top priority
	if l.HTML != "" {
		html = []byte(l.HTML)
	}

	// No inline HTML. Get filename.
	filename := l.File

	// Locate it in the theme directory
	filename = filepath.Join(app.Site.siteThemesPath, filename)

	// TODO: Should probably return an error
	// Quit silently if file can't be found
	if !fileExists(filename) {
		// TODO: Handle error properly & and document error code
		// T
		app.Debug("\t\t\tlayoutEl() Can't find layout element file %v", filename)
		return ""
	}

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
	//html := app.interps(app.Page.filePath, string(body) + app.sidebar())
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
	return (app.layoutElementToHTML("sidebar"))
}

func (app *App) footer() string {
	return app.layoutElementToHTML("footer")
}

// darkMode() returns true if dark mode has been specified
func (app *App) darkMode() bool {
	// TODO: Maket this a cfg value, bcause like Theme it can also be
	// set in other areas
	if strings.ToLower(app.Page.FrontMatter.Mode) == "dark" {
		return true
	}
	return false
}

// sidebarType() determines what sidebar to use,
// if any. Returns either "left" or "right",
// forced to lowercase, or "none" if there isn't any.
// It's then written back to
// app.Page.FrontMatter.Sidebar
// TODO: If no value has been set for this page,
// it assigns the sidebar value set in
// Site.Sidebar.
func (app *App) isidebarType() string {
	// TODO: Make this a cfg value, bcause like Theme it can also be
	// set in other areas
	sidebar := strings.ToLower(app.Page.FrontMatter.Sidebar)
	if sidebar != "left" && sidebar != "right" {
		sidebar = "none"
	}
	//app.Debug("\t\t\t\tsidebarType(%v)", sidebar)
	app.Page.FrontMatter.Sidebar = sidebar
	return sidebar
}

func (app *App) publishStylesheet(source string, dest string) error {
	app.Debug("\t\t\t\tpublishStylesheet(%v, %v)", source, dest)
	if source == dest {
		return ErrCode("0217", source)
	}
	if source == "" {
		return ErrCode("1004", "")
	}
	if dest == "" {
		return ErrCode("1005", source)
	}
	err := Copy(source, dest)
	if err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	return nil
}

// publishStylesheets() copies the stylesheets required
// by this theme to be published.
// It must be called
// after normalizeStylesheetList().
func (app *App) publishStylesheets() error {
	app.Debug("\t\t\tpublishStylesheets()")
	var source, dest string
	for name, theme := range app.Page.Themes {
		for _, stylesheet := range theme.Stylesheets {
			source = filepath.Join(app.siteThemesDir(name), stylesheet)
			dest = filepath.Join(app.themePublishDir(name), stylesheet)
			app.Debug("\t\t\t\t\t%v", source)
			app.Debug("\t\t\t\t\t%v", dest)
			if err := app.publishStylesheet(source, dest); err != nil {
				return ErrCode("PREVIOUS", err.Error())
			}
		}
	}
	return nil
} // publishStylesheets

// publishPageAssets() makes copies of stylesheets,
// graphics files, and other assets required to
// publish this page.
func (app *App) publishPageAssets() error {
	app.Debug("\t\tpublishPageAssets()")
	// Take raw list of stylesheets from theme and obtain
	// a list of only those that need to be published,
	// in the correct sequence.
	app.normalizeStylesheetList()
	if err := app.publishStylesheets(); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	return nil

}
