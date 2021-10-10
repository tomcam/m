package main

import (
	"github.com/google/go-cmp/cmp"
	"github.com/tomcam/m/pkg/texts"
	//"github.com/yuin/goldmark/text"
	"testing"
)

var test1 = `# 1`

// Very beginnings of end-to-end testing
func iTestConversion(t *testing.T) {
	got := string(mdToHTML([]byte(test1)))
	want := "<h1 id=\"1\">1</h1>\n"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}



func TestConversion(t *testing.T) {
	tests := []struct {
		mdSrc string
		want  string
	}{
		{texts.Dedent(`
		   +++
       Stuff: here
       +++
       # 1
		`), "<h1 id=\"1\">1</h1>\n"},
		{texts.Dedent(`
       hello
		`), "<p>hello</p>\n"},
		{texts.Dedent(`
       hello


		`), "<p>hello</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.mdSrc, func(t *testing.T) {
      //got := mdToHTML(tt.mdSrc)
      got := string(mdToHTML([]byte(tt.mdSrc)))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf(" mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

