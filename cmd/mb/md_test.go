package main

import (
	"github.com/google/go-cmp/cmp"
	"github.com/tomcam/m/pkg/texts"
	"testing"
)

func TestConversion(t *testing.T) {
	tests := []struct {
    // Line of Markdown source to test
		mdSrc string

    // Expected HTML output of mdSrc after conversion
    // to Markdown
		want  string
	}{
    // Generate an h1.
    // Ensure front matter isn't included in output 
    // Ensure an ID is added to the header.
		{`
		   +++
       Stuff: here
       +++
       # 1
		`, "<h1 id=\"1\">1</h1>\n"},

    // Ensure whitespace in input is ignored properly
    // (compare output with next test)
		{`
       hello, world.
		`, "<p>hello, world.</p>\n"},

    // Ensure whitespace in input is consistent
    // with previous test even though leading
    // and trailing whitespace differ between
    // the two test cases.
		{`
       hello, world.


		`, "<p>hello, world.</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.mdSrc, func(t *testing.T) {
      // Go through each test case.
      // Dedent() removes semantically unnecessary
      // whitespace from output so it can be reliably
      // compared to the expected output.

      // Convert an example fragment of Markdown source to HTML.
      got := string(mdToHTML([]byte(texts.Dedent(tt.mdSrc))))
      // Compare input to expected output
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf(" mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

