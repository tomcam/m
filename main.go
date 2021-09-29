package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tomcam/m/pkg/mdext"
	"github.com/tomcam/m/pkg/util"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"io"
	"os"
)

func main() {
	filename := os.Args[1]
	fmt.Printf("Filename: %#v\n", filename)
	markdown := goldmark.New(
		goldmark.WithExtensions(
			mdext.New(mdext.WithTable()),
			extension.Table,
		),
	)
	source := `+++
Title: Front matter
Summary: Add YAML metadata to the document
Tags:
    - markdown
    - goldmark
+++

`
	var buf bytes.Buffer
	s := util.FileToBytes(filename)
	if err := markdown.Convert(s, &buf); err != nil {
		panic(err)
	}
	fmt.Print(buf.String())
	os.Exit(0)

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
