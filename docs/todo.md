# To do

## Priority 1: Showstoppers--required for the next release
* Not returning errors when specified stylesheets aren't found
some kind onf indicator in the theme. yaml 
* If a themefile has a defect, say somethinglike .FrontMatter.PageType (which no longer exists) you don't know exactly where the defect occurred. Or what source file was being processed at the time.
* Add idea of post and specfiications like YYYY-MM-DD or y-m-d etc, using dirs or strings as needed . That way mb new post "/blog/avengers review" would expand to something like "/blog/2022/04/21/avengers-review.html" or "/blog/2022-March-1-avengers-review.html" and so on 
* Mention in docs that should default to post
* Search for "TODO:" in source
* Documenting themes:
  Image properties are likely based on adjancent headers, which can be added empty (without text for the header)
* Document how Frontmatter Mode determines whether
theme-light.css or theme-dark.css is used.
* Incomplete list of things that need to be handled once I start accepting
options other than the front matter:
  - Site.Mode sets default for FrontMatter.Mode
  - Site.Language sets default for front Theme.Language
* Document all error codes
* Make these changes when returning to the standard mb directory
  - Fix hardcoded paths in the files `gf`, `bu`
* Ensure mdToHTML and mdFileToHTML return errors and use application error handling
* Look for occurrences of App.Note(), which is only meant for prerelease usage
* Understand whether I should create empty index.html files for dirs
that don't have anything else to oprevent directory traversla attacks
* Search and replqce almost all QuitError calls because
everything should return errors, displaying to stdout at the last
possible moment.
That'll be important for the interactive website version
* Add versioning for themes (already did it in the older version). That should also mean:
  - new theme should increment version and allow an optional version, something like this: 

```
mb new theme test 0.2.0 from wide 1.1.0

```
  - Check for conflicting versions with new theme command
* Add to glossary: `command` means a command-line verb such as `new theme` or `new site`

### Testing
* Deformed YAML file

## Credits
* https://iconduck.com/icons/84801/tools
* https://iconduck.com/icons/22486/tools

## Marketing
* Promise that in early versions while data stuctures and features may change, the
source file directory structure remains sacrosanct: a tree of Markdown files

## Priority 2: Desired but not required for the next release
* wide theme using to have {{ toc }} in the sidebar. Hve to revivi that aftermaking parser options more detailed
* Append all stylesheets to a single file as encountered (for all levels of theme, so inheritance works correctly)
* pub.go: stylesheetTags()  Last stylesheet tag always gets duplicated
* RSS support
* Sitemap support
* Error in YAML file doesn't identify the YAML filename
* Change readThemeFile to readThemeConfig. Also write- version
* If nothing is avaialble for header, footer, and so on, 
publish nothing. Right now I'm publishing empty tags.
* Make most or all goldmark extensions and parsers optional
* Consider creating NewSite(filename) and rolling up site.New() into it
* Ensure Dedent is working properly. The generated source always has a newline appended. Does Dedent account for that, or should it?
* Introduce idea of drafts so you don't publish something by accident
* Move util.go to pkg/util
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
