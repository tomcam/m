# Using the Metabuzz command line 
# TODO: Unfinished

* [Create a new project](#new-site)
* [Generate HTML for a site (build)](#build)  

## New site

The syntax for creating a new site is:

```
mb new site {dir}
```

Where you'd specify a folder name for the optional
parameter `dir`, which should limit itself to letters, the 
hypen (&ndash;) charater, and numbers. 
If you omit it you'll be asked for the name.

There's no output on success and an exit code of 0 is 
returned to the operating system. On error, a
message is displayed and an exit code of 1 is
returned.

### Example 1

If your project were going to be named `foo`, you'd
do something like this:

```
mb new site foo
```

### Example 2

Omit the optional parameter. Metabuzz asks you
for the site name.

```
mb new site
Name of site to create?  foo
```

This creates the project in the local subdirectory `foo`.

### Example 3

```
mb new site /Users/rajiv/test
```

### Common end to end usage

The following example creates a site in the
local directory `foo` with stub `index.md`
file, then the `build` command generates HTML for it,
yielding the `1 file` message explainining how
many files were built. 
You can now open the file foo/index.html in a web
browser.

```
mb new site foo
mb build foo
1 file
```
This creates the project in different a directory entirely.


## Build 

Build generates HTML for the site whose directory is
named in the optional `{dir}` parameter.
If omitted, Metabuzz assumes the current directory
is the project to build.

```
mb build {dir}
```

### Example 1

