package app

import (
	//"text/transform"
	"bytes"
	//"fmt"
	//"github.com/tomcam/m/pkg/default"
	//toc "github.com/abhinav/goldmark-toc"
  "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"io/ioutil"
	"os"
	"path/filepath"
)

// mdToHTML converts a Markdown source file in a byte
// slice to HTML.
func (app *App) mdToHTML(source []byte) ([]byte, error) {
  //app.src = source
	var buf bytes.Buffer
	if err := app.parser.Convert(source, &buf, parser.WithContext(app.parserCtx)); err != nil {
		return buf.Bytes(), ErrCode("0920", err.Error())
	}
	// Obtain the parsed front matter as a raw
	// interface
	app.Page.frontMatterRaw = meta.Get(app.parserCtx)
	return buf.Bytes(), nil
}

// build() is what it's all about!
// It converts the project to HTML files.
// pathname isn't known to be good. It's
// for situation such as mb build ~/foo
// when you happen to be in
// directory ~/something/else/bar
// Pre:
//  setPaths()
func (app *App) build(path string) error {
	var err error
	app.Debug("\tbuild(%v)", path)
	// Change to specified directory.
	// Update app.Site.path and build all related directories
	if err := app.setWorkingDir(path); err != nil {
		app.Debug("\t\tUnable to change to directory (%v)", path)
		return ErrCode("1107", path)
	}

	if !isProject(app.Site.path) {
		return ErrCode("1002", path)
	}

	// Get a list of all files & directories in the site.
	if _, err = app.getProjectTree(app.Site.path); err != nil {
		app.Debug("\t\tError reading project tree")
		return ErrCode("0913", app.Site.path)
	}

	// Delete any existing publish dir
	if err := os.RemoveAll(app.Site.publishPath); err != nil {
		app.Debug("\t\tError deleting existing publish dir %v", app.Site.publishPath)
		return ErrCode("0302", app.Site.publishPath)
	}

	// xxx
	if err := app.readSiteConfig(); err != nil {
		app.Debug("\t\tError reading siteConfig %v", app.Site.siteFilePath)
		return ErrCode("PREVIOUS", err.Error())
	}
	// Build the target publish dir so there should be
	// no trouble copying files over
	app.buildPublishDirs()

	// Loop through the list of permitted directories for this site.
	for dir := range app.Site.dirs {
		// Change to each directory
		if err := os.Chdir(dir); err != nil {
			// TODO: Handle error properly & and document error code
			return ErrCode("1101", dir)
		}
		// Get the files in just this directory
		files, err := ioutil.ReadDir(".")
		if err != nil {
			return ErrCode("0703", dir)
		}

		// https://github.com/tomcam/mb/blob/master/pkg/app/pub.go
		// Go through all the Markdown files and convert.
		// Start search index JSON file with opening '['
		// TODO: Add this back
		// TODO: I think this will be superseded by Yuin's toc feature
		//app.DelimitIndexJSON(a.Site.SearchJSONFilePath, true)
		filename := ""
		commaNeeded := false
		for _, file := range files {
			if !file.IsDir() {
				if isMarkdownFile(file.Name()) {
					// It's a Markdown file
					app.Site.fileCount++
					// It's a Markdown file, not a dir or anything else.
					if commaNeeded {

						// TODO: Add this back when I add search
						// app.AddCommaToSearchIndex(app.Site.SearchJSONFilePath)
						commaNeeded = false
					}
					filename = filepath.Join(dir, file.Name())
					if err = app.publishMarkdownFile(filename); err != nil {
						return ErrCode("PREVIOUS", err.Error())
					}
					commaNeeded = true
				} else {
					// It's not a Markdown file. Copy if it's a graphic
					// asset or something.
					filename = filepath.Join(dir, file.Name())
					if err = app.publish(filename); err != nil {
						return ErrCode("PREVIOUS", err.Error())
					}
				}
			}
		}

		// Close search index JSON file with ']'
		// TODO: Add this back
		// DelimitIndexJSON(a.Site.SearchJSONFilePath, false)

	}
	if app.Site.fileCount != 1 {
		app.Print("%v files", app.Site.fileCount)
	} else {
		app.Print("1 file")
	}

	// Return with success code.
	return nil
}

// newGoldmark returns the a goldmark object with a parser and renderer.
func (app *App) newGoldmark() goldmark.Markdown {

	exts := []goldmark.Extender{

		// YAML support
		meta.Meta,
		// Support GitHub tables
		extension.Table,
		extension.GFM,
		extension.DefinitionList,
		extension.Footnote,
		highlighting.NewHighlighting(
			highlighting.WithStyle(app.Site.markdownOptions.HighlightStyle),
			highlighting.WithFormatOptions()),
	}

	parserOpts := []parser.Option{
		parser.WithAttribute(),
		parser.WithAutoHeadingID()}

	renderOpts := []renderer.Option{
		// WithUnsafe is required for HTML templates to work properly
		html.WithUnsafe(),
		html.WithXHTML(),
	}
	return goldmark.New(
		goldmark.WithExtensions(exts...),
		goldmark.WithParserOptions(parserOpts...),
		goldmark.WithRendererOptions(renderOpts...),
	)
}

func (a *App) markdownAST(input []byte) ast.Node {
	ctx := parser.NewContext()
	p := a.newGoldmark().Parser()
	return p.Parse(text.NewReader(input), parser.WithContext(ctx))
	//return p.Parse(text.NewReader(input), parser.WithContext(ap.parserCtx))
}




