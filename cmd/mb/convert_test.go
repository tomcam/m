package main

import (
	"github.com/google/go-cmp/cmp"
	"github.com/tomcam/m/pkg/texts"
	"testing"
)

func TestConversion(t *testing.T) {
	tests := []struct {
		mdSrc string
		want  string
	}{
		{`
		   +++
       Stuff: here
       +++
       # 1
		`, "<h1 id=\"1\">1</h1>\n"},
		{`
       hello
		`, "<p>hello</p>\n"},
		{`
       hello


		`, "<p>hello</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.mdSrc, func(t *testing.T) {
      //got := mdToHTML(tt.mdSrc)
      got := string(mdToHTML([]byte(texts.Dedent(tt.mdSrc))))
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf(" mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

