package mb

import (
	"errors"
	"fmt"
	"github.com/tomcam/m/pkg/app"
	"github.com/tomcam/m/pkg/mark"
	"io"
	"os"
)

func main() {
	filename := "."
	command := "build"
	switch len(os.Args) {
	case 4: // e.g. ./mb new site foo
		command = os.Args[1] + os.Args[2]
		filename = os.Args[3]
	case 3: // e.g. ./mb new foo
		command = os.Args[1]
		filename = os.Args[2]
	case 2: // e.g. ./mb new
		command = os.Args[1]
	case 1:
	default: // program name only
		// Same as build
	}
	a := app.NewApp(filename)
	fmt.Printf("\tCommand: %s\n", command)
	fmt.Printf("Args: %v\nArg count: %v\nCommand: %s\nFilename: %s\n",
		os.Args, len(os.Args), command, filename)
	switch command {
	case "build":
		a.HTML = mark.MdFileToHTML(filename)
		fmt.Println(string(a.HTML))
	case "new", "newsite":
		a.NewSite()
		fmt.Printf("\tProject path: %s\n", a.Site.Path)
	}

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
