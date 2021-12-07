package app

import (
	"embed"
	"github.com/tomcam/m/pkg/default"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
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

	// Full path to site config file
	siteFilePath string

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
// the directory specified in app.Site.path
func (app *App) newSite(pathname string) error {
	app.Debug("newSite(%v)", pathname)
	var err error
	// Create a project at the specified path
	// xxx
	err = os.MkdirAll(pathname, defaults.ProjectFilePermissions)
	if err != nil {
		return ErrCode("0401", pathname)
	}
	// Update app.Site.path and build all related directories
	if err = app.changeWorkingDir(pathname); err != nil {
		return ErrCode("PREVIOU", err.Error())
	}

	// Change to specified directory.
	pathname = app.Site.path

	// Exit if there's already a project at specified location.
	if isProject(pathname) {
		return ErrCode("0951", pathname)
	}

	// Copy files required to populate the .mb directory
	if err = app.copyMbDir(); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}

	// Get factory themes and copy to project. They will then
	// be copied on demand to the publish directory as needed.
	// This makes it easy to find themes and modify theme.
	filename := ""
	// If user supplied a site configuration file, use it
	if app.Flags.Site != "" {
		app.Debug("\t\tUse site file %v", app.Flags.Site)
		filename = app.Flags.Site
	}
	// TODO: Populate
	if err := app.writeSiteConfig(); err != nil {
		app.Debug("\t\tError after app.writeSiteConfig()")
		// TODO: Handle error properly & and document error code
		//return ErrCode("0220", err.Error(), filename)
		return ErrCode("PREVIOUS", err.Error(), filename)
	}

	// Generate stub pages/sections if specified
	if app.Flags.Starters != "" {
		if err := app.generate(app.Flags.Starters); err != nil {
			return ErrCode("PREVIOUS", err.Error())
		}
	} else {
    // if not just generate a very simple index page
 		if err := app.createStubIndex(); err != nil {
			return ErrCode("PREVIOUS", err.Error())
		}
 }
	return nil
}

// readSiteConfig() obtains site config info from the
// site configuration file, i.e. site.yaml
// Probably need to call Site.newSite first.
func (app *App) readSiteConfig() error {
	var err error
	var b []byte
	if app.Site.siteFilePath == "" {
		return ErrCode("1063", "")
	}
	if b, err = ioutil.ReadFile(app.Site.siteFilePath); err != nil {
		// TODO: Handle error properly & and document error code
		return ErrCode("PREVIOUS", err.Error(), app.Site.siteFilePath)
	}

	err = yaml.Unmarshal(b, &app.Site)
	if err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}
	app.Debug("readSiteConfig(%v): Site is %#v", app.Site.siteFilePath, app.Site)

	return nil
}

// writeSiteConfig() writes out the contents of App.Site
func (app *App) writeSiteConfig() error {
	app.Debug("writeSiteConfig()")
	if err := writeYamlFile(app.Site.siteFilePath, app.Site); err != nil {
		// TODO: Better error handling?
		return ErrCode("PREVIOUS", app.Site.siteFilePath, err.Error())
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
func (app *App) copyMbDir() error {
	app.Debug("\tcopyMbDir to %v", app.Site.path)
	return app.embedDirCopy(mbfiles, app.Site.path)
}
