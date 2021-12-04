package app

import (
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
	Type string `yaml:"Type"` // Page, Posts, Gallery, Category
	//Name string `yaml:"Name"`

	// Derived from the map key name if not given here
	Filename string `yaml:"Filename"`

	// Directory this should appear in
	Folder string `yaml:"Folder"`

	// Sort order if gallery, category, or posts
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

// via newSite() and that we're in the
// project site specified in app.Site.path
func (app *App) generate(pathname string) error {
	var starters map[string]Starter
	//pathname = filepath.Join(app.Site.path, "starter.yaml")
	b, err := ioutil.ReadFile(pathname)
	if err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	err = yaml.Unmarshal(b, &starters)
	if err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
	/* why? Just why?
		if err = readStarterConfig(pathname, &starter); err != nil {
			app.QuitError(ErrCode("0115", pathname))
			return ErrCode("PREVIOUS", err.Error())
		}
		if err := readYAMLFile(pathname, starters); err != nil {
	 		app.QuitError(ErrCode("0115", pathname))
			return ErrCode("PREVIOUS", err.Error())
		}
	*/
	for k, v := range starters {
		switch strings.ToLower(v.Type) {
		case "page":
			if err = app.starterPage(k, v); err != nil {
				return ErrCode("PREVIOUS", err.Error())
			}
		}
	}
	return nil
}

// starterPage() creates a stub page from a description
// in a YAML file with startup pages
func (app *App) starterPage(name string, starter Starter) error {
	app.Debug("starterPage(%v). Folder: %v", name, starter.Folder)
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
