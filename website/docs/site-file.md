# The Site Configuration file, aka site.toml

The site configuration named `site.toml` holds settings used by your whole site, for example, your company name and URL. It's in the `.site` subirectory of your project. In some operating systems, for example, MacOS and Linux, it's hidden but you can easily edit it without doing anything special.

Metabuzz does some things for you automatically if you've filled out the site file settings. Depending on the theme you're using, just doing that may completely eliminate any need to touch the header, footer, and navbar files.

It's explained below but you can see an example site file
at [Example site file](site-file-example.html)

[[[Author]]](#author): The document's primary creators and maintainer
[[[CoAuthors]]](#coauthors): The document's additional creators and maintainers  
[Branding](#branding): the site's full name  
[[Company]](#company): Your Organization name and logo  
[DefaultTheme](#defaulttheme): Theme used by all pages of the site unless specified otherwise in the page's [front matter](front-matter.html)
[DefaultSidebar](#defaultsidebar): Sidebar setting used by all pages of the site unless specified otherwise in the page's [front matter](front-matter.html)
[ExcludeExtensions](#exclude-extensions): Prevent types of files from being copied to the publish directory

## DefaultTheme
Theme used by all pages of the site unless specified otherwise in the page's [front matter](front-matter.html#theme)

## DefaultSidebar

Sidebar setting used by all pages of the site unless specified otherwise in the page's [front matter](front-matter.html#sidebar)

[[Social] section: Social media URLs](#social)  


{{- /* Not ready for documentation 
	AssetDir string
	BaseDir string
*/ -}}


## [[Author]] section: the document's primary creator

<a id="author"></a>

Credits the person who originated and maintains the document.

A list of 0 or more authors. Because there can be any number of them use `[[Authors]]` instead of `[Authors]`.
`author` is the full name of the document creator.  
`URL` is the web address at which the author can be found.  
`Role` describes the author's part in the document's creation. 

## [[Authors]] section: the document creator(s)

Credits the person who originated and maintains the document.  A list of 0 or more authors. Because there can be any number of them use `[[Authors]]` instead of `[Authors]`.
`author` is the full name of the document creator.  

`URL` is the web address at which the author can be found.  
`Role` describes the author's part in the document's creation. 


## [[CoAuthors]] section: the document's secondary creators and maintainers

##### file .site/site.toml

```
[[Authors]]
FullName = "Tom Campbell"
URL = "https://Metabuzz.com"
Role = "Dev lead"

[[Authors]]
FullName = "John Woloschuk"
URL = "https://klaatu.org"
Role = "Composer"

[[Authors]]
FullName = "Ryan Flynn"
URL = "https://brickandchrome.com"
Role = "Photographer"

```

##### output

```
Show authors as a list: [{Tom Campbell https://Metabuzz.com Dev lead} {John Woloschuk https://klaatu.org Resident Genius} {Ryan Flynn https://brickandchrome.com Photographer}]

First author in the list: {Tom Campbell https://Metabuzz.com Dev lead}

Enumerate list of authors and precede each with an asterisk:

* {Tom Campbell https://Metabuzz.com Dev lead}  
* {John Woloschuk https://klaatu.org Resident Genius}  
* {Ryan Flynn https://brickandchrome.com Photographer}  

```
<a id="branding"></a>
## Branding: the site's full name

Your site name is its url, for example, metabuzz.com. The Branding string lets you give it a fuller name. So, for example, if the name is `my-project` this might be `My Insanely Cool Project`

##### Example

```
branding = "My Insanely Cool Project"
```

<a id="company"></a>
## [Company] section: Company name and logo

##### Example

```
[company]
name = "METABUZZ"
HeaderLogo = "metabuzz-red.svg"
URL = "https://metabuzz.com"
```

These three items are typically used in image links or when your organization's logo gets displayed.

`name` is expected to be one word because it's sometimes used as a filename for logos or URLs. 

`HeaderLogo` is the name of an image file and is expected to be used for image files in image links.

`URL` is the organization's website address.

<a id= "exclude-dirs"></a>

## ExcludeDirs: List of directories to prevent from being published

List of directories in the source project directory that should be
excluded, things like ".git" and "node_modules".

Normally if Metabuzz sees a directory, it assumes that the directory contains Markdown documents and copies that directory to the to the [publish directory](publish-directory). `ExcludeDirs` prevents that copy.

### Note: Dot directories are always ignored

Metabuzz excludes all directories with names that start with a "." period/dot character.

##### Example

```
ExcludeDirs=[ "fullresimages", "node_modules" ]
```
<a id= "exclude-extensions "></a>

## ExcludeExtensions: Prevent types of files from being copied to the publish directory

`ExcludeExtensions` allows you to prevent classes of files by file extension from being copied to the [publish directory](publish-directory.md). Give it a list of extensions including the dot character as shown below.

##### Example

```
ExcludeExtensions=[ ".bak" ".go" ".php" ]
```

	// Google Analytics tracking ID
	Ganalytics string

	// All these files are copied into the HTML header.
	// Example: favicon links.
	Headers string

	// Language tag for html lang=
	Language string

	// Mode ("dark" or "light") used by this site unless overridden in front matter
	Mode string

	// Directory for finished site--rendered HTML & asset output
	Publish string

<a id="social"></a>
## [Social] section: Social media URLs

The optional `[Social]` part of your site file has slots for most popular social media sites. These examples all use the stunningly creative name `replace_this`, which of course you would replace with your own profile name.

```
[Social]
deviantart="https://www.deviantart.com/replace_this"
facebook="https://www.facebook.com/replace_this"
github="https://github.com/tomcam/replace_this"
gitlab="https://gitlab.com/replace_this"
instagram="https://www.instagram.com/replace_this/"
linkedin="https://www.linkedin.com/in/replace_this/"
twitter="https://twitter.com/replace_this"
pinterest="wwww.pinterest.com/replace_this"
reddit="https://reddit.com/u/replace_this"
tumblr="https://replace_thils.tumblr.com
twitter="twitter.com/replace_this"
weibo="https://www.weibo.com/replace_this"
youtube = "https://youtube.com/replace_this"
```



	// Name (not path) of Theme used by this site unless overridden in front matter.
	Theme string

	// ThemeDir is where all the themes are stored.
	// It is computed at startup based on configuration values.
	ThemeDir string



}


	// THIS ALWAYS GOES AT THE END OF THE FILE/DATA STRUCTURE
	// User data.
	List interface{}



type companyConfig struct {
	// Company name, like "Metabuzz" or "Example Inc."
	Name string
	URL  string

	// Logo file for the header
	HeaderLogo string
}
type authorConfig struct {
	FullName string
	URL      string
}



## List

TODO: Explain it needs to be last




