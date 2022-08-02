# Front matter

A Metabuzz markdown file may

## Example file
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

```
Description: "hello, world."
## excludedFiles

List of files in the current directory you don't want
copied to the Publish directory.
Must be literal filenames, not wildcards.

## TODO: Check this

```
excludedfiles = [ "clientid.src", "productkey.txt" ]
```

## Theme
Allows you to set the visual appearance on a per-page basis
by naming a Metabuzz theme. If you don't
name a theme, Metabuzz looks for a default theme
set in the [site file](site-file.html#defaulttheme). If
it can't find one there, it uses the 
[default theme](themes.html#default-theme).


## templates

## TODO: Verify

For documentation purposes. If you're writing documentation that uses the template language, setting `templates="off"` prevents templates on that page from
being executed. Helps when you're documenting, well, templates

```
Templates: false
```

For example, since there's no template function called `world` this
would normally produce an [0917]error if used in your Markdown, but if you 
set `Templates: false` you won't have that problem.
```
hello, {{ world. }}
```
.
