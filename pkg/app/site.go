package app

import (
	"github.com/tomcam/m/pkg/default"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

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
	// Was assetDir
	assetPath string

	// Make it easy if you just have 1 author.
	Author author `yaml:"Author"`

	// List of authors with roles and websites in site.toml
	Authors []author `yaml:"Authors"`

	// Base directory for URL root, which may be different
	// from its actual root. For example, GitHub Pages prefers
	// the blog to start in /docs instead of root, but
	// a URL would omit it.
	// Was BaseDir
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
	// Was cssDir string
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

	// Google Analytics tracking ID specified in site.toml
	Ganalytics string `yaml:"Ganalytics"`

	// All these files are copied into the HTML header.
	// Example: favicon links.
	HeadTags []string `yaml:"Head-Tags"`

	// Full path of header tags for "code injection"
	headTagsPath string

	HTMLStartFile string `yaml:"HTML-start-file"`
	HTMLEndFile   string `yaml:"HTML-end-file"`
	// Subdirectory under the AssetDir where image files go
	// Was imageDir
	imagePath string

	// for HTML header, as in "en" or "fr"
	Language string `yaml:"Language"`

	// Flags indicating which non-CommonMark Markdown extensions to use
	markdownOptions MarkdownOptions

	// Mode ("dark" or "light") used by this site unless overridden in front matter
	Mode string `yaml:"Mode"`

	// Site's project name, so it's a filename.
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
	path string

	// Directory for finished site--rendered HTML & asset output
	// Was publishDir
	publishPath string

	// Full path of shortcode dir for this project. It's computed
	// at runtime using SCodeDir, also in this struct.
	sCodePath string

	// Full path of file containing JSON version of site text
	// to be indexed
	searchJSONFilePath string

	// Full path to site config file
	siteFilePath string

	// Social media URLs
	Social social `yaml:"Social"`

	// Site defaults to using this sidebar setting unless
	// a page specifies otherwise
	// Was DefaultSidebar
	Sidebar string `yaml:"Sidebar"`

	// Name (not path) of Theme used by this site unless overridden in front matter.
	// Was DefaultTheme
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
	webPages map[string]WebPage

	// THIS ALWAYS GOES AT THE END OF THE FILE/DATA STRUCTURE
	// User data.

	List interface{} `yaml:"List"`

	// Pages to generate when site is created
  Generate map[string]Stub `yaml:"Generate"`
}

type Stub struct {
  Title string `yaml:"Title"`
  Filename string `yaml:"Filename"`
  Description string `yaml:"Description"`
}

type company struct {
	// Company name, like "Metabuzz" or "Example Inc."
	Name string `yaml:"Name"`
	URL  string `yaml:"URL"`

	// Logo file for the header
	HeaderLogo string
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
	HardWraps bool `yaml:"Hard-wraps"`

	// Name of color scheme used for code highlighting,
	// for example, "monokai"
	// List of supported languages:
	// https://github.com/alecthomas/chroma/blob/master/README.md
	// I believe the list of themes is here:
	// https://github.com/alecthomas/chroma/tree/master/styles
	HighlightStyle string `yaml:"Highlight-style"`

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
func (a *App) addMdOption(dir string, mdOption MdOptions) {
	d := a.Site.dirs[dir]
	d.mdOptions |= mdOption
	a.Site.dirs[dir] = d
}

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

// createSite() generates an empty site at
// the directory specified in app.Site.path
func (app *App) createSite(pathname string) error {
	app.Debug("\tcreateSite(%v)", pathname)
	var err error
	// Create a project at the specified path
	err = os.MkdirAll(pathname, defaults.ProjectFilePermissions)
	if err != nil {
		return ErrCode("0401", pathname)
	}
	// Update app.Site.path and build all related directories
	if err = app.setWorkingDir(pathname); err != nil {
		return err
	}

	// Change to specified directory.
	pathname = app.Site.path

	// Exit if there's already a project at specified location.
	if isProject(pathname) {
		return ErrCode("0951", pathname)
	}

	// Create minimal directory structure: Publish directory
	// .site directory, .themes, etc.
	if err = createDirStructure(&defaults.SitePaths); err != nil {
		app.Debug("\t\tcreateDirStructure() failed during createSite()")
		return ErrCode("PREVIOUS", err.Error())
	}
	// Get factory themes and copy to project. They will then
	// be copied on demand to the publish directory as needed.
	// This makes it easy to find themes and modify theme.
	app.Debug("About to copy factory themes")
	//if err = app.copyFactoryThemes(); err != nil {
	err = app.copyFactoryThemes()
	app.Debug("\terr after calling app.copyFactoryThemes(): %v", err)
	if err != nil {
		// TODO: Improve error handling?
		app.Note("TODO: DUDE!!!")
		app.Debug("\t\tcopyFactoryThemes() failed during createSite()")
		return ErrCode("PREVIOUS", err.Error())

	}
	app.Note("Copied factory themes")

	// TODO: Populate
	if err := app.writeSiteConfig(); err != nil {
		// TODO: Handle error properly & and document error code
		app.Debug("Error writing site file %v", app.Site.siteFilePath)
		return ErrCode("PREVIOUS", err.Error())
	}
	return nil
}

// readSiteConfig() obtains site config info from the
// site configuration file, i.e. site.yaml
// Pre: call Site.newSite()
func (app *App) readSiteConfig() error {
	var err error
	var b []byte
	if b, err = ioutil.ReadFile(app.Site.siteFilePath); err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}

	err = yaml.Unmarshal(b, &app.Site)
	if err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}

	err = yaml.Unmarshal(b, &app.Site)
	if err != nil {
		// TODO: Handle error properly & and document error code
		return err
	}
	return nil
}

// writeSiteConfig() writes the contents of App.Site
// to .mb/site.yaml.
// and creates or replaces a TOML file in the
// project's site subdirectory.
// Assumes you're in the project directory.
func (app *App) writeSiteConfig() error {
	// Populate site with default values from config info.
	app.setSiteDefaults()
	return writeYamlFile(app.Site.siteFilePath, app.Site)
}


// generate() creates files specified in the 
// site file. It assumes a site has been created
// via createSite() and that we're in the 
// project site specified in app.Site.path
func (app *App) generate() error {
  pathname := app.Site.path
  app.readSiteConfig()
	app.Note("\tgenerate(%v)", pathname)
  app.Note("\tSite %#v\n\n", app.Site)
  app.Note("\tSite.Generate: %#v\n\n", app.Site.Generate)
	//var err error

	// Exit if there's already a project at specified location.
	if !isProject(pathname) {
		return ErrCode("1026", pathname)
	}
  for k, v := range app.Site.Generate {
    app.Note("%v: %v\n", k, v)
  }


	return nil
}

