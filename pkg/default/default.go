package defaults

import "github.com/tomcam/m/pkg/util"

const (
	Version = "0.4.0"
)

var (

	// Directory configuration for a project--a new site.
	SitePaths = [][]string{
		{CfgDir, DefaultPublishPath},
		{CfgDir, CommonPath},
		{CfgDir, HeadTagsPath},
		{CfgDir, SCodePath},
		{CfgDir, ScriptClosePath},
		{CfgDir, ScriptOpenPath},
		{CfgDir, ThemesDir},
	}
	// Markdown file extensions
	// They don't in lexical order because they
	// get sorted. That's because it's possible
	// more will be added via config file on startup.
	MarkdownExtensions = util.NewSearchInfo([]string{
		".Rmd",
		".markdown",
		".md",
		".mdown",
		".mdtext",
		".mdtxt",
		".mdwn",
		".mkd",
		".mkdn",
		".sidebar",
		".navbar",
		".header",
		".footer",
		".text"})

	// ExcludedAssetExtensions are the extensions of files in the asset
	// directory that should NOT be copied to the publish directory.
	// The contents of a theme directory mix both things you want copied,
	// like CSS files, and things you don't, like TOML or Markdown files.
	ExcludedAssetExtensions = util.NewSearchInfo([]string{
		".bak",
		".html",
		".toml",
		".yaml",
		".yml",
	})
)

const (
	// Name of the subdirectory that holds shared text.
	// Excluded from publishing.
	CommonPath = "common"

	// Tiny starter file for index.md
	IndexMd = `
# %s

Welcome to %s
`
	HTMLStartFile = "<!DOCTYPE html><HTML lang="
	HTMLEndFile   = `</body>` + "\n" + `</html>`

	// Name of subdirectory within the publish directory for CSS, theme files.
	// for that theme.
	DefaultAssetPath = "assets"

	// Name of the subdirectory containing files that get copied
	// into the header of each HTML file rendered by the site
	// Excluded from publishing.
	HeadTagsPath = "headtags"

	// Language used for `HTML lang=`
	Language = "en"

	// Name of the subdirectory the rendered files get rendered
	// to. It can't be changed because it's used to determine
	// whether a site is contained within its parent directory.
	// Excluded from publishing.
	DefaultPublishPath = "pub"

	// Name of subdirectory containing shortcode files
	// Excluded from publishing.
	SCodePath = "scodes"

	// Location of directory containing Javascript
	// that goes at the end of the HTML file, near
	// the closing <body> tag.
	// The files MUST supply <script> tags.
	// It is possible that somehting other
	// than Javascript will be used.
	ScriptClosePath = "scriptclose"

	// Location of directory containing Javascript
	// that goes at the begining of the HTML file, near
	// the opening <body> tag.
	// The files MUST <script> tags.
	ScriptOpenPath = "scriptopen"

	// Name of subdirectory within the theme that holds help & sample files
	// for that theme.
	ThemeHelpPath = ".help"

	// Name of subdirectory under the publish directory for style sheet assets
	// (Currently not well thought out nor in use, though assets directory is
	// being used)
	DefaultPublishCssPath = "css"

	// Name of subdirectory under the publish directory for image assets
	// (Currently not well thought out nor in use, though assets directory is
	// being used)
	DefaultPublishImgPath = "img"

	// Use this theme if none is specified, and also
	// as the theme used to generate new themes if
	// not otherwise specified.
	DefaultThemeName = "wide"

	// Name of the directory that holds items used by projects, such
	// as themes and shortcodes.
	// TODO: Change this when I settle on a product name
	// TC: Formerly GlobalConfigurationDirName.
	CfgDir = ".mb"

	// Default file extension used by configuration files.
	// See https://yaml.org/faq.html
	ConfigFileDefaultExt = "yaml"

	// A configuration file passed to the command line.
	ConfigFilename = ProductName + "." + ConfigFileDefaultExt

	// The configuration file in the user's application
	// data directory, without the path.
	AppDataConfigFilename = ProductName + "." + ConfigFileDefaultExt

	// The local configuration file name without the path.
	LocalConfigFilename = ProductName + "." + ConfigFileDefaultExt

	// Name of file containing .JSON database of text used for
	// search purposes.
	SearchJSONFilename = ProductName + "-" + "search" + ".json"

	// A starters configuration file passed to the command line.
	ConfigStartersFilename = "starters" + "." + ConfigFileDefaultExt

	// By default, the published site gets its theme from a local copy
	// within the site directory--this directory.
	// It then copies from that copy to
	// generate pages in the Publish directory. Helps prevent unintended changes
	// from being made to the originals, and makes it much easier to
	// make theme changes, especially if you're a noob or just want to
	// type less.
	ThemesDir = "themes"

	// These are the themes named in the front matter
	// get copied from the global factory themes directory
	// to this particular site.
	SiteThemesDir = "themes"

	// Configuration file found in the current site source directory
	SourcePathConfigFilename = ProductName + "." + ConfigFileDefaultExt

	// Actual colloquial name for this product
	// but used in directories & other
	// purposes, like storing config files.
	// Make it in lowercase. One word,
	// like docset or metabuzz.
	// If this changes update CfgDir
	ProductName = "metabuzz"

	// Name of the product as it would appear in published copy.
	// Capitalization, spaces, special characters, etc.
	// are respected
	ProductBranding = "Metabuzz" + "™ " + "v." + Version

	// Abbreviation, used for name command line program.
	ProductShortName = "mb"

	// Values set through the environment as opposed to config files
	// or command line options.
	// A short version of the product name
	// used as a prefix for environment variables.
	ProductEnvPrefix = "MBZ_"

	// The permissions given to output files, and also to
	// configuration files.
	// 0755 means the owner can read, write and execute (first 7)
	// Other people can only read (5 and 5). That makes sense
	// for a web server
	PublicFilePermissions = 0755

	// Objects that must be visible to the project, but not the public
	ProjectFilePermissions = 0700

	// Name of the file that holds site configuration information
	SiteConfigFilename = "site" + "." + ConfigFileDefaultExt

	// String that precedes error codes
	ErrorCodePrefix = "mbz"
)
