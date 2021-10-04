# Command line operation

## TODO: Replace hardcoded names such as `site.yaml` with variable names

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





