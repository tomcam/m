# To do

## Updating themes
* Old theme directory
https://github.com/tomcam/mb/tree/master/.mb/.themes
* Delete themename.toml
* Copy either w/layout.css or pillar/layout.css
* add --branding-weight to sizes.css
* Update theme-light.css, theme-dark.css
    /* Same as article > h1 { color: } */
    --article-h1-fg:var(--fg);
    --article-h2-fg:var(--fg);
    --article-h3-fg:var(--fg);
    --article-h4-fg:var(--fg);
    --article-h5-fg:var(--fg);
    --article-h6-fg:var(--fg);

    /* Same as article > h1 { background-color: } */
    --article-h1-bg:var(--bg);
    --article-h2-bg:var(--bg);
    --article-h3-bg:var(--bg);
    --article-h4-bg:var(--bg);
    --article-h5-bg:var(--bg);
    --article-h6-bg:var(--bg);


* YAML file Format looks like this:
Branding: "W by Metabuzz"
Description: "Minimal wide theme"
Stylesheets: 
- "reset.css"
- "fonts.css"
- "layout.css"
- "bind.css"
- "sizes.css"
- "theme-light.css"
- "w.css"
- "responsive.css"
Nav: {File: nav.md, HTML:}
Header: {File: header.md, HTML:}
Article: {File:, HTML:}
Sidebar: {File: sidebar.md, HTML:}
Footer: {File: footer.md, HTML:}
Language: en

## Priority 1: Showstoppers--required for the next release
* Bug: pillar.css is infested with colors and may even have trim-fg etc
* Bug: delete a theme directory, then try to build wih that theme.
you get this error. Handle that errcode condition.
Error building Can't find a theme named /Users/tom/code/m/cmd/mb/theme-test/.mb/themes/simplify (error code mbz1028)
 (error code mbz0923)

* Bug in sizes.css: I had got rid of --text-start.
  --sidebar-padding-left:var(--text-start);
* BUG: errors aren't gettingreported correctly, though they seem to work OK fi the extra
parameter is empty
* BUG: .Page.FrontMatter.Theme doesn't work correctly in an article, instead yielding asterisks
* Bug: Site.HighlightStyle aka Highlight-style  doesn't seem to work
* Add: toc
* Add: search
* Add: scriptclose directory
* Add: idea of post and specfications like YYYY-MM-DD or y-m-d etc, using dirs or strings as needed . That way mb new post "/blog/avengers review" would expand to something like "/blog/2022/04/21/avengers-review.html" or "/blog/2022-March-1-avengers-review.html" and so on 
* Make these changes when returning to the standard mb directory
  - Fix hardcoded paths in the files `gf`, `bu`
* Understand whether I should create empty index.html files for dirs
that don't have anything else to oprevent directory traversla attacks
* Add: dark themes
* Search and replqce 
  - TODO: in source
  - almost all QuitError calls because of a future interactive version
  - app.Print
  - app.Note

### Testing
* Should be able to create a project with a leading . in the name, but
that project should not publish files inside one if its subdirectoris with
a leading . in the name
* Deformed YAML file

## Credits
* https://iconduck.com/icons/84801/tools
* https://iconduck.com/icons/22486/tools

## Marketing
* Promise that in early versions while data stuctures and features may change, the
source file directory structure remains sacrosanct: a tree of Markdown files


## Document
* Document. Given a site.yaml with this:
Company:
    Name: "Eastside Emerald Home Repair"
    address: ""
    address2: ""
    city: ""
    country: ""
    postalcode: ""
    URL: ""
    HeaderLogo: "eastside-emerald-64x64.png"
You can do a header like this, which gives you an optional logo
{{- if .Site.Company.Name -}}
{{- $name := .Site.Company.Name -}}
{{- if .Site.Company.HeaderLogo -}}
* ![Logo](.Site.Copmany.HeaderLogo) [{{ $name -}}](/)
{{- else -}}
* [{{ $name -}}](/)
{{- end }} 
* [Services](services.html)
* [Rates](rates.html)
* [Contact](contact.html)
* [About](about.html)
* Create a page with this intentional error (no closing quote), then document the resulting error message
```
{{ inc "common|mdemo.md }}
```
* HighlightStyle styles: See https://github.com/alecthomas/chroma/tree/master/styles
* To change the root text (rem) size, 
## Priority 2: Desired but not required for the next release
* Add: copy theme feature
* Add versioning for themes (already did it in the older version). That should also mean:
  - new theme should increment version and allow an optional version, something like this: 

```
mb new theme test 0.2.0 from wide 1.1.0

```
  - Check for conflicting versions with new theme command
* Add to glossary: `command` means a command-line verb such as `new theme` or `new site`
* Ensure mdToHTML and mdFileToHTML return errors and use application error handling
* Incomplete list of things that need to be handled once I start accepting
options other than the front matter:
  - Site.Mode sets default for FrontMatter.Mode
  - Site.Language sets default for front Theme.Language
* Document all error codes
* I think I neeed to make stylesheet paths relative not absolute
* Documenting themes:
  Image properties are likely based on adjancent headers, which can be added empty (without text for the header)
* Document how Frontmatter Mode determines whether
theme-light.css or theme-dark.css is used.
* Mention in docs that one should default to post if one dones't know the difference between page and post
* Theme that's named as a number doesn't seem to work well
* I am not using the assets path
* Add a way to generaate empty
  - site config file
  - theme config file
* BUG-ish: Settle on naming convention for yaml portion of structs. Sometimes it looks
like HeaderLogo, sometimes it looks like Full-Name, sometimes it looks like Hard-wraps
* funcs.go articlefunc() doesn't work because I don't actually store the output. Should probably do it.
  don't forget that  getProjectTree() returns a list of all files on the site but discards it.
* Documente that Inhereted themes still need (empty) sidebar-left.css, sidebar-right.css, theme-light.css, theme-dark.css for the test suite only. Or maybe test suite should generate them.
* wide theme using to have {{ toc }} in the sidebar. Hve to revivi that aftermaking parser options more detailed
* The test used to show all features of a theme should include
  - HTML forms
* Append all stylesheets to a single file as encountered (for all levels of theme, so inheritance works correctly)
* pub.go: stylesheetTags()  Last stylesheet tag always gets duplicated
* RSS support
* Sitemap support
* Create a page with this intentional error (no closing quote), then document the resulting error message
```
{{ inc "common|mdemo.md }}
```
* Error in YAML file doesn't identify the YAML filename
* Change readThemeFile to readThemeConfig. Also write- version
* If nothing is avaialble for header, footer, and so on, 
publish nothing. Right now I'm publishing empty tags.
* Make most or all goldmark extensions and parsers optional
* Consider creating NewSite(filename) and rolling up site.New() into it
* Ensure Dedent is working properly. The generated source always has a newline appended. Does Dedent account for that, or should it?
* Introduce idea of drafts so you don't publish something by accident
* Move util.go to pkg/util
* createDirStructure() is no longer used, but create an example from it before deleting?
* In util.go, see if I need all the cfgPath code
* Create a test case for each error code
* Ensure each error code is documented
* Make `/docs` directory configurable
* Supplement QuitError to take only the error number for clarity
* Write tests for slice pkg
* Support for TOML front matter. See pkg/mdedxt/tomltc.go.sav and 
[Reddit RFP for TOML](https://www.reddit.com/r/golang/comments/pthh4p/paying_gig_for_foss_project_extending_the/)
* Ability to handle links with `.md` files instead of `.html`,
e.g. instead of

* I want to create an error convention that gives a clear indication of where an error happened, 
but also the original Goo-provided system messae. The convention would b esomething like
```
if err != nil {
  return(ErrCode("1033", "PREVIOUS", from, err.Error()))
}
```

instead of:

```
if err != nil {
  return(ErrCode("1023", from, err.Error()))
}
```

Reason: This would allow the reader to know precisely where the error occurred (the 1033 part, which
means that a theme config aka yaml file couldn't be opened), but also the exact nature of the
system error that generated it. Probalby that would mean a fourth parameter as shown:

The error message generated from it would be something like: "Unable to open theme file debut.yaml'. Cause: a system returned an error of "no permissions".


* If you're missing s stylesheet in the themename.yaml this section of code in publishstylesheet doesn't give enough info:
```
	err := Copy(source, dest)
	if err != nil {
		return ErrCode("PREVIOUS", err.Error(), "Publishing stylesheet")
	}
```	
The message is: Unable to copy file /Users/tom/code/m/cmd/mb/theme-test/.mb/themes/w/mw.css (error code mbz0112)
It would be nice to be able return the previous error but have a specific error for this so the online help
could point the user to this situation.  I think I may have 
a way for that. Right now I'm obscuing the previous error cause.
10/25/2021
* Work on marshalling front matter to a FrontMatter struct
* mb new foo does something, and it shouldn't

10/24/2021
* ./mb kitchen foo no longer creates a sample site. If I improve
the error handling I'll know what's happening
* My io.ReadAll() usage might be revisited for performance reasons. https://haisum.github.io/2017/09/11/golang-ioutil-readall/
* isproject and and friends don't work correctly *
* ./mb info food where food doesnt' exist isn't handling the error well
* Also target files aren't being written to foo/.mb/themes but foo/foo/.mb/themes even though the themes directory is being corectly written to foo/.mb/themes
* isProject isn't working
* I think if there's a dir starting with "." in the factory themes directory
a runtime error occurs. Fix that but also add a test for it.
10/17/2021
* Bug: `mb new site foo` calls setSiteDefaults() twice
* Bug: `mb new site /Users/tom/code/deleteme` calles setSiteDefaults() twice, the 
first time for the current directory and the second time for the target directory
* Bug: Run `mb -i` that is without build. The data structures don't get default initialization
* Themes
  - Wiring up the YAML data to internals
  - Adding theme support
    - Create a theme directory
    - convert 1 theme to yaml and put it in the directory
    - copy it to the publish directory
* Bind flags & other values to viper
* Document order of execution on startup in regards to Cobra and Viper

### Code smell

Here are some potential problems in the code.

* `md_test.go` assumes automatic header ID generation is on by default. 
Probably need to move to a more complicated test harness that dealw
with different settings for output

`[To do](todo.html)`

the markup would be:

`[To do](todo.md)`

* Investigate speed of converting byte array to string. See https://stackoverflow.com/questions/40632802/how-to-convert-byte-array-to-string-in-go . Important in code like this: 

```
fmt.Println(string(mdFileToHTML(filename)))
```
