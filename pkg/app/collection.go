package app

// Repreents
type Collection struct {
	// Format used to generate a path to the post.
	permalink string

	// Directory path to the post.
	// Derives its formatting from permalink.
	path string
}
