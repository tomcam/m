# The Site Configuration file

The site configuration named, normally named `site.yaml` holds settings used by your whole site, for example, your company name and URL. It's in the `.mb` subirectory of your project. In some operating systems, for example, MacOS and Linux, it's hidden but you can easily edit it without doing anything special.

Metabuzz does some things for you automatically if you've filled out the site file settings. Depending on the theme you're using, just doing that may completely eliminate any need to touch the header, footer, and navbar files.

It's explained below but you can see an example site file
at [Example site file](site-file-example.html)

## TODO: Finish
[Branding](#branding): the site's full name  
[Company](#company): Your Organization name and logo  
[Theme](#theme): Theme used by all pages of the site unless specified otherwise in the page's [front matter](front-matter.html)
[Sidebar](#sidebar): Sidebar setting used by all pages of the site unless specified otherwise in the page's [front matter](front-matter.html)
[ExcludeExtensions](#exclude-extensions): Prevent types of files from being copied to the publish directory

<a id="theme"></a>
## Theme
Theme used by all pages of the site unless specified otherwise in the page's [front matter](front-matter.html#theme)

Example:

```
Theme: genuine
```
<a id="sidebar"></a>
## Sidebar

Sidebar setting used by all pages of the site unless specified otherwise in the page's [front matter](front-matter.html#sidebar)


<a id="author"></a>
## Author
Credits the person who originated and maintains the document.

<a id="branding"></a>
## Branding: the site's full name

Your site name is its url, for example, metabuzz.com. The Branding string lets you give it a fuller name. So, for example, if the name is `my-project` this might be `My Insanely Cool Project`

##### Example

```
Branding: My Insanely Cool Project
```

<a id="company"></a>
## Company section: Company name and logo

##### Example

```
Company:
```

<a id= "exclude-dirs"></a>

## ExcludeDirs: List of directories to prevent from being published

List of directories in the source project directory that should be
excluded, things like ".git" and "node_modules".

Normally if Metabuzz sees a directory, it assumes that the directory contains Markdown documents and copies that directory to the to the [publish directory](publish-directory). `ExcludeDirs` prevents that copy.

### Note: Dot directories are always ignored

Metabuzz excludes all directories with names that start with a "." period/dot character.

