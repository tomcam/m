package app

import (
	"time"
)

// type Page contains read-only information about the Markdown page currently
// being processed.
type Page struct {
	// Directory location of this page
	dir string

	// TODO: Marshal in front matter as a real struct
	frontMatterRaw map[string]interface{}
	FrontMatter    FrontMatter

	// Theme used by this page.
	Theme Theme

	// In case of a nested theme, such as "debut/gallery"
	// or "debut/gallery/item", include them all here.
	themes []Theme

	// Fully qualified filename of this source file
	filePath string

	// List of stylesheets actually published
	// (for example, only sidebar-left.css
	// or sidebar-right.css will be published)
	stylesheets []string

	// List of stylesheets with full path designations and
	// enclosed in HTML stylesheet tags.
	stylesheetTags string
}

type FrontMatter struct {
	// Theme specified by user
	Theme string `yaml:"Theme"`

	// Generates a Description metatag on output
	Description string `yaml:"Description"`

	// Date this document was created
	Created time.Time `yaml:"Created"`

	// Filenames to skip when publishing a theme
	ExcludeFiles []string `yaml:"ExcludeFiles"`

	// Generates a Title tag on output
	Title string `yaml:"Title"`

	// If Mode is "dark", use a dark theme.
	Mode string `yaml:"Mode"`

	// Disable features as needed on a per-page basis
	Suppress string `yaml:"Suppress"`

	// Determine whether sidebasr is on the
	// "right", "left", or "none" on per-page basis
	Sidebar string `yaml:"Sidebar"`

	// If set to "off", don't execute templates on this page.
	// Used for documentation purposes.
	Templates bool `yaml:"Templates"`
}
