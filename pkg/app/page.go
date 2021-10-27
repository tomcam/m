package app

// type Page contains read-only information about the Markdown page currently
// being processed.
type Page struct {
	// TODO: Marshal in front matter as a real struct
	frontMatterRaw map[string]interface{}
	frontMatter    FrontMatter
	theme          Theme
	// Location of source theme files computed at
	// runtime
	themePath string
	// List of stylesheets actually published
	// (for example, only sidebar-left.css
	// or sidebar-right.css will be published)
	stylesheets []string
}

type FrontMatter struct {
	// Theme specified by user
	Theme string `json:"name"`

	// Generates a Description metatag on output
	Description string `json:"description"`

	// Generates a Title tag on output
	Title string `json:"title"`

	// If Mode is "dark", use a dark theme.
	Mode string

	// Determine whether aside is on the
	// right, left, or none
	Sidebar string `json:"sidebar"`
}
