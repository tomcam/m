package mark

import (
	//"bytes"
	//"errors"
	//"fmt"
	//highlighting "github.com/yuin/goldmark-highlighting"
	//"github.com/tomcam/m/pkg/app"
	"github.com/tomcam/m/pkg/mdext"
	//"github.com/tomcam/m/pkg/util"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	//"github.com/yuin/goldmark/text"
	//"io"
	//"os"
)

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
