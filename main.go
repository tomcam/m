package main

import (
  "github.com/tomcam/pkg/mdext/tom"
	"bytes"
	"errors"
	"fmt"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"io"
)


func main() {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.New(meta.WithTable()),
			extension.Table,
		),
	)
	source := `---
Title: goldmark-meta
Summary: Add YAML metadata to the document
Tags:
    - markdown
    - goldmark
---

# Hello goldmark-meta
`

	var buf bytes.Buffer
	if err := markdown.Convert([]byte(source), &buf); err != nil {
		panic(err)
	}
	fmt.Print(buf.String())
}

// run() is used for testing instead of main(). See:
// https://pace.dev/blog/2020/02/12/why-you-shouldnt-use-func-main-in-golang-by-mat-ryer.html
func run(args []string, stdout io.Writer) error {
	if len(args) < 2 {
		return errors.New("no names")
	}
	for _, name := range args[1:] {
		fmt.Fprintf(stdout, "Hi %s", name)
	}
	return nil
}





