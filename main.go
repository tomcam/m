package main

import (

	//"github.com/tomcam/m/goldmark-y"
	"bytes"
	"errors"
	"fmt"
	"github.com/yuin/goldmark"
	//"github.com/yuin/goldmark-meta"
	"github.com/tomcam/m/pkg/mdext"
	"github.com/yuin/goldmark/extension"
	"io"
  "io/ioutil"
  "os"
)

// fileToBytes
func fileToBytes(filename string) []byte {
  bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}
	}
  return bytes
}

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
  s := fileToBytes(filename)
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
