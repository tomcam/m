package app

import (
	"embed"
	"fmt"
	"github.com/tomcam/m/pkg/default"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

//go:embed .mb
var mbfiles embed.FS

// Site contains configuration specific to each site, such as
// its directory location, title, publish directory,
// branding information, etc.
type Site struct {
	// Target subdirectory for assets such as CSS and images.
	// It's expected to be a child of the Publish directory.
	// The function App.assetDir() computes the full
	// path of that directory, based on the app path,
	// the current theme, etc.
	// See also its subdirectories, CSSDir and ImageDir
	assetPath string

	// Make it easy if you just have 1 author.
	Author author `yaml:"Author"`

	// List of secondary authors
	CoAuthors []author `yaml:"CoAuthors"`

	// Base directory for URL root, which may be different
	// from its actual root. For example, GitHub Pages prefers
	// the blog to start in /docs instead of root, but
	// a URL would omit it.
	BasePath string `yaml:"Base-Path"`

	// Site's branding, any string, that user specifies in site.toml.
	// So, for example, if the Name is 'my-project' this might be
	// 'My Insanely Cool Project"
	Branding string `yaml:"Branding"`

	// Full pathname of common directory. Derived from CommonSubDir
	commonPath string

	// Company name & other info user specifies in site.toml
	Company company `yaml:"Company"`

	// Subdirectory under the AssetDir where
	// CSS files go when published.
	cssPublishPath string

	// List of all directories in the site
	dirs map[string]dirInfo

	// List of directories in the source project directory that should be
	// excluded, things like ".git" and "node_modules".
	// Note that directory names starting with a "." are excluded too.
	// DO NOT ACCESS DIRECTLY:
	// Use excludedDirs() because it applies other information such as PublishDir()
	ExcludeDirs []string `yaml:"Exclude-Dirs"`

	// List of file extensions to exclude. For example. [ ".css" ".go" ".php" ]
	ExcludeExtensions []string `yaml:"Exclude-Extensions"`

	// Number of markdown files processed
	fileCount uint

	// Full pathname of site file so it can be read
	// using {{ Site.Filename }}.
	// READ ONLY
	Filename string `'yaml:"Filename"`

	// Google Analytics tracking ID specified in site.toml
	Ganalytics string `yaml:"Ganalytics"`

	// All these files are copied into the HTML header.
	// Example: favicon links.
	HeadTags []string `yaml:"Head-Tags"`

	// Full path of header tags for "code injection"
	headTagsPath string

	// Subdirectory under the AssetDir where image files go
	imagePath string

	// for HTML header, as in "en" or "fr"
	Language string `yaml:"Language"`

	// Flags indicating which non-CommonMark Markdown extensions to use
	MarkdownOptions MarkdownOptions `yaml:"Markdown-Options"`

	// Mode ("dark" or "light") used by this site unless overridden in front matter
	Mode string `yaml:"Mode"`

	// Site's project name, so it's a filename
	// (not a pathname).
	// It's an identifier so it should be in slug format:
	// Preferably just alphanumerics, underline or hyphen, and
	// no spaces, for example, 'my-project'
	name string

	// Home directory for the project. All other
	// paths, such as location of publish directory,
	// graphics & javascript assets, etc., are based on
	// this location.
	// If you run:
	//   mb new site /Users/tom/html/foo
	// It would be /Users/tom/html/foo
	// If you just run:
	//   mb new site
	// It's initialized to the name of the current directory.
	// Site.Filename is generated from this.
	path string

	// Directory for finished site--rendered HTML & asset output
	publishPath string

	// Full path of shortcode dir for this project. It's computed
	// at runtime using SCodeDir, also in this struct.
	sCodePath string

	// Full path of file containing JSON version of site text
	// to be indexed
	searchJSONFilePath string

	// Full path to directory containing scripts
	// to be copied in just before the closing HTML tag
	scriptClosePath string

	// Social media URLs
	Social social `yaml:"Social"`

	// Site defaults to using this sidebar setting unless
	// a page specifies otherwise
	Sidebar string `yaml:"Sidebar"`

	// Name (not path) of Theme used by this site unless overridden in front matter.
	Theme string `yaml:"Theme"`

	// Location of complete set of themes included
	// with product release. A subset of these
	// gets copied to the siteThemesPath directory.
	factoryThemesPath string

	// Tracks which themes are used by this site
	// to avoid copying over themes that have
	// already been copied.
	publishedThemes map[string]bool

	// Location of theme files copied over for this
	// particular site
	siteThemesPath string

	// All the rendered pages on the site, plus meta information.
	// Index by the fully qualified path name of the source .md file.
	// TODO: Not yet used
	webPages map[string]WebPage

	// List of all directories of type Collection
	Collections map[string]Collection

	// IMPORTANT
	// LIST ALWAYS GOES AT THE END OF THE FILE/DATA STRUCTURE
	// User data.
	List interface{} `yaml:"List"`
}

type company struct {
	// Company name, like "Metabuzz" or "Example Inc."
	Name       string `yaml:"Name"`
	Address    string `yaml:Address1`
	Address2   string `yaml:Address2`
	City       string `yaml:City`
	Country    string `yaml:Country`
	PostalCode string `yaml:PostalCode`
	URL        string `yaml:"URL"`

	// Logo file for the header
	HeaderLogo string `yaml:"Header-Logo"`
}
type author struct {
	FullName string `yaml:"Full-Name"`
	URL      string `yaml:"URL"`
	Role     string `yaml:"Role"`
}

type authors struct {
	Authors []author `yaml:"Authors"`
}

// MarkdownOptions consists of goldmark
// options used when converting markdown to HTML.
type MarkdownOptions struct {
	// If true, line breaks are significant
	HardWraps bool `yaml:"Hard-Wraps"`

	// Name of color scheme used for code highlighting,
	// for example, "monokai"
	// List of supported languages:
	// https://github.com/alecthomas/chroma/blob/master/README.md
	// I believe the list of themes is here:
	// https://github.com/alecthomas/chroma/tree/master/styles
	HighlightStyle string `yaml:"Highlight-Style"`

	// Create id= attributes for all headers automatically
	HeadingIDs bool `yaml:"Heading-IDs"`
}

// Indicates whether it's directory, a directory containing
// markdown files, or file, or a Markdown file.
// Used for bit flags
type MdOptions uint8

type dirInfo struct {
	mdOptions MdOptions
}

type social struct {
	DeviantArt string `yaml:"Deviant-Art"`
	Facebook   string `yaml:"Facebook"`
	GitHub     string `yaml:"GitHub"`
	Gitlab     string `yaml:"GitLab"`
	Instagram  string `yaml:"Instagram"`
	LinkedIn   string `yaml:"LinkedIn"`
	Pinterest  string `yaml:"Pinterest"`
	Reddit     string `yaml:"Reddit"`
	TikTok     string `yaml:"TikTok"`
	Tumblr     string `yaml:"Tumbler"`
	Twitter    string `yaml:"Twitter"`
	Weibo      string `yaml:"Weibo"`
	YouTube    string `yaml:"YouTube"`
}

// Everything relevant about the page to be published,
// namely its rendered text and what's in the front matter, but
// potentially also other stuff like file create date.
type WebPage struct {
	// Rendered text, the HTML after going through templates
	html []byte
}

// TODO: Document
const (
	// Known to be a directory with at least 1 Markdown file
	MarkdownDir MdOptions = 1 << iota

	// Known to be a filename with a Markdown extension
	MarkdownFile

	// Directory. Don't know yet if it contains Markdown files.
	NormalDir

	// File. Don't know if it's a markdown file.
	NormalFile

	// Set if directory has a file named "index.md", forced to lowercase
	HasIndexMd

	// Set if directory has a file named "README.md", case sensitive
	HasReadmeMd
)

// TODO: Document
func (a *App) setMdOption(dir string, mdOption MdOptions) {
	d := a.Site.dirs[dir]
	d.mdOptions = mdOption
	a.Site.dirs[dir] = d
}

// IsOptionSet returns true if the opt bit is set.
func (m MdOptions) IsOptionSet(opt MdOptions) bool {
	return m&opt != 0
}

// newSite() generates an empty site at
// the directory specified in pathname.
// It creates the site in a temp directory,
// then renames the temp directory to the desired name
// when completed.
func (app *App) newSite(pathname string) error {
  app.Print("newSite(%v)", pathname)
	// Exit if there's already a project at specified location.
	if isProject(pathname) {
		return ErrCode("0951", pathname)
	}

	// Change to specified directory.
	requested := pathname

	// First job is to create a temporary site in
	// the current directory. If anything goes wrong
	// it gets deleted. Extract the desired directory,
	// then replace the filename with a temp fileame.
	dir := filepath.Dir(pathname)

	// For case in which you already have the directory,
	// are inside it, but it's not an existing project.
	// For example, mkdir foo && cd foo && mb new site .
	inProjectDir := false
	if dir == "." /* || dir == ".." */ {
		dir = currDir()
    inProjectDir = true
		requested = filepath.Join(dir, pathname)
	}
	if requested == currDir() {
    inProjectDir = true
  }
  app.Note("inProjectDir: %v", inProjectDir) 
  var tmpDir string
	var err error
	// Create the temporary directory. It starts with the
	// Metabuzz product name abbreviation.
	if !inProjectDir {
		if tmpDir, err = os.MkdirTemp(dir, defaults.ProductShortName); err != nil {
			msg := fmt.Sprintf("%s for project %s: %s", dir, pathname, err.Error())
			return ErrCode("0414", msg)
		}
	}
	app.Debug("newSite(%v)", pathname)
	// Create a project at the specified path
	if !inProjectDir {
		err = os.MkdirAll(tmpDir, defaults.ProjectFilePermissions)
		if err != nil {
			return ErrCode("0401", tmpDir)
		}
		if err = app.changeWorkingDir(tmpDir); err != nil {
			return ErrCode("PREVIOUS", err.Error())
		}
	} else {
		if err = app.changeWorkingDir(requested); err != nil {
			return ErrCode("PREVIOUS", err.Error())
		}
  }
	// Copy files required to populate the .mb directory
	if inProjectDir {
		if err = app.copyMbDir(requested); err != nil {
			return ErrCode("PREVIOUS", err.Error())
		}
	} else {
		if err = app.copyMbDir(tmpDir); err != nil {
			return ErrCode("PREVIOUS", err.Error())
		}
	}

	// Get factory themes and copy to project. They will then
	// be copied on demand to the publish directory as needed.
	// This makes it easy to find themes and modify theme.
	// If user supplied a site configuration file, use it
	filename := ""
	if app.Flags.Site != "" {
		// User specified a site file such "--site foo.yaml",
		// so turn it into a fully qualified pathname if
		// no path is supplied
		dir := filepath.Dir(app.Flags.Site)
		if dir == ".." || dir == "." {
			// Only filename, e.g. "--site foo.yaml", no path, so add path
			if inProjectDir {
				filename = filepath.Join(requested, defaults.CfgDir, app.Flags.Site)
			} else {
				filename = filepath.Join(tmpDir, defaults.CfgDir, app.Flags.Site)
			}
		} else {
			// Fully qualified filename, e.g. "--site ~/foo.yaml",
			// so no need to add path
			filename = app.Flags.Site
		}
	} else {
		// No site file specified, so use default
		if inProjectDir {
			filename = filepath.Join(requested, defaults.CfgDir, defaults.SiteConfigFilename)
		} else {
			filename = filepath.Join(tmpDir, defaults.CfgDir, defaults.SiteConfigFilename)
		}
	}
	if err = app.writeSiteConfig(filename); err != nil {
		msg := fmt.Sprintf("%s: %s", filename, err.Error())
		return ErrCode("0227", msg)
	}
  if inProjectDir {
	  app.Site.path = requested
  } else {
	  app.Site.path = tmpDir
  }
	// Generate stub pages/sections if specified
	app.Site.name = filepath.Base(requested)
	if app.Flags.Starters != "" {
		// User specified a site file such "--site foo.yaml",
		// so turn it into a fully qualified pathname if
		// no path is supplied
		dir := filepath.Dir(app.Flags.Starters)
		if dir == ".." || dir == "." {
			// Only filename, e.g. "--site foo.yaml", no path, so add path
			filename = filepath.Join(filepath.Dir(requested), app.Flags.Starters)
		} else {
			// Fully qualified filename, e.g. "--site ~/foo.yaml",
			// so no need to add path
			filename = app.Flags.Starters
		}
		if err = app.generate(filename); err != nil {
			return ErrCode("PREVIOUS", err.Error())
		}
	} else {
		// if not just generate a very simple index page
		if err = app.createStubIndex(); err != nil {
			return ErrCode("PREVIOUS", err.Error())
		}
	}
	// Restore temp dir name to pathname  passed in.
	// Handle the case where the directory already exists.
  if !inProjectDir {
    if err = os.Rename(tmpDir, requested); err != nil {
      msg := fmt.Sprintf("%s to %s: %s", tmpDir, requested, err.Error())
      return ErrCode("0226", msg)

      // TODO: This can't execute.
      // It shoiuld run when os.Rename fails.
      // At the moment I can't see why Rename is failing but
      // I'm not getting an error
      if err := os.RemoveAll(tmpDir); err != nil {
        msg := fmt.Sprintf("System error attempting to delete temporary site directory %s: %s", tmpDir, err.Error())
        return ErrCode("0601", msg)
      }
    }
  }
	app.Site.path = requested
	if err := app.changeWorkingDir(requested); err != nil {
		msg := fmt.Sprintf("System error attempting to change to new site directory %s: %s", requested, err.Error())
		return ErrCode("1111", msg)
	}
	// xxx
	return nil
}
func fullPath(defaultDir string, filename string) string {
	dir := filepath.Dir(filename)
	if dir == ".." || dir == "." {
		filename = filepath.Join(defaultDir, filename)
	}
	return filename
}

// readSiteConfig() obtains site config info from the
// site configuration file, i.e. site.yaml
// Probably need to call Site.newSite first.
func (app *App) readSiteConfig() error {
	var err error
	var b []byte
	if app.Site.Filename == "" {
		return ErrCode("1063", "")
	}
	if b, err = ioutil.ReadFile(app.Site.Filename); err != nil {
		// TODO: Handle error properly & and document error code
		return ErrCode("PREVIOUS", err.Error(), app.Site.Filename)
	}

	err = yaml.Unmarshal(b, &app.Site)
	if err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}
	app.Debug("readSiteConfig(%v): Site is %#v", app.Site.Filename, app.Site)

	return nil
}

// writeSiteConfig() writes out the contents of App.Site
// Optional param is when a temp directory is passed in
// during site creation
func (app *App) writeSiteConfig(path ...string) error {
	var filename string
	if len(path) < 1 {
		filename = app.Site.Filename
	} else {
		filename = path[0]
	}
	if err := writeYamlFile(filename, app.Site); err != nil {
		// TODO: Better error handling?
		app.Note("writesiteconfig failed")
		return err
	}
	return nil
}

// setSiteDefaults() obtains starting values
// for a fresh Site object before it's written
// to a site config file, or initialized
// by another site config file.
func (app *App) setSiteDefaults() {
	app.Site.Language = defaults.Language
	app.Site.MarkdownOptions.HighlightStyle = defaults.ChromaDefault
	app.setPaths()
}

// copyMbDir() copies the .mb directory to the new site.
// Typically should go to app.Site.path
// But you can pass in the name of a temp dir too.
func (app *App) copyMbDir(params ...string) error {
	var path string
	if len(params) < 1 {
		path = app.Site.path
	} else {
		path = params[0]
	}
	app.Debug("\tcopyMbDir to %v", path)
	return app.embedDirCopy(mbfiles, path)
}
