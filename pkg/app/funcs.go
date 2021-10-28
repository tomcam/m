package app

import (
	//"bytes"
	//"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	//"strconv"
	"strings"
	"time"
)

// articlefunc() returns the contents of the Markdown file itself.
// It can only be used from one of the page regions, not inside
// the markdown text, because that would cause a Markdown inception.
func (app *App) articlefunc(params ...string) string {
	if len(params) < 1 {
		return string(app.Site.webPages[app.Page.filePath].html)
	} else {
		return string(app.Site.webPages[app.Page.filePath].html)
	}
}

// dirNames() returns a directory listing of the specified
// file names in the document's directory
func (app *App) dirNames(params ...string) []string {
	files, err := ioutil.ReadDir(app.Page.dir)
	if err != nil {
		return []string{}
	}
	var ret []string
	for _, file := range files {
		ret = append(ret, file.Name())
	}
	return ret
}

// files() obtains a slice of filenames in the specified
// directory, using a wildcard specified in suffix.
// Example: {{ files "." "*.jpg" }}
// TODO:
// 1. Consider security issues
// 2. Handle insufficent input
func (a *App) files(dir, suffix string) []string {
	files, err := filepath.Glob(filepath.Join(dir, suffix))
	if err != nil {
		return []string{}
	} else {
		return files
	}
}

// ftime() returns the current, local, formatted time.
// Can pass in a formatting string
// https://golang.org/pkg/time/#Time.Format
func (a *App) ftime(param ...string) string {
	var ref = "Mon Jan 2 15:04:05 -0700 MST 2006"
	var format string

	if len(param) < 1 {
		format = ref
	} else {
		format = param[0]
	}
	t := time.Now()
	return t.Format(format)
}

// hostname() returns the name of the machine
// this code is running on
func (a *App) hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	} else {
		return hostname
	}
}

// inc inserts the named file into the current Markdown file.
/* Treats it as a Go template, so either HTML or Markdown
work fine.
The location of the file appears first, before a pipe character.
It can be one of:

"article" for the current markdown file's directory
"common" for the Site.commonSubDir subdirectory

So it might look like :

  inc "articles|kitchen.md"

*/
func (app *App) inc(filename string) template.HTML {

	// Read the HTML file into a byte slice.
	var input []byte
	var err error
	if filename == "" {
		return template.HTML("")
	}
	parsed := strings.Split(filename, "|")
	// If nothing specified, look in article directory
	if len(parsed) < 2 {
		filename = filepath.Join(app.Page.dir, parsed[0])
	} else {
		location := parsed[0]
		filename = parsed[1]

		switch strings.ToLower(location) {
		case "article":
			filename = filepath.Join(app.Page.dir, filename)
		case "common":
			filename = filepath.Join(app.Site.commonPath, filename)
		default:
			app.QuitError(ErrCode("0119", location))
		}
	}
	if !fileExists(filename) {
		app.QuitError(ErrCode("0120", filename))
	}

	input, err = ioutil.ReadFile(filename)
	if err != nil {
		app.QuitError(ErrCode("0121", filename))
	}

	// Apply the template to it.
	// The one function missing from fewerFuncs is shortcode() itself.
	s := app.execute(filename, string(input), app.fewerFuncs)
	return template.HTML(s)
}

// path() returns the current markdown document's directory
func (app *App) path() string {
	return currDir()
}

// scode() provides a not-very-good shortcode feature. Can't figure
// out how to do a better job considering a Go template function
// can take only a map, but you can't pass map literals.
// You need to pass it a map with a key named "filename"
// that matches a file in ".scodes". Currently the
// only thing that works with Javascript is youtube.html
/*
   Example:
   ====
   [List]
   youtube = { filename="youtube.html", id = "dQw4w9WgXcQ" }
   ===
   ## Youtube?
   {{ scode .FrontMatter.List.youtube }}


*/
func (app *App) scode(params map[string]interface{}) template.HTML {
	filename, ok := params["filename"].(string)
	if !ok {
		return template.HTML("filename missing")
	}
	var input []byte
	var err error
	if len(params) < 1 {
		return ("ERROR0")
	}

	// If no extension specified assume HTML
	if filepath.Ext(filename) == "" {
		filename = replaceExtension(filename, "html")
	}

	// Find that file in the shortcode file directory
	filename = filepath.Join(app.Site.sCodePath, filename)

	if !fileExists(filename) {
		app.QuitError(ErrCode("0122", filename))
	}

	input, err = ioutil.ReadFile(filename)
	if err != nil {
		app.QuitError(ErrCode("0123", filename))
	}

	// Apply the template to it.
	// The one function missing from fewerFuncs is shortcode() itself.
	s := app.execute(filename, string(input), app.fewerFuncs)
	return template.HTML(s)
}

func (a *App) addTemplateFunctions() {
	a.funcs = template.FuncMap{
		"article":  a.articlefunc,
		"dirnames": a.dirNames,
		"files":    a.files,
		"ftime":    a.ftime,
		"hostname": a.hostname,
		"inc":      a.inc,
		"path":     a.path,
		"scode":    a.scode,
		//"toc":      a.toc,
	}
}
