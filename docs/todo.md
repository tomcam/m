# To do

## Priority 1: Showstoppers--required for the next release
* Seqrch for "TODO:" in source
* Make these changes when returning to the standard mb directory
  - Fix hardcoded paths in the files `gf`, `bu`
* Ensure mdToHTML and mdFileToHTML return errors and use application error handling
* Look for occurrences of App.Note(), which is only meant for prerelease usage

## Priority 2: Desired but not required for the next release
* Consider creating NewSite(filename) and rolling up site.New() into it
* Ensure Dedent is working properly. The generated source always has a newline appended. Does Dedent account for that, or should it?
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
