# To do

10/24/2021
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

## Priority 1: Showstoppers--required for the next release
* Seqrch for "TODO:" in source
* Document all error codes
* Make these changes when returning to the standard mb directory
  - Fix hardcoded paths in the files `gf`, `bu`
* Ensure mdToHTML and mdFileToHTML return errors and use application error handling
* Look for occurrences of App.Note(), which is only meant for prerelease usage

## Priority 2: Desired but not required for the next release
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
