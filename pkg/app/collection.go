package app

// Repreents
type Collection struct {
	// Format used to generate a path to the post.
	Permalink string `yaml:"Permalink"`

	// Directory path to the post.
	// Derives its formatting from permalink.
	Path string `yaml:"Path"`

	// Sort directory by filename, date, whatever
	Sort string `yaml:"Sort"`
}
