package app

// type Page contains read-only information about the Markdown page currently
// being processed.
type Page struct {
	// TODO: Load in front matter as a real struct
	frontMatter map[string]interface{}
	// Currently loaded theme
	//Theme Theme
}
