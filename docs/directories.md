# Directories

Your project lives in a directory that must meet a few
criteria. 

Suppose the directory is `/Users/tom/html/foo`. The project's source 
Markdown files are expencted to appear in `/docs` 
(e.g. `/Users/tom/html/foo/docs`. The supporting files
for configuration, themes, site information, and so
on live in the `/.mb` directory (e.g. `/Users/tom/html/foo/.mb`). 

The reason for this is that it's natural for Metabuzz to be used
for software projects, and leaving the Markdown source files
in the root directory clutters it, making the software
source files hard to find.

Here are some examples of what the directories look like.

## /docs 

```
docs
 ├── Section 1.md
 │    ├── Chapter 1.md
 │    └── Chapter 2.md
 ├── Section 2.md
 ├── Section 3.md
 │    └── Chapter 3.md
 │        └── Subchapter.md
 └─── Appendix.html
```

## /.mb

```
.mb
 ├── allthemes
 │   ├─── aventurine
 │   ├─── debut 
 │   │    ├─── blog
 │   │    └─── news
 │   ├── wide
 │   │   └─── blog
 │   └── debut 
 │       ├─── blog
 │       └─── news
 ├── common
 ├── headers
 ├── pub
 │    ├── Section 1.html
 │    │    ├── Chapter 1.html
 │    │    └── Chapter 2.html
 │    ├── Section 2.html
 │    ├── Section 3.html 
 │    │    └── Chapter 3.html
 │    │        └── Subchapter.html
 │    ├─── Glossary.html
 │    │ 
 │    └── themes
 │        └── wide
 │            └─── blog
 ├── scodes
 ├── site
 └── themes
      └── wide 
          └─── blog
```
## .mb

The `/.mb` directory contains project configuration information,
especially the ::

## /docs

The `/docs` directory contains the source Markdown files. They
can be nested into subdirectories.



