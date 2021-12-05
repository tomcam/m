package mdext

import (
	"github.com/google/go-cmp/cmp"
	"github.com/tomcam/mb/pkg/texts"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"testing"
)

func TestExtractTOCs(t *testing.T) {
	tests := []struct {
		mdSrc string
		level int
		want  []TOCEntry
	}{
		{texts.Dedent(`
		     # h1.1
		     body
		     ## h2.1
		`), 2, []TOCEntry{{"h11", "h1.1", 1}, {"h21", "h2.1", 2}}},
		{texts.Dedent(`
		     # h1.1 *foo* bar
		     body
		     ## h2.1
		     ## h2.2
		     ### h3.1
		     # h1.2
		`), 2, []TOCEntry{
			{"h11-foo-bar", "h1.1 <em>foo</em> bar", 1},
			{"h21", "h2.1", 2},
			{"h22", "h2.2", 2},
			{"h12", "h1.2", 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.mdSrc, func(t *testing.T) {
			gm := goldmark.New(goldmark.WithParserOptions(parser.WithAutoHeadingID()))
			root := gm.Parser().Parse(text.NewReader([]byte(tt.mdSrc)))
			tocs, err := ExtractTOCs(gm.Renderer(), root, []byte(tt.mdSrc), tt.level)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, tocs); diff != "" {
				t.Errorf("ExtractTOCs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
