package main

import (
  "testing"
)

var test1 = `# 1`

// Very beginnings of end-to-end testing
func TestConversion(t *testing.T) {
  got := string(mdToHTML([]byte(test1)))
  want := "<h1 id=\"1\">1</h1>\n"
  if got != want {
          t.Errorf("got %q, wanted %q", got, want)
  }
}

