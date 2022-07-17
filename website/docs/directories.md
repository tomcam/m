# Directories

Your project lives in a directory that must meet a few
criteria. 

Suppose the directory is `/Users/myproject`. The project's source 
Markdown files are expected to appear in `/docs`,
e.g. `/Users/myproject/docs`. The supporting files
for configuration, themes, site information, and so
on live in the `/.mb` directory (e.g. `/Users/myproject/.mb`). 

The reason for this is that it's natural for Metabuzz to be used
for software projects, and leaving the Markdown source files
in the root directory clutters it, making the software
source files hard to find.

Here are some examples of what the directories look like.

## /myproject

Suppose you create a project directory like this, where you'd
replace `/Users/myproject` with the full path and directory
name of your own project:

```
mb new site /Users/myproject
cd /Users/myproject
```
### Note: What if you're adding a Metabuzz project to an existing directory?

It's just as easy to create a Metabuzz project in an existing directory.
Change to that directory and run `mb new site .`, as shown:

```
cd /Users/existingproject
mb new site .
```

### First directory level

When you create a new project you'll see a directory tree 
looking like this at the first level, 
where the trailing slash in `docs/` and `.mb/` indicate directories:


```
/Users/myproject
 ├── index.md
 ├── docs/
 └── .mb/
```

If you're experienced with other static site generations 
and prefer to use `README.md` instead of `index.md` you can. Just replace
the `index.md` file with `README.md`.)

## /docs 

All your Markdown source files, image files, and other assets you want
to be published go in `/docs`, which starts out empty.

Here's a simple example of a directory. Source files with Markdown
extensions such as `md`, `mdown`, and so on get converted to HTML, 
then copied to the publish directory. All other files, such as
existing HTML, PNGs, JPGs, etc., get copied over unchanged.

```
docs
 ├── Section 1.md
 ├── Chapter 1.md
 ├── Chapter 2.md
 ├── illustration1.png
 ├── Section 2.md
 ├── /Section 3
 │    ├── index.md
 │    ├── Chapter 3.md
 │    ├── illustration1.png
 │    ├── illustration9.png
 │    └── Subchapter.md
 └─── Appendix.html
```

## /.mb

```
.mb
 ├── /themes
 │    ├─── /aventurine
 │    ├─── /debut 
 │    │     ├─── /blog
 │    │     └─── /news
 │    ├── /wide
 │    │    └─── /blog
 │    └── /debut 
 │         └─── news
 ├── /common
 ├── /headers
 ├── /pub
 │    ├── Section 1.html
 │    ├── Chapter 1.html
 │    ├── Chapter 2.html
 │    ├── illustration1.png
 │    ├── Section 2.md
 │    ├── /Section 3
 │    │    ├── index.html
 │    │    ├── Chapter 3.html
 │    │    ├── illustration1.png
 │    │    ├── illustration9.png
 │    │    └── Subchapter.html
 │    └─── Appendix.html
 ├── /scodes
 └─── site.yaml

```
## .mb

The `/.mb` directory contains project configuration information,
such as all themes, the actual published HTML files, and boilerplate
files. 

## /docs

The `/docs` directory contains the source Markdown files that
actually comprise the text of the HTML or documentation site you create. 
They can be nested into subdirectories.


