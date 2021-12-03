package app

// type Page contains read-only information about the Markdown page currently
// being processed.
type Page struct {
	// Directory location of this page
	dir string

	// TODO: Marshal in front matter as a real struct
	frontMatterRaw map[string]interface{}
	FrontMatter    FrontMatter
	Theme          Theme

	// Fully qualified filename of this source file
	filePath string

	// List of stylesheets actually published
	// (for example, only sidebar-left.css
	// or sidebar-right.css will be published)
	stylesheets []string
}

type FrontMatter struct {
	// Theme specified by user
	Theme string `json:"theme"`

	// Generates a Description metatag on output
	Description string `json:"description"`

	// Filenames to skip when publishing a theme
	ExcludeFiles []string `json:"exclude-files"`

	// Generates a Title tag on output
	Title string `json:"title"`

	// If Mode is "dark", use a dark theme.
	Mode string

	// Determine whether aside is on the
	// right, left, or none
	Sidebar string `json:"sidebar"`

	// If set to "off", don't execute templates on this page.
	// Used for documentation purposes.
	Templates string
}
