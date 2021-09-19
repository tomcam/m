package defaults

import "github.com/tomcam/mb/pkg/slices"

var (
	Version = ProductName + " version " +
		"0.2.0"

	// Directory configuration for a project--a new site.
	SiteDirs = [][]string{
		{PublishDir},
		{GlobalConfigurationDirName, CommonDir},
		{GlobalConfigurationDirName, HeadTagsDir},
		{GlobalConfigurationDirName, SCodeDir},
		{GlobalConfigurationDirName, ScriptCloseDir},
		{GlobalConfigurationDirName, ScriptOpenDir},
		{GlobalConfigurationDirName, SiteConfigDir},
		{GlobalConfigurationDirName, ThemeDir},
	}
	// Markdown file extensions
	MarkdownExtensions = slices.NewSearchInfo([]string{
		".Rmd",
		".md",
		".mdown",
		".mdtext",
		".mdtxt",
		".markdown",
		".mdwn",
		".mkd",
		".mkdn",
		".text"})

	// ExcludedAssetExtensions are the extensions of files in the asset
	// directory that should NOT be copied to the publish directory.
	// The contents of a theme directory mix both things you want copied,
	// like CSS files, and things you don't, like TOML or Markdown files.
	ExcludedAssetExtensions = slices.NewSearchInfo([]string{
		".html",
		".toml",
		".bak",
	})
)

const (
	// Name of the subdirectory that holds shared text.
	// Excluded from publishing.
	CommonDir = "common"

	// Tiny starter file for index.md
	IndexMd = `
# %s

Welcome to %s
`

	// Name of subdirectory within the publish directory for CSS, theme files.
	// for that theme.
	DefaultAssetDir = "assets"

	// Name of the subdirectory the rendered files get rendered
	// to. It can't be changed because it's used to determine
	// whether a site is contained within its parent directory.
	// Excluded from publishing.
	PublishDir = ".pub"

	// Name of the subdirectory containing files that get copied
	// into the header of each HTML file rendered by the site
	// Excluded from publishing.
	HeadTagsDir = "headtags"

	// Name of subdirectory containing shortcode files
	// Excluded from publishing.
	SCodeDir = "scodes"

  // Location of directory containing Javascript 
  // that goes at the end of the HTML file, near
  // the closing <body> tag.
  // The files MUST supply <script> tags.
  // It is possible that somehting other
  // than Javascript will be used. 
  ScriptCloseDir = "scriptclose" 

  // Location of directory containing Javascript 
  // that goes at the begining of the HTML file, near
  // the opening <body> tag.
  // The files MUST <script> tags.
  ScriptOpenDir = "scriptopen"

	// Name of subdirectory within the theme that holds help & sample files
	// for that theme.
	ThemeHelpSubdirname = ".help"

	// Name of subdirectory under the publish directory for style sheet assets
	// (Currently not well thought out nor in use, though assets directory is
	// being used)
	DefaultPublishCssSubdirname = "css"

	// Name of subdirectory under the publish directory for image assets
	// (Currently not well thought out nor in use, though assets directory is
	// being used)
	DefaultPublishImgSubdirname = "img"

	// Name of theme used to create other themes, or theme to be
	// used if for some reason no others are present
	DefaultThemeName = "wide"

	// Name of the directory that holds items used by projects, such
	// as themes and shortcodes.
	// TODO: Change this when I settle on a product name
	GlobalConfigurationDirName = ".mb"

	// Default file extension used by configuration files.
	ConfigFileDefaultExt = "toml"

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

	// By default, the published site gets its theme from a local copy
	// within the site directory. It then copies from that copy to
	// generate pages in the Publish directory. Helps prevent unintended changes
	// from being made to the originals, and makes it much easier to
	// make theme changes, especially if you're a noob or just want to
	// type less.
	ThemeDir = ".themes"

	// Configuration file found in the current site source directory
	SourceDirConfigFilename = ProductName + "." + ConfigFileDefaultExt

	// Actual colloquial name for this product
	// but used in directories & other
	// purposes, like storing config files.
	// Make it in lowercase. One word,
	// like docset or metabuzz.
	// TODO: If this changes update GLOBAL_CONFIG_DIRNAME
	// TODO: Change this when I settle on a product name, and also change PRODUCT_SHORT_NAME
	ProductName = "metabuzz"

	// Abbreviation, used for name command line program.
	ProductShortName = "mb"

	// Values set through the environment as opposed to config files
	// or command line options.
	// A short version of the product name
	// used as a prefix for environment variables.
	// TODO: Change this when I settle on a product name
	ProductEnvPrefix = "MBZ_"
	// Examples:
	// PRODUCT_ENV_PREFIX+"DEFAULT_THEME"
	// PRODUCT_ENV_PREFIX+"SC_DIR"

	// The permissions given to output files, and also to
	// configuration files.
	// 0755 means the owner can read, write and execute (first 7)
	// Other people can only read (5 and 5). That makes sense
	// for a web server
	PublicFilePermissions = 0755

	// Objects that must be visible to the project, but not the public
	ProjectFilePermissions = 0700

	// Name of the subdirectory in the project where the site info is held.
	// That includes the site.toml file and also the publish directory.
	SiteConfigDir = "site"

	// Name of the file that holds site configuration information
	SiteConfigFilename = "site" + "." + ConfigFileDefaultExt

	// String that precedes error codes
	ErrorCodePrefix = "mbz"
)
