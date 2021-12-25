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

	// Sort order if collection
	Sort string `yaml:"Sort"`

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
	app.Debug("generate(%s)", pathname)
	var starters map[string]Starter
	b, err := ioutil.ReadFile(pathname)
	if err != nil {
		// TODO: Improve error handling
		return ErrCode("PREVIOUS", err.Error())
	}
	err = yaml.Unmarshal(b, &starters)
	if err != nil {
		msg := fmt.Sprintf("%s: %s", pathname, err.Error())
		return ErrCode("0135", msg)
	}
	app.Site.starterFile = pathname
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
			if err = app.newCollection(k, v, pathname); err != nil {
				return ErrCode("PREVIOUS", err.Error())
			}
		}
	}
	return nil
}

// newCollection() creates a collection from the path and/or
// permalink described in name, for example,
// "/blog/:year/:monthnum/:daynum"
// And adds it to Site.Collections
// filename is the name of the starter file. If interactive,
// just pass ""
func (app *App) newCollection(name string, starter Starter, filename string) error {
	app.Debug("newCollection(%v)", name)
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

	// The key to the collection will be the base directory
	// (everything up to the first permalink variable)
	base := permalinkBase(permalink)

	// Quit if this is a duplicate
	if app.Site.Collections[base].Permalink == permalink {
		msg := ""
		if filename != "" {
			msg = fmt.Sprintf("Starter file %s already has a collection named %s", filename, base)
		} else {
			msg = fmt.Sprintf("Collection named %s already exists", base)
		}
		return ErrCode("0954", msg)
	}
	c.Permalink = permalink
	app.Debug("\tAbout to add %#v to %#v at %v", c, app.Site.Collections, base)
	app.Site.Collections[base] = c
	app.Debug("\tapp.Site.Collections[%v] is now %v", base, app.Site.Collections[base])
	//app.Print("Permalink: %v\n app.Site.Collections.[permalink]: %v\n. FirstDir: %v\n", permalink, app.Site.Collections[permalink], permalinkBase(permalink))

	// Create the specified folder as a subdirectory
	// of the current project.
	// The base directory is everything up to the first
	// colon. Permalink is guaranteed to start
	// with a directory separator.
	// TODO: refactor with permalinkBase()?
	dir := permalink[1:strings.IndexRune(permalink, ':')]
	err = os.MkdirAll(dir, defaults.ProjectFilePermissions)
	if err != nil {
		return ErrCode("0415", dir)
	}
	//xxx

	if err = app.writeSiteConfig(); err != nil {
		return ErrCode("1301", name)
	}
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

// permalinkBase takes a permalink such as "/site/news/:year/:month/:postname"
// and yields its bath bath, for example,  "/site/news/"
func permalinkBase(permalink string) string {
	return permalink[:strings.IndexRune(permalink, ':')]
}

// fixPermalink() ensures that
// the proposed permalink can be used in a
// collection-style directory
// structure reliably. Among other things it ensures
// that the first part of path is fixed, that it begins
// with the path separator, and that :postname
// appears last. It will rewrite the permalink
// if these conventions aren't followed.
// So, "blog" would be transformed into "/blog/:permalink",
// ":permalink/blog" would be transformed into
// "/blog/:permalink", the empty string becomes
// ":year/:monthnum/:daynum/:postname",
// "news/:year/:monthnum/:daynum/" is transformed
// into "/news/:year/:monthnum/:daynum/:postname", etc.
// TODO: Use above comment to generate test cases, and
// include test cases where the input is expected to be
// the same as the output.
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
		case ":year", ":month", ":monthnum", ":daynum", ":day", ":hour",
			":minute", ":second", ":postname", ":author":
			if i == 0 {
				return "", ErrCode("1208", "")
			}
		default:
			segments[i] = slug.Make(seg)
		}
		// If the permalink didn't include post name
		// then append it.
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

	// Collapse everything back to what should now be
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
