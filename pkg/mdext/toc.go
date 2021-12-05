package mdext

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
)

// TOCEntry is an individual entry in the table of contents.
type TOCEntry struct {
	// The HTML ID of the source header or the empty string if the header had no ID.
	ID string
	// The text of the header
	Header string
	// Its level 1-6 for h1-h6
	Level int
}

// ExtractTOCs returns TOC entries corresponding to each header in the markdown
// document ordered by appearance.
func ExtractTOCs(rn renderer.Renderer, node ast.Node, mdSrc []byte, level int) ([]TOCEntry, error) {
	const numTocs = 8 // assume a reasonable amount of headers per md document
	tocNodes := make([]*ast.Heading, 0, numTocs)
	// Safe to ignore error because only our walk function can error and we don't
	// error.
	_ = ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || n.Kind() != ast.KindHeading {
			return ast.WalkContinue, nil
		}
		h := n.(*ast.Heading)
		if h.Level <= level {
			tocNodes = append(tocNodes, h)
		}
		return ast.WalkContinue, nil
	})

	// Render each TOC into an HTML string in a TOCEntry.
	tocs := make([]TOCEntry, len(tocNodes))
	buf := new(bytes.Buffer)
	const headerSize = 64
	buf.Grow(headerSize)
	for i, node := range tocNodes {
		// Render the children of the heading because we don't want the enclosing
		// <h1> tags.
		for c := node.FirstChild(); c != nil; c = c.NextSibling() {
			if err := rn.Render(buf, mdSrc, c); err != nil {
				return nil, fmt.Errorf("extract TOC render: %w", err)
			}
		}
		id := ""
		if attr, ok := node.AttributeString("id"); ok {
			id = string(attr.([]byte))
		}
		tocs[i] = TOCEntry{
			ID:     id,
			Header: buf.String(),
			Level:  node.Level,
		}
		buf.Reset()
	}
	return tocs, nil
}
