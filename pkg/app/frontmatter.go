package app

import (
	//"fmt"
	"reflect"
)

type FrontMatter struct {
	// Theme specified by user
	Theme string `yaml:"Theme"`

	// Generates a Description metatag on output
	Description string `yaml:"Description"`

	// Date this document was created (separate from file date)
	//DateCreated time.Time `yaml:"DateCreated"`
	DateCreated string `yaml:"DateCreated"`

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

	// User data--MUST REMAIN AT END
	List interface{} `yaml:"List"`

}

// frontMatterToString generates the front matter
// section of a page in "sparse" format, that is,
// without a bunch of empty fields.
// So it might create something like this if called
// from a starter. Could have even fewer
// fields; simply depends on what nonempty values
// are in the FrontMatter struct.
//
//   ---
//   Theme: hero
//   Title: Assemble
//   Sidebar: left
//   ---
//
// Extract only the string fields with contents
// and include those, for example,
// FrontMatter.Theme or FrontMatter.Mode
// If nothing in the front matter is set, returns
// an empty string.
// Hmm... see https://stackoverflow.com/a/66511341
func frontMatterToString(f FrontMatter) string {
	fields := reflect.ValueOf(f)
	frontMatter := ""
	for i := 0; i < fields.NumField(); i++ {
		// Get name (key) of next FrontMatter struct field.
		fieldName := fields.Type().Field(i).Name
		contents := structFieldByNameStrMust(f, fieldName)
		//fmt.Printf("\t%v: %v\n", fieldName, contents)
		if contents != "" && contents != "[]" {
			// TODO: stringbuilder
			frontMatter += fieldName + ": " + contents + "\n"
		}
	}
	if frontMatter != "" {
		frontMatter = "---" + "\n" + frontMatter + "---" + "\n"
	}
	//fmt.Printf("\tfrontMatterToString results: %#v\n", frontMatter)

	return frontMatter
}
