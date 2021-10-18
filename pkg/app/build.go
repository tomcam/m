package app

import (
	"bytes"
	"fmt"
	"github.com/tomcam/m/pkg/default"
	"github.com/tomcam/m/pkg/mdext"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"io/ioutil"
	"os"
	"path/filepath"
)

// mdToHTML converts a Markdown source file in a byte
// slice to HTML.
func (app *App) mdToHTML(source []byte) []byte {
	node := app.parser.Parser().Parse(text.NewReader(source), parser.WithContext(app.parserCtx))
	// Create variable-sized buffer for parsed output.
	var buf bytes.Buffer
	// Convert the Markdown file to HTML.
	if err := app.parser.Renderer().Render(&buf, source, node); err != nil {
		// TODO: update error handling
		// a.QuitError(errs.ErrCode("0920", err.Error()))
		panic("Parse error")
		// TODO: Hmmm... this should probably return an error
		return nil
	}
	// Return the HTML.
	return buf.Bytes()
}

func (app *App) build(pathname string) error {
	var err error
	if pathname != "" {
		// Change to the specified directory.
		if err = os.Chdir(pathname); err != nil {
			return ErrCode("0901", err.Error())
		}
	}

	// Determine current fully qualified directory location.
	// Can't use relative paths internally.
	pathname = currPath()

	if !isProject(pathname) {
		return ErrCode("1002", "")
	}
	// Changed directory successfully so
	// pass it to initialize the site and update internally.
	app.setSiteDefaults(pathname)

	// Create minimal directory structure: Publish directory,
	// site directory, .themes, etc.
	if err = createDirStructure(&defaults.SitePaths); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}

	// Get a list of all files & directories in the site.
	if _, err = app.getProjectTree(app.site.path); err != nil {
		return ErrCode("0913", app.site.path)
	}

	// Delete any existing publish dir
	if err := os.RemoveAll(app.site.publishPath); err != nil {
		return ErrCode("0302", app.site.publishPath)
	}

	// Build the target publish dir so there should be
	// no trouble copying files over
	app.buildPublishDirs()

	// Loop throuIh the list of permitted directories for this site.
	for dir := range app.site.dirs {
		// Change to each directory
		if err := os.Chdir(dir); err != nil {
			// TODO: Document this error code
			return ErrCode("1101", dir)
		}
		// Get the files in just this directory
		files, err := ioutil.ReadDir(".")
		if err != nil {
			// TODO: Document this error code
			return ErrCode("0703", dir)
		}

		// Go through all the Markdown files and convert.
		// Start search index JSON file with opening '['
		// TODO: Add this back
		//app.DelimitIndexJSON(a.Site.SearchJSONFilePath, true)
		commaNeeded := false
		for _, file := range files {
			if !file.IsDir() && isMarkdownFile(file.Name()) {
				app.site.fileCount++
				// It's a Markdown file, not a dir or anything else.
				if commaNeeded {

					// TODO: Add error checking
					// TODO: Add this back
					// app.AddCommaToSearchIndex(app.site.SearchJSONFilePath)
					commaNeeded = false
				}
				if err = app.publishFile(filepath.Join(dir, file.Name())); err != nil {
					return ErrCode("PREVIOUS", err.Error())
				}
				commaNeeded = true
			}
		}

		// Close search index JSON file with ']'
		// TODO: Add this back
		// DelimitIndexJSON(a.Site.SearchJSONFilePath, false)

	}
	fmt.Printf("%v ", app.site.fileCount)
	if app.site.fileCount != 1 {
		fmt.Println("files")
	} else {
		fmt.Println("file")
	}

	if app.Flags.Info {
		app.ShowInfo()
	}
	// Return with success code.
	return nil
}

// TODO: Move this to mark package or eliinate mark
func parserWithOptions() goldmark.Markdown {
	exts := []goldmark.Extender{
		//mdext.New(mdext.WithTable()), extension.Table,

		// YAML support
		mdext.New(),

		// Support GitHub tables
		extension.Table,
		extension.GFM,
		extension.DefinitionList,
		extension.Footnote,
		// TC: Add highlighting options
		/*
			highlighting.NewHighlighting(
				highlighting.WithStyle(a.Site.MarkdownOptions.HighlightStyle),
				highlighting.WithFormatOptions()),
		*/

	}

	parserOpts := []parser.Option{parser.WithAttribute(), parser.WithAutoHeadingID()}

	renderOpts := []renderer.Option{
		// WithUnsafe is required for HTML templates to work properly
		html.WithUnsafe(),
		html.WithXHTML(),
	}
	// TC: Add as option?
	/*
		if a.Site.MarkdownOptions.hardWraps {
			renderOpts = append(renderOpts, html.WithHardWraps())
		}
	*/

	return goldmark.New(
		goldmark.WithExtensions(exts...),
		goldmark.WithParserOptions(parserOpts...),
		goldmark.WithRendererOptions(renderOpts...),
	)
}
