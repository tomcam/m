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
			if err = app.createStarterPage(k, v); err != nil {
				return ErrCode("PREVIOUS", err.Error())
			}
		case "collection":
			if err = app.newCollection(k /* v, */, pathname); err != nil {
				return ErrCode("PREVIOUS", err.Error())
			}
		}
	}
	return nil
}

// Retrieve the first directory in the URL-like
// string passed in
// https://stackoverflow.com/a/70342730/478311
// Thanks for saving my remaining sanity on this, stt106!
// Turned out not to be useful because I forgot about
// the fact that one may have something like
// dir1/dir2/:year/:monthum/:daynum
func firstDir(permalink string) string {
	split := strings.Split(permalink, string(os.PathSeparator))
	return split[1]
}

// permalinkBase takes a permalink such as "/site/news/:year/:month/:postname"
// and yields its bath bath, for example,  "/site/news/"
func permalinkBase(permalink string) string {
	if strings.Contains(permalink, ":") {
		return permalink[:strings.IndexRune(permalink, ':')]
	} else {
		return permalink
	}
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
// filename is the name of the starter file, if any
func fixPermalink(permalink, filename string) (string, error) {
	defaultPermalink := ":year/:monthnum/:daynum/:postname"
	if permalinkBase(permalink) == permalink {
		permalink = filepath.Join(permalink, defaultPermalink)
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
				msg := fmt.Sprintf( /* Unknown permalink variable*/ " %s (filename %s)", seg, filename)
				return "", ErrCode("1208", msg)
			}
		default:
			if strings.HasPrefix(seg, ":") {
				// Unrecognized permalink variable
				msg := fmt.Sprintf("%s has unknown permalink variable %s", filename, seg)
				return "", ErrCode("1209", msg)

			}
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

// createStarterPage() generates a stub page from a description
// in a YAML file with startup pages
// Creates it directory if need be.
func (app *App) createStarterPage(name string, starter Starter) error {
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
		return ErrCode("0417", dir)
	}

	filename = slug.Make(filename)

	// Get the fully qualified filename to generate
	filename = filepath.Join(dir, filename)

	var frontMatter FrontMatter
	frontMatter.Theme = starter.Theme
	frontMatter.Title = starter.Title
	frontMatter.Description = starter.Description
	frontMatter.Sidebar = starter.Sidebar
	// See if the filename has a Markdown extension
	if !isMarkdownFile(filename) {
		filename = filename + ".md"
	}
	return app.createPageFrontMatter(filename, starter.Article, frontMatter)

}
