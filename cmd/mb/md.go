package main

import (
	"bytes"
	//"errors"
	"fmt"
	//highlighting "github.com/yuin/goldmark-highlighting"
	//"github.com/tomcam/m/pkg/app"
	"github.com/tomcam/m/pkg/mdext"
	"github.com/tomcam/m/pkg/util"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	//"io"
	//"os"
)

func mdToHTML(source []byte) []byte {
	// Get a parser object.
	ctx := parser.NewContext()
	p := newParser().Parser()
	node := p.Parse(text.NewReader(source), parser.WithContext(ctx))
	// Create variable-sized buffer for parsed output.
	buf := new(bytes.Buffer)
	// Convert the Markdown file to HTML.
	if err := newParser().Renderer().Render(buf, source, node); err != nil {
		// TC: update error handling
		// a.QuitError(errs.ErrCode("0920", err.Error()))
		panic("Parse error")
		// TODO: Hmmm... this should probably return an error
		return nil
	}
	// Return the HTML.
	fmt.Println("mdToHTML: \n" + string(buf.Bytes()))
	return buf.Bytes()
}

// mdFileToHTML converts the markdown file in filename to HTML.
// It may include optional front matter.
func mdFileToHTML(filename string) []byte {
	// Read file into a byte slice.
	s := util.FileToBytes(filename)
	return mdToHTML(s)

	// Get a parser object.
	ctx := parser.NewContext()
	p := newParser().Parser()
	node := p.Parse(text.NewReader(s), parser.WithContext(ctx))
	// Create variable-sized buffer for parsed output.
	buf := new(bytes.Buffer)
	// Convert the Markdown file to HTML.
	if err := newParser().Renderer().Render(buf, s, node); err != nil {
		// TC: update error handling
		// a.QuitError(errs.ErrCode("0920", err.Error()))
		panic("Parse error")
		// TODO: Hmmm... this should probably return an error
		return nil
	}
	// Return the HTML.
	return buf.Bytes()
}

func newParser() goldmark.Markdown {
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


