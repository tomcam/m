# Metabuzz Theme Framework 

* Is the `background` grey?

## How code is displayed (inline code and code fences)


## Metabuzz Theme Framework FAQ

I modified my theme but the changes aren't showing up
: TODO: You probably edited .mb/themes/mytheme (or whatever theme you modified) locally and didn't copy it back to the .mb/themes directory of an older project. Remember that each time you create a new project using `mb new site`  it copies the theme from a subdirectory under the Metabuzz executable.

## See also

* [theme-framework.txt](theme-framework.txt)

Most themes packaged with Metabuzz employ The Metabuzz Theme Framework (MTF) to
style the HTML output from a a metabuzz [source file](glossary.html#source-file-1)


| Tag              | YML          |
| ---------------- | ------------ |
| `<header>`       | `Header:`    |
| `<nav>`          | `Nav:`       |
| `<article`       | `Article:`   |
| `<sidebar>`      | `Sidebar:`   |
| `<footer>`       | `Footer:`    |
| `<nav>`          | `Nav:`       |


```
Header:
  HTML: "<header>News of the Day</header>"
  File: ''
Nav:
  HTML: ''
  File: nav.md
Article:
  HTML: ''
  File: ''
Footer:
  HTML: ''
  File: footer.md
Sidebar:
  HTML: ''
  File: sidebar.md
```

