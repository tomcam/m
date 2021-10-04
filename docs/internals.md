# Internals: How Metabuzz builds, starts, and runs

## REMINDERS FOR TC:
* mb alone builds in the current directory, 
or in the specified directory if a directory name is supplied 
on the command line
* mb filename.md should generate an HTML file in .mb/pub
* mb run starts a web server in the current directory 
or in the specified directory if a directory name is supplied 
on the command line
* `/cmd/mb` CLI app to compile projects in bulk
* `/cmd/web` Web app to create pages interactively

## Creating a project

* A Metabuzz website consists of at least one Markdown source
files and a site file, which is located in the the `.mb`
directory and is named `site.yaml`. 
It almost certainly includes 
other assets such as themes, graphics files,
media files, and HTML fragments.
All these together are 
called a Metabuzz *project*. One project generates
one website.


## Startup


* 


## Project structure

# TC: CONFUSING.
* Metabuzz assumes the tree of files that make up a
project are not in the root directory, but in the
/docs directory. If your site is named example.com
and it lives in a directory named `example`,
Metabuzz looks in `example/docs` for the site
file, source files, and so on.
* Metabuzz first looks for an `/.mb` subirectory
under the `/docs` directory. Without that subdirectory
Metabuzz won't generate a website.

