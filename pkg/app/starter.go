package app

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/tomcam/m/pkg/default"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Generate a list of pages, collections
// from a config file. Normally used when
// the site is created (or later, but it makes most
// sense upon site creation).
type Starter struct {
	Type string `yaml:"Type"` // e.g. Page, Collection
	//Name string `yaml:"Name"`

	// Derived from the map key name if not given here
	//Filename string `yaml:"Filename"`

	// Directory this should appear in
	//Folder string `yaml:"Folder"`

	// Sort order if collection
	Sort string `yaml:"Sort"`

	// If specified, the Permalink template makes it post
	// and not a page
	//Permalink string `yaml:"Permalink"`

	// For title tag
	Title string `yaml:"Title"`

	// For description meta tag
	Description string `yaml:"Description"`

	// Name of theme to use for this page/section
	Theme string `yaml:"Theme"`

	// Specify sidebar direction, if any
	Sidebar string `yaml:"Sidebar"`

	// Text of the article portion, if desired
	Article string `yaml:"Article"`
}

// generate() creates files specified in the
// specified starter config file.
func (app *App) generate(pathname string) error {
	var starters map[string]Starter
	b, err := ioutil.ReadFile(pathname)
	if err != nil {
		// TODO: Improve error handling
		return ErrCode("PREVIOUS", err.Error())
	}
	err = yaml.Unmarshal(b, &starters)
	if err != nil {
		// TODO: Improve error handling
		return ErrCode("PREVIOUS", err.Error())
	}
	for k, v := range starters {
		switch strings.ToLower(v.Type) {
		default:
			return ErrCode("1207", v.Type)
		case "page":
			// TODO: Improve error handling
			if err = app.starterPage(k, v); err != nil {
				return ErrCode("PREVIOUS", err.Error())
			}
		case "collection":
			if err = app.starterCollection(k, v); err != nil {
				return ErrCode("PREVIOUS", err.Error())
			}
		}

	}
	return nil
}

func (app *App) starterCollection(name string, starter Starter) error {
	app.Note("\n\nstarterCollection(%v)", name)
	// The name is a path to the file or collection.
	// It may also be a permalink.
	path := name
	// If no folder is given, assume project root.
	// Remember Go uses Unix folder conventions even
	// under Windows
	if path == "" {
		path = "/"
	}

	// Ensure it has a leading "/" unless preceded with a dot
	if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, ".") {
		path = filepath.Join("/", path)
	}

	var c Collection
	var permalink string
	var err error
	// TODO: Need test case
	if permalink, err = fixPermalink(path); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	// TODO: This should come later after it's been cleaned
	if app.Site.Collections[permalink].path == permalink {
		return ErrCode("0953", path)
	}
	c.path = permalink
	app.Site.Collections[permalink] = c

	// Create the specified folder as a subdirectory
	// of the current project.
  // The base directory is everything up to the first
  // colon. Since permalink is guaranteed to start
  // with a directory separater
  dir := permalink[1:strings.IndexRune(permalink, ':')]

	err = os.MkdirAll(dir, defaults.ProjectFilePermissions)
	if err != nil {
		return ErrCode("0415", dir)
	}
	//xxx

	fmt.Println(app.Site.Collections)
	return nil
}

// Retrieve the first directory in the URL-like
// string passed in
// https://stackoverflow.com/a/70342730/478311
// Thanks for saving my remaining sanity on this, stt106!
// Turned out not to be useful because I forgot about
// the fact that yiou may have something like
// dir1/dir2/:year/:monthum/:daynum
func firstDir(permalink string) string {
	split := strings.Split(permalink, string(os.PathSeparator))
	return split[1]
}

// fixPermalink() ensures that
// the proposed permalink can be used in a directory
// structure reliably.
func fixPermalink(permalink string) (string, error) {
	defaultPermalink := ":year/:monthnum/:daynum/:postname"
	if permalink == "" {
		permalink = defaultPermalink
	}
	// Break the description up into segments.
	segments := strings.Split(permalink, "/")
	segs := len(segments)
	postnamePresent := false
	postnameIndex := 0
	var seg string
	// Explode the string into ":" delineated
	// path segments. Remove any slashes.
	for i := range segments {
		seg = segments[i]
		last := strings.HasSuffix(seg, "/")
		if last {
			// Remove trailing slash if any
			seg = firstN(seg, len(seg)-1)
			segments[i] = seg
		}
		switch seg {
		case ":year", ":monthnum", ":daynum", ":day", ":hour",
			":minute", ":second", ":postname", ":author":
			if i == 0 {
				return "", ErrCode("1208", "")
			}
		default:
			segments[i] = slug.Make(seg)
		}
		// If the permalink didn't include post name
		// we'll have to fix append it.
		if segments[i] == ":postname" {
			postnamePresent = true
			postnameIndex = i
		}
	}

	// Detect if postname is used anywhere but the end
	// If so remove it and replace with empty string.
	// This results in 2 slash characters at that position.
	// Should work OK because Go uses Unix conventions.
	if postnameIndex <= segs-2 {
		segments[postnameIndex] = ""
		postnamePresent = false
	}
	// Ensure :postname comes last
	if !postnamePresent {
		segments = append(segments, ":postname")
	}

	// Collapse everything pack to what should now be
	// a usable directory designation.
	clean := strings.Join(segments, "/")
	return clean, nil
}

// starterPage() creates a stub page from a description
// in a YAML file with startup pages
func (app *App) starterPage(name string, starter Starter) error {
	dir := filepath.Dir(name)
	if name == "" {
		return ErrCode("1104", dir)
	}
	// If no folder is given, assume project root.
	filename := filepath.Base(name)
	// Remember Go uses Unix folder conventions even
	// under Windows
	if dir == "" {
		dir = "/"
	}
	// Create the specified folder as a subdirectory
	// of the current project.
	dir = slug.Make(dir)
	dir = filepath.Join(app.Site.path, dir)
	err := os.MkdirAll(dir, defaults.ProjectFilePermissions)
	if err != nil {
		return ErrCode("0410", dir)
	}
	app.Debug("\tDir: %v", dir)

	filename = slug.Make(filename)

	// Get the fully qualified filename to generate
	filename = filepath.Join(dir, filename)

	// Create the front matterl
	theme := ""
	if starter.Theme != "" {
		theme = "Theme: " + starter.Theme + "\n"
	}
	title := ""

	if starter.Title != "" {
		title = "Title: " + starter.Title + "\n"
	}

	// TODO: Stuff these things into a read FrontMatter to get it right
	description := ""
	if starter.Description != "" {
		title = "Description: " + starter.Description + "\n"
	}

	frontMatter :=
		"---\n" +
			theme +
			title +
			description +
			"---\n"

	// See if the filename has a Markdown extension
	if !isMarkdownFile(filename) {
		filename = filename + ".md"
	}
	article := frontMatter + starter.Article

	if err := writeTextFile(filename, article); err != nil {
		return ErrCode("0410", filename)
	}
	return nil
}
