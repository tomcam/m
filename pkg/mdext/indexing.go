package mdext

import (
	"bytes"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
)

// Doc is the markdown text content suitable for passing to an index system.
type Doc struct {
	// The URL path, excluding the domain, to the page represented by this doc,
	// i.e. `/foo/bar/baz` and  not `http://example.com/foo/bar/baz`.
	Path string `json:"path"`
	// The text content of the title of a markdown source file.
	Title string `json:"title"`
	// The text content of a markdown source file with all formatting directives
	// removed, e.g. `foo *bar*` is represented as `foo bar`.
	Body string `json:"body"`
}

// BuildDoc creates a new Doc representing the current page that app
// is processing.
func BuildDocBody(root ast.Node, mdSrc []byte) string {
	b := new(bytes.Buffer)
	b.Grow(1024)
	appendNodeText(b, root, mdSrc)
	return b.String()
}

// appendNodeText recursively adds the text content of node and all descendants
// to the buffer. We cant use ast.Node.Text() because it doesn't add space
// between block elements.
func appendNodeText(b *bytes.Buffer, node ast.Node, mdSrc []byte) {
	hasContent := b.Len() > 0
	if node.Type() == ast.TypeBlock && hasContent {
		b.WriteByte('\n') // space out block nodes
	} else if hasContent {
		lastByte := b.Bytes()[b.Len()-1]
		if lastByte != ' ' && lastByte != '\n' {
			b.WriteByte(' ') // space out anything else
		}
	}

	switch n := node.(type) {
	case *ast.String:
		b.Write(n.Value)
		return
	case *ast.Text:
		b.Write(n.Segment.Value(mdSrc))
		return
	case *parser.Delimiter:
		b.Write(n.Segment.Value(mdSrc))
		return
	}

	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		appendNodeText(b, c, mdSrc)
	}
}
