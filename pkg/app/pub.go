package app

import (
	//"bytes"
	"github.com/tomcam/m/pkg/default"
	//"github.com/tomcam/m/pkg/util"
	//"github.com/yuin/goldmark-meta"
	//"github.com/yuin/goldmark/parser"
  "fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// publish() copies the specified file to the publish directory,
// as long as it's not excluded. It's known not to
// be a Markdown file.
func (app *App) publish(filename string) error {
	rel, err := filepath.Rel(app.Site.path, filepath.Dir(filename))
	if err != nil {
		// TODO: Perhaps better error context
		return ErrCode("PREVIOUS", err.Error())
	}
	dest := filepath.Join(app.Site.publishPath, rel, filepath.Base(filename))
	app.Debug("\t\tpublish(%v) to %v", filename, rel)
	return Copy(filename, dest)
}

func (app *App) publishMarkdownFile(filename string) error {
	app.Page = Page{}
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
	// Save full pathname to file
	app.Page.filePath = filename
	// Save just the filename
	app.Page.fileBaseName = filepath.Base(filename)
	// Save the current directory
	app.Page.dir = currDir()

	// Take the input file name, e.g. myarticle.md or whatever,
	// and use it to form the fully qualified name of the
	// output file, e.g.
	// /Users/tom/mb/mysite/.pub/myarticle.html
	target := replaceExtension(filename, "html")
	target = filepath.Join(app.Site.publishPath, rel, filepath.Base(target))

	var b []byte
	b = fileToBuf(filename)
	// Convert Markdown file to a byte slice of HTML
	// The body contains the text of a file, for example, 'foo.md', and
	// the contents of its front matter. Other page layout elements, such
	// as the footer and header, will be parsed later in App.layoutEl()
	if b, err = app.mdYAMLToHTML(b); err != nil {
    // TODO: Real error code
		return err
	}
  // xxx
	// Get YAML front matter and copy to app.Page.FrontMatter
  app.frontMatterRawToStruct()

  var body string
	if body, err = app.doTemplateFuncs(filename, string(b)); err != nil {
    // TODO: Real error code
	  return err
	}

  fmt.Printf("app.metaData: %v\n", app.metaData)
  fmt.Printf("app.Page.FrontMatter: %+v\n", app.Page.FrontMatter)
  fmt.Printf("app.Page.FrontMatter.Theme: %v\n", app.Page.FrontMatter.Theme)
  fmt.Printf("app.Page.FrontMatter.List: %+v\n", app.Page.FrontMatter.List)
  fmt.Printf("app.Site: %+v\n", app.Site)

	// Theme has been named in Page.FrontMatter so load it.
	if err = app.loadTheme(); err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}
	//app.Print("after loadTheme(): %+v", app.Page.FrontMatter.List)

	// Copy out stylesheets, graphics, and other assets
	// required for this page.
	if err := app.publishPageAssets(); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}

	// Contents of all files to insert into <head> section
	var h string
	if h, err = app.headFiles(); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	var header string
	if header, err = app.header(); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	var nav string
	if nav, err = app.nav(); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	var sidebar string
	if sidebar, err = app.sidebar(); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	var footer string
	if footer, err = app.footer(); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	var closeScripts string
	if closeScripts, err = app.insertScript(app.Site.scriptClosePath); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	// Write HTML text of the body
	fullPage := "<!DOCTYPE html><html lang=" +
		"\"" + app.Site.Language + "\"" + ">" + "\n" +
		"<head>" + "\n" +
		app.titleTag() +
		app.metatags() +
		// Avoid flash of unstyled content (FOUT)
		"<style>.no-js {visibility: visible;}</style>" +
		h + "\n" +
		app.Page.stylesheetTags +
		"</head>" + "\n" +
		// Avoid flash of unstyled content (FOUT)
		"<body style='visibility: hidden;'class='no-js'>" +
		"<script>document.querySelector('body').classList.remove('no-js');</script>" +
		header +
		nav +
		app.article([]byte(body), "article") +
		sidebar +
		footer +
		closeScripts +
		"</body>" + "\n" +
		"</html>"

	if err = os.WriteFile(target, []byte(fullPage), defaults.PublicFilePermissions); err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}

	// TODO: Write out assets in same dir as page
	// TODO: If  you have a lot of source files in the directory,
	// you'd be copying them all mulitple times. That's a perf issue.

	app.Page.stylesheets = []string{}
	// TODO: May be unnecessary
	app.Page.Theme = Theme{}
	app.Site.publishedThemes = map[string]bool{}
	return nil
}

// normalizeStylesheetList() builds the list of stylesheets
// required to publish this theme in the correct order
// and with the right filenames. It transforms the raw list of
// stylesheets needed by this theme from the raw
// collection of filesheets from the theme config file
// (and stored in app.Page.Theme.Stylesheets),
// into app.Page.Theme.publishStylesheets
// Pre: loadTheme() has been called so all nested themes are present
func (app *App) normalizeStylesheetList() {
	app.Debug("\t\t\tnormalizeStylesheetList()")
	for level := 0; level < len(app.Page.Theme.levels); level++ {
		theme := app.Page.themes[level]
		app.Debug("\t\t\t\t%#v", app.Page.themes[level])
		responsive := false
		darkMode := app.darkMode()
		if darkMode {
			app.Debug("\t\t\t\tDark mode")
		}
		// Is this page light (system default) or dark mode?
		// Get the list of stylesheets specified for this theme.
		for _, stylesheet := range theme.Stylesheets {
			app.Debug("\t\t\t\t%v", stylesheet)

			switch stylesheet {
			case "sidebar-right.css":
			case "sidebar-left.css":
				stylesheet = ""
			case "theme-dark.css":
				app.Debug("\t\t\t\tConvert dark mode to light")
				if !darkMode {
					stylesheet = "theme-light.css"
				}
			case "theme-light.css":
				app.Debug("\t\t\t\t\tConvert light mode to dark")
				if darkMode {
					stylesheet = "theme-dark.css"
				}
			case "responsive.css":
				app.Debug("\t\t\t\tresponsive.css found")
				responsive = true
				stylesheet = ""
			}
			if stylesheet != "" {
				theme.publishStylesheets =
					append(theme.publishStylesheets, stylesheet)
			}
		}
		// sidebar-right.css or sidebar-left.css must be
		// penultimate, followed by responsive.css
		sidebar := app.sidebarType()
		app.Debug("\t\t\t\tsidebar type: %v", sidebar)
		var stylesheet string
		switch sidebar {
		case "left", "right":
			stylesheet = "sidebar-" + sidebar + ".css"
			theme.publishStylesheets =
				append(theme.publishStylesheets, stylesheet)
		}
		// responsive.css is the final stylesheet to add
		if responsive == true {
			app.Debug("\t\t\t\tadding responsive.css")
			theme.publishStylesheets = append(theme.publishStylesheets, "responsive.css")
		}
		app.Page.themes[level] = theme
	}
}

// descriptionTag() reads Description from front matter
// and returns as the full
// TODO: Get as []byte and also soee FullDescriptionTag
// descriptionTag() reads Description from front matter and
// if it can't find any, does whatever it can to come up
// with a worthwhile description
func (app *App) descriptionTag() string {
	description := app.frontMatterMust("Description")
	if description != "" {
		return description
	}

	// TODO: Create test case for this
	if app.Site.Branding != "" {
		return app.Site.Branding
	}
	// TODO: Create test case for this
	if app.Site.Company.Name != "" {
		return app.Site.Company.Name
	}
	// TODO: Create test case for this
	if app.Site.name != "" {
		return app.Site.name
	}
	return "Powered by " + defaults.ProductName
}

// buildPublishDirs() creates a mirror of the source
// directory in the publish directory and also adds
// paths defined at at startup.
func (app *App) buildPublishDirs() error {

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
	// TODO: Not actually using this
	app.Debug("About to create %v", app.Site.cssPublishPath)
	if err := os.MkdirAll(app.Site.cssPublishPath, defaults.PublicFilePermissions); err != nil {
		app.QuitError(err)
	}

	return nil
}

// stylesheetTag() produces just that.
// Given the full pathname of a stylesheet,
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
func (app *App) layoutElementToHTML(tag string) (string, error) {
	app.Debug("\tlayoutElementToHTML(%v)", tag)
	// tag is forced lowercase already
	html := ""
	var err error
	switch tag {
	default:
		html = ""
	case "header", "nav", "footer":
		app.Debug("\t\tAttempting to generate HTML for %v", tag)
		if html, err = app.layoutElement(tag); err != nil {
			app.Debug("\t\t%s failed", tag)
			return "", err
		}
		if html != "" {
			app.Debug("\t\tHTML: %s", html)
			return wrapTag("<"+tag+">", html, true), nil
		}
		app.Debug("\t\tNo HTML generated")
	case "sidebar":
		html, err = app.layoutElement(tag)
		if html != "" {
			return wrapTag("<aside id='sidebar'>", html, true), nil
		}
	}
	return html, nil
}

// layoutElement() takes one of the layout elements (which
// are "header", "nav", "article", "sidebar", and "footer")
// and figures out its HTML representation. According
// to the theme config it maybe in inline HTML from the
// config file, or it may be a file containing either
// markdown or HTML.
func (app *App) layoutElement(tag string) (string, error) {
	var l layoutElement
	switch tag {
	case "header":
		l = app.Page.Theme.Header
	case "nav":
		l = app.Page.Theme.Nav
	case "sidebar":
		sidebar := app.sidebarType()
		if sidebar != "left" && sidebar != "right" {
			return "", nil
		}
		l = app.Page.Theme.Sidebar
	case "footer":
		l = app.Page.Theme.Footer
	}
	app.Debug("\t\t\tlayoutElement(%v). %#v", tag, l)
	return app.layoutEl(l)
}

// TODO: Probably want to return a byte slice
func (app *App) layoutEl(l layoutElement) (string, error) {
	app.Debug("\t\t\tlayoutEl(%#v)", l)
	var html []byte
	// Inline HTML is top priority
	if l.HTML != "" {
		html = []byte(l.HTML)
	}

	// No inline HTML. Get filename.
	filename := l.File

	// If no file and no HTML have been specified, no
	// sweat. Only the article is required.
	if filename == "" {
		return "", nil
	}

	// Locate it in the theme directory
	filename = filepath.Join(app.Page.Theme.sourcePath, filename)

	app.Debug("\t\t\t\tlayoutEl filename: %v", filename)
	if !fileExists(filename) {
		return "", ErrCode("1034", filename)
	}

	var err error
	if html, err = app.mdYAMLFileToHTML(filename); err != nil {
		return "", err
	} else {
		return string(html), nil
	}
}

// article() takes the already-generated HTML and returns
// the text wrapped in an "<article>" tag.
// You can optionally include an id tag with it.
func (app *App) article(body []byte, params ...string) string {
	html := string(body)
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

func (app *App) header() (string, error) {
	// If this feature isn't supported by the
	// Metabuzz Theme Framework, don't bother.
	if !app.Page.Theme.Supports.MTF || !app.Page.Theme.Supports.Header || strings.Contains(strings.ToLower(app.Page.FrontMatter.Suppress), "header") {
		return "", nil
	}
	return app.layoutElementToHTML("header")
}

func (app *App) nav() (string, error) {
	// If this feature isn't supported by the
	// Metabuzz Theme Framework, don't bother.
	if !app.Page.Theme.Supports.MTF || !app.Page.Theme.Supports.Nav || strings.Contains(strings.ToLower(app.Page.FrontMatter.Suppress), "nav") {
		return "", nil
	}
	return app.layoutElementToHTML("nav")
}

func (app *App) sidebar() (string, error) {
	// If this feature isn't supported by the
	// Metabuzz Theme Framework, don't bother.
	if !app.Page.Theme.Supports.MTF || !app.Page.Theme.Supports.Sidebar || strings.Contains(strings.ToLower(app.Page.FrontMatter.Suppress), "sidebar") {
		return "", nil
	}
	return (app.layoutElementToHTML("sidebar"))
}

func (app *App) footer() (string, error) {
	// If this feature isn't supported by the
	// Metabuzz Theme Framework, don't bother.
	if !app.Page.Theme.Supports.MTF || !app.Page.Theme.Supports.Footer || strings.Contains(strings.ToLower(app.Page.FrontMatter.Suppress), "footer") {
		return "", nil
	}
	return app.layoutElementToHTML("footer")
}

// darkMode() returns true if dark mode has been specified
func (app *App) darkMode() bool {
	// If this feature isn't supported by the
	// Metabuzz Theme Framework, don't bother.
	if !app.Page.Theme.Supports.MTF || !app.Page.Theme.Supports.Mode {
		return false
	}
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
func (app *App) sidebarType() string {
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

// publishStylesheet copies a file specified in the theme
// configuration file to the publish directory used for
// stylesheets.
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
		return ErrCode("1027", app.Page.Theme.filename+" specifies a file named "+filepath.Base(source)+", which can't be found")
	}
	// Keep list of stylesheets that got published
	app.Page.stylesheets = append(app.Page.stylesheets, dest)
	return nil
}

// publishStylesheets() copies the stylesheets required
// by this theme to be published, omitting those it
// doesn't need (for example, "theme-light.css" if
// Mode has been set to "dark").
// Pre: It must be called
// after normalizeStylesheetList().
func (app *App) publishStylesheets() error {
	app.Debug("\t\t\tpublishStylesheets()")
	var source, dest string
	// Go through the list of stylesheets for this theme.
	// Copy stylesheets for this theme from the local
	// theme directory to the publish
	// CSS directory for stylesheets.
	var stylesheets strings.Builder
	for level := 0; level < len(app.Page.Theme.levels); level++ {
		theme := app.Page.themes[level]
		app.Debug("\t\t\t\ttheme is: %#v", theme.level)
		app.Debug("\t\t\t\tpublishStylesheets: %#v", theme.publishStylesheets)
		for _, stylesheet := range theme.publishStylesheets {
			if strings.HasPrefix(strings.ToLower(stylesheet), "http") {
				stylesheet = stylesheetTag(stylesheet)
				stylesheets.WriteString(stylesheet)
				continue
			}
			source = filepath.Join(theme.sourcePath, stylesheet)
			dest = filepath.Join(app.themePublishDir(theme.level), stylesheet)

			app.Debug("\t\t\t\t\tsheet: %#v", stylesheet)
			if err := app.publishStylesheet(source, dest); err != nil {
				return ErrCode("PREVIOUS", err.Error())
				//return ErrCode("1024", source)
			}

			stylesheet = stylesheetTag(filepath.Join(app.themePublishDir(theme.level), stylesheet))
			stylesheets.WriteString(stylesheet)

		}
		app.Page.stylesheetTags = stylesheets.String()
		app.Debug("\t\t\t\tstylesheets.String() %v", stylesheets.String())
	}
	return nil
}

// publishPageAssets() makes copies of stylesheets,
// graphics files, and other assets required to
// publish this page.
// Pre: loadTheme() has been called, so all nested stylesheets are present
func (app *App) publishPageAssets() error {
	app.Debug("\t\tpublishPageAssets() for %v", app.Page.filePath)
	// Take raw list of stylesheets from theme and ensure
	// they're in the right order, right
	app.normalizeStylesheetList()
	if err := app.publishStylesheets(); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	return nil
}

// titleTag() uses the tag specified in the front matter.
// If it can't find one it tries other ideas.
// TODO: Restore inferTitle
func (app *App) titleTag() string {
	title := app.Page.FrontMatter.Title
	if title == "" {
		title = defaults.ProductName + ": Title needed here, squib"
	}
	return wrapTag("<title>", title, true)
}

func (app *App) metatags() string {
	return ("<meta charset=\"utf-8\">" + "\n" +
		metatag("description", app.descriptionTag()) +
		metatag("viewport", "width=device-width,initial-scale=1") +
		metatag("generator", defaults.ProductBranding))
}

// headFiles() (formerly headerFiles) finds
// all the files in the headers subdirectory
// and copies them into the HMTL head section
func (app *App) headFiles() (string, error) {
	var h string
	headers, err := ioutil.ReadDir(app.Site.headTagsPath)
	if err != nil {
		return "", ErrCode("0706", app.Site.headTagsPath)
	}
	for _, file := range headers {
		h += fileToString(filepath.Join(app.Site.headTagsPath, file.Name()))
	}
	return h, nil
}

// insertScript() injects Javascript (technically,
// any thing inside script tags) into the
// output stream.
// dir is the fully qualified directory
// name containing the scripts.
// The scripts supplied MUST provide their
// own script tags.
func (app *App) insertScript(dir string) (string, error) {
	var script string
	scripts, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", ErrCode("0709", dir)
	}
	for _, file := range scripts {
		script += fileToString(filepath.Join(dir, file.Name()))
	}
	return script, nil
}
