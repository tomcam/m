package app

// type Page contains read-only information about the Markdown page currently
// being processed.
type Page struct {
	// TODO: Load in front matter as a real struct
	frontMatterRaw map[string]interface{}
	theme          Theme
  // Location of source theme files computed at
  // runtime
  themePath string
}

type FrontMatter struct {
	// Theme specified by user
	Theme string `json:"name"`

	// Generates a Description metatag on output
	Description string `json:"description"`

	// Generates a Title tag on output
	Title string `json:"title"`

	// Determine whether aside is on the
	// right, left, or none
	Sidebar string `json:"sidebar"`
}
