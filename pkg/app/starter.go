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

// Description makes up what you need for a Description metatag.
// The Tag field is the most important, but if you want to
// append something like "| blog" you'd use After for that.
// Likewise for Before, but it creates a suffix.
type Description struct {
	Before string `yaml:"Before"`
	Tag    string `yaml:"Description"`
	After  string `yaml:"After"`
}

// Generate a list of pages, posts, galleries, or categories
// to avoid copy pasta. Normally used when
// the site is created (or later, but it makes most
// sense upon site creation).
type Starter struct {
	Type string `yaml:"Type"` // e.g. Page, Collection
	//Name string `yaml:"Name"`

	// Derived from the map key name if not given here
	Filename string `yaml:"Filename"`

	// Directory this should appear in
	Folder string `yaml:"Folder"`

	// Sort order if collection
	Sort string `yaml:"Sort"`

	// If specified, the Permalink template makes it post
	// and not a page
	Permalink string `yaml:"Permalink"`

	// For title tag
	Title string `yaml:"Title"`

	// For description meta tag
	Description Description `yaml:"Description"`

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

// starterPages() creates a collection directory from a description
// in a YAML file with startup pages
func (app *App) starterCollection(name string, starter Starter) error {
	app.Note("\n\nstarterCollection(%v) Folder: %v", name, starter.Folder)
	dir := starter.Folder
	// If no folder is given, assume project root.
	// Remember Go uses Unix folder conventions even
	// under Windows
	if dir == "" {
		dir = "/"
	}

	if app.Site.Collections[dir].path == dir {
		return ErrCode("0953", dir)
	}

	// Create the specified folder as a subdirectory
	// of the current project.
	dir = filepath.Join(app.Site.path, dir)
	err := os.MkdirAll(dir, defaults.ProjectFilePermissions)
	if err != nil {
		return ErrCode("0415", dir)
	}

	var c Collection
	// TODO: Need test case
	if permalink, err := validatePermalink(starter.Permalink); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	} else {
		c.permalink = permalink
	}
	c.path = dir
	// xxx

	app.Site.Collections[dir] = c
  fmt.Println(app.Site.Collections)
	return nil
}

// validatePermalink() ensures that
// the proposed permalink can be used in a directory
// structure reliably.
func validatePermalink(permalink string) (string, error) {
	defaultPermalink := ":year/:monthnum/:day/:postname"
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
		case ":year", ":monthnum", ":day", ":hour",
			":minute", ":second", ":postname", ":author":
			/* Do nothing */
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
	//app.Note("starterPage(%v) Folder: %v", name, starter.Folder)
	dir := starter.Folder
	if name == "" {
		return ErrCode("1104", dir)
	}
	// If no folder is given, assume project root.
	// Remember Go uses Unix folder conventions even
	// under Windows
	if dir == "" {
		dir = "/"
	}
	// Create the specified folder as a subdirectory
	// of the current project.
	dir = filepath.Join(app.Site.path, dir)
	err := os.MkdirAll(dir, defaults.ProjectFilePermissions)
	if err != nil {
		return ErrCode("0410", dir)
	}
	app.Debug("\tDir: %v", dir)

	var filename string
	// Convert the name to a filename.
	if starter.Filename == "" {
		filename = slug.Make(name)
	} else {
		filename = starter.Filename
	}

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
	if starter.Description.Tag != "" {
		title = "Description: " + starter.Description.Tag + "\n"
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
