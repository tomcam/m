package mdext

import (
	"github.com/google/go-cmp/cmp"
	"github.com/tomcam/mb/pkg/texts"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"testing"
)

func TestBuildDocBody(t *testing.T) {
	tests := []struct {
		mdSrc string
		want  string
	}{
		{texts.Dedent(`
		   foo bar
		   qux baz 
		`), "foo bar qux baz"},
		{texts.Dedent(`
       # h1
		   foo *bar* _baz_.
		   qux qux 
		`), "h1\nfoo bar  baz . qux qux"},
	}
	for _, tt := range tests {
		t.Run(tt.mdSrc, func(t *testing.T) {
			gm := goldmark.New(goldmark.WithParserOptions(parser.WithAutoHeadingID()))
			root := gm.Parser().Parse(text.NewReader([]byte(tt.mdSrc)))
			got := BuildDocBody(root, []byte(tt.mdSrc))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf(" mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
