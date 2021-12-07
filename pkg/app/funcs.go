package app

import (
	"bytes"
	"fmt"
	"github.com/tomcam/m/pkg/mdext"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// articlefunc() returns the contents of the Markdown file itself.
// It can only be used from one of the page regions, not inside
// the markdown text, because that would cause a Markdown inception.
// TODO: I think this doesn't work
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
		// TODO: return an error instead
		app.QuitError(ErrCode("0120", filename))
	}

	input, err = ioutil.ReadFile(filename)
	if err != nil {
		// TODO: return an error instead
		app.QuitError(ErrCode("0121", filename))
	}

	// Apply the template to it.
	s := app.execute(filename, string(input), app.funcs)
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
		"toc":      a.toc,
	}
}

// generateTOC reads the Markdown source and returns a slice of TOC entries
// corresponding to each header less than or equal to level.
func (a *App) generateTOC(level int) []mdext.TOCEntry {
	node := a.markdownAST(a.src)
	tocs, err := mdext.ExtractTOCs(a.newGoldmark().Renderer(), node, a.src, level)
	if err != nil {
		// TODO: this should return an error
		a.QuitError(ErrCode("0926", err.Error()))
	}
	return tocs
}

// toc generates a table of contents and includes all headers with a level less
// than or equal level. Level must be 1-6 inclusive.
func (a *App) toc(params ...string) string {
	pcount := len(params)
	var listType string
	var level int
	var err error
	switch pcount {
	case 0:
		{
			level = 6
			listType = "ul"
		}
	case 1:
		{
			level, err = strconv.Atoi(params[0])
			listType = "ul"
		}
	default:
		{
			level, err = strconv.Atoi(params[0])
			listType = params[1]
			if strings.Contains(listType, "ol") {
				listType = "ol"
			} else {
				listType = "ul"
			}
		}
	}

	// Please leave this error code as is
	if err != nil {
		a.QuitError(ErrCode("1205", err.Error()))
	}
	// Ditto
	if level <= 0 || level > 6 {
		// TODO: Return an error
		a.QuitError(ErrCode("1206", params[0]))
	}
	tocs := a.generateTOC(level)
	b := new(bytes.Buffer)
	b.Grow(256)
	writeTOCLevel(listType, b, tocs, 1)
	return b.String()
}

// writeTOCLevel writes a single TOC level and recursively delegates for child
// levels. The result is nested HTML lists corresponding to TOC levels.
func writeTOCLevel(listType string, b *bytes.Buffer, tocs []mdext.TOCEntry, level int) int {
	openTag := "<" + listType + ">"
	closeTag := "</" + listType + ">"
	b.WriteString(openTag)
	i := 0 // explicit index because recursive calls advance i by variable amount
loop:
	for {
		if i >= len(tocs) {
			break
		}
		toc := tocs[i]
		switch {
		case toc.Level < level:
			break loop
		case toc.Level == level:
			b.WriteString("<li>")
			_, _ = fmt.Fprintf(b, `<a href="#%s">`, toc.ID)
			b.WriteString(toc.Header)
			b.WriteString("</a>")
			b.WriteString("</li>")
		case toc.Level > level:
			b.WriteString("<li>")
			// We're adding i instead of assigning because we pass a smaller slice
			// to the recursive call.
			i += writeTOCLevel(listType, b, tocs[i:], level+1)
			b.WriteString("</li>")
			continue // skip i++ since the child must have made progress
		}
		i++
	}
	b.WriteString(closeTag)
	return i
}
