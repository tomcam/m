package main

import (
	"bytes"
	"errors"
	"fmt"
	//"github.com/tomcam/m/pkg/mdext"
  "github.com/tomcam/m/pkg/tomltc"
	"github.com/tomcam/m/pkg/util"
	"github.com/yuin/goldmark"
  "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
  //highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
  "github.com/yuin/goldmark/renderer/html"
  "github.com/yuin/goldmark/text"
	"io"
	"os"
)

// mdFileToHTML converts the markdown file in filename to HTML.
// It may include optional front matter.
func mdFileToHTML(filename string) []byte{
	s := util.FileToBytes(filename)
  ctx := parser.NewContext()
	p := newGoldmark().Parser()
  node := p.Parse(text.NewReader(s), parser.WithContext(ctx))
  buf := new(bytes.Buffer)
  if err := newGoldmark().Renderer().Render(buf, s, node); err != nil {
    // TC: update error handling
    // a.QuitError(errs.ErrCode("0920", err.Error()))
    panic("Parse error")
		return nil
	}
	return buf.Bytes()
}

func newGoldmark() goldmark.Markdown {
	exts := []goldmark.Extender{
    mdext.New(mdext.WithTable()),extension.Table,
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

	parserOpts := []parser.Option{parser.WithAttribute(),parser.WithAutoHeadingID()}

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




func main() {
	filename := os.Args[1]
	fmt.Printf("Filename: %#v\n", filename)
  /*
	markdown := goldmark.New(
		goldmark.WithExtensions(
			mdext.New(mdext.WithTable()),
			extension.Table,
		),
	)
  */
/*
source := `+++
Title: Front matter
Summary: Add YAML metadata to the document
Tags:
    - markdown
    - goldmark
+++

`
*/

  fmt.Println(string(mdFileToHTML(filename)))
}



// run() is used for testing instead of main(). See:
// https://pace.dev/blog/2020/02/12/why-you-shouldnt-use-func-main-in-golang-by-mat-ryer.html
func run(args []string, stdout io.Writer) error {
	if len(args) < 2 {
		return errors.New("no names")
	}
	for _, name := range args[1:] {
		fmt.Fprintf(stdout, "Hi %s", name)
	}
	return nil
}
