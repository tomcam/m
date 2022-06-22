# Metabuzz Command line operation
# TODO: Unfinished

* [Create a new project](#new-site)
* [Generate HTML for a site (build)](#build)  


## TODO: Replace hardcoded names such as `site.yaml` with variable names?

## new

### mb new [pathname]

Creates a new project in the directory named
in pathname. If no pathname is supplied, 
creates a project in the current directory.

Fails if there's already a project in that
directory.

#### Example: Create a project in another directory

Look for a site file in the directory `/Users/tom/html/foo`.
If it exists, exit with an appropriate message.
If it doesn't exist:
* Create a directory at `/Users/tom/html/foo`
* Generate a directory structure at `/Users/tom/html/foo/.mb`,
including themes directory, header tags directory, etc. 
* Generate a site file in `/Users/tom/html/foo/.mb/site.yaml`
* Create a minimal home page file at `/Users/tom/html/foo/index.md`
 
```
mb new /Users/tom/html/foo
```

#### Example: Create a project the current directory

```
mb new
```
* Obtain the full pathname of the current directory.
For example, suppose the current directory
is `/Users/tom/html/foo` 
* Look for a site file in the directory `/Users/tom/html/foo`.
If it exists, exit with an appropriate message.
If it doesn't exist, follow the procedure listed in [Example: Create a project in another directory](cli.html#example-create-a-project-in-another-directory)





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

