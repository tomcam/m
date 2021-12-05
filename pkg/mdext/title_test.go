package mdext

import (
	"github.com/google/go-cmp/cmp"
	"github.com/tomcam/mb/pkg/texts"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"testing"
)

func TestInferTitle(t *testing.T) {
	tests := []struct {
		mdSrc string
		want  string
	}{
		{texts.Dedent(`
		     # h1.1
		     body
		     ## h2.1
		`), "h1.1"},
		{texts.Dedent(`
		     # h1.1 *foo* bar
		     body
		     ## h2.1
		`), "h1.1 foo bar"},
		{texts.Dedent(`
		     ## h2.1
		     # h1.1 *foo* bar
		`), "h1.1 foo bar"},
		{texts.Dedent(`
		     ###### h6.1
		`), "h6.1"},
	}
	for _, tt := range tests {
		t.Run(tt.mdSrc, func(t *testing.T) {
			gm := goldmark.New(goldmark.WithParserOptions(parser.WithAutoHeadingID()))
			root := gm.Parser().Parse(text.NewReader([]byte(tt.mdSrc)))
			title := InferTitle(root, []byte(tt.mdSrc))
			if diff := cmp.Diff(tt.want, title); diff != "" {
				t.Errorf("InferTitle() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
