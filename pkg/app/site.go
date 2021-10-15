package app

import (
	"github.com/tomcam/m/pkg/default"
	"os"
	"path/filepath"
)

// Allocates a Site object, leaving everything clean & empty
func (site *Site) New() (*Site, error) {
	s := Site{}
	return &s, nil
}

func (site *Site) NewSite() error {
	if site.path != "" {
		// Change to the specified directory.
		if err := os.Chdir(site.path); err != nil {
			//app.QuitError(ErrCode("1101", err.Error(), pathname))
			return ErrCode("0902", err.Error())
		}
	}

	// Create minimal directory structure: Publish directory
	// .site directory, .themes, etc.
	if err := createDirStructure(&defaults.SitePaths); err != nil {
		//return errs.ErrCode("PREVIOUS", err.Error())
		return ErrCode("PREVIOUS", err.Error())
	}
	return nil
}

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
	Author author

	// List of authors with roles and websites in site.toml
	Authors []author

	// Base directory for URL root, which may be different
	// from its actual root. For example, GitHub Pages prefers
	// the blog to start in /docs instead of root, but
	// a URL would omit it.
	// Was BaseDir
	BasePath string

	// Site's branding, any string, that user specifies in site.toml.
	// So, for example, if the Name is 'my-project' this might be
	// 'My Insanely Cool Project"
	Branding string

	// Full pathname of common directory. Derived from CommonSubDir
	commonPath string

	// Company name & other info user specifies in site.toml
	Company company

	// Subdirectory under the AssetDir where CSS files go
	// Was cssDir string
	cssPath string

	// List of all directories in the site
	dirs map[string]dirInfo

	// List of directories in the source project directory that should be
	// excluded, things like ".git" and "node_modules".
	// Note that directory names starting with a "." are excluded too.
	// DO NOT ACCESS DIRECTLY:
	// Use excludedDirs() because it applies other information such as PublishDir()
	ExcludeDirs []string

	// List of file extensions to exclude. For example. [ ".css" ".go" ".php" ]
	ExcludeExtensions []string

  // Number of markdown files processed
	fileCount   uint

	// Google Analytics tracking ID specified in site.toml
	Ganalytics string

	// All these files are copied into the HTML header.
	// Example: favicon links.
	HeadTags []string

	// Full path of header tags for "code injection"
	headTagsPath string

	// Subdirectory under the AssetDir where image files go
	// Was imageDir
	imagePath string

	// for HTML header, as in "en" or "fr"
	Language string

	// Flags indicating which non-CommonMark Markdown extensions to use
	markdownOptions MarkdownOptions

	// Mode ("dark" or "light") used by this site unless overridden in front matter
	Mode string

	// Site's project name, so it's a filename.
	// It's an identifier so it should be in slug format:
	// Preferably just alphanumerics, underline or hyphen, and
	// no spaces, for example, 'my-project'
	name string

	// Fullly qualified pathname of output file, e.g. index.html
	outfile string

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
	// xxx
	searchJSONFilePath string

	// Full path to site config file
	siteFilePath string

	// Social media URLs
	Social social

	// Site defaults to using this sidebar setting unless
	// a page specifies otherwise
	DefaultSidebar string

	// Name (not path) of Theme used by this site unless overridden in front matter.
	DefaultTheme string

	// Directory this site uses to copy themes from. If the -d option was
	// set, use the global factory themes directory. Otherwise, use local copy
	themesPath string

	// All the rendered pages on the site, plus meta information.
	// Index by the fully qualified path name of the source .md file.
	webPages map[string]WebPage

	// THIS ALWAYS GOES AT THE END OF THE FILE/DATA STRUCTURE
	// User data.
	List interface{}
}

// defaults() intializes the Site object
func (site *Site) defaults(home string) {
	site.setPaths(home)
}

type company struct {
	// Company name, like "Metabuzz" or "Example Inc."
	Name string
	URL  string

	// Logo file for the header
	HeaderLogo string
}
type author struct {
	FullName string
	URL      string
	Role     string
}

type authors struct {
	Authors []author
}

// MarkdownOptions consists of goldmark
// options used when converting markdown to HTML.
type MarkdownOptions struct {
	// If true, line breaks are significant
	hardWraps bool

	// Name of color scheme used for code highlighting,
	// for example, "monokai"
	// List of supported languages:
	// https://github.com/alecthomas/chroma/blob/master/README.md
	// I believe the list of themes is here:
	// https://github.com/alecthomas/chroma/tree/master/styles
	HighlightStyle string

	// Create id= attributes for all headers automatically
	HeadingIDs bool
}

// Indicates whether it's directory, a directory containing
// markdown files, or file, or a Markdown file.
// Used for bit flags
type MdOptions uint8

type dirInfo struct {
	mdOptions MdOptions
}
type social struct {
	DeviantArt string
	Facebook   string
	GitHub     string
	Gitlab     string
	Instagram  string
	LinkedIn   string
	Pinterest  string
	Reddit     string
	Tumblr     string
	Twitter    string
	Weibo      string
	YouTube    string
}

// Everything relevant about the page to be published,
// namely its rendered text and what's in the front matter, but
// potentially also other stuff like file create date.
type WebPage struct {
	// Rendered text, the HTML after going through templates
	html []byte
}

// setPaths computes values for location of site
// theme files, publish directory, etc.
// Most of them are relative to the site directory.
// It must be called after command line flags, env
// variables, and other application configuration has been done.
// home is the fully qualified directory name
// the project lives in.
// By now it is the current directory.
// This is based on App.SiteDefaults() in the previous
// version of Metabuzz.
func (site *Site) setPaths(home string) {
	// This is the fully qualified path of the current
	// directory, which is also guaranteed to be the
	// root directory of the project.
	site.path = home

  cfgPath := filepath.Join(site.path, defaults.CfgPath)

	// Build up the fully qualified pathname of the site file.
	site.siteFilePath = filepath.Join(cfgPath,
		defaults.SiteConfigFilename)

	// Build up the fully qualified pathname of the publish 
  // directory.
	site.publishPath = filepath.Join(cfgPath,
		defaults.PublishPath)

  // Create a new, empty map to hold the
  // source directory tree.
  site.dirs = make(map[string]dirInfo)
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
	d := a.site.dirs[dir]
	d.mdOptions |= mdOption
	a.site.dirs[dir] = d
}

// TODO: Document
func (a *App) setMdOption(dir string, mdOption MdOptions) {
	d := a.site.dirs[dir]
	d.mdOptions = mdOption
	a.site.dirs[dir] = d
}

// IsOptionSet returns true if the opt bit is set.
func (m MdOptions) IsOptionSet(opt MdOptions) bool {
	return m&opt != 0
}


