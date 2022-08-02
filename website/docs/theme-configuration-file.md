# Theme configuration file

The Metabuzz theme configuration file tells Metabuzz how to employ
the CSS and HTML files that make up a theme's appearance and
some of its behavior. For example, if a theme doesn't support sidebars,
add `Sidebar: false` under `Supports` and no sidebars will be
generated.


Here's a sample Metabuzz theme configuration file, and below
you'll see how its contents are used.


```yaml
Branding: 'Metabuzz Pillar'
Description: 'A flagship Metabuzz theme'
Supports:
  MTF: true
  Mode: true
  Header: true
  Nav: true
  Sidebar: true
  Footer: true
Stylesheets:
  - reset.css
  - fonts.css
  - bind.css
  - sizes.css
  - theme-light.css
  - layout.css
  - pillar.css
  - responsive.css
Version: 1.2.3
Nav:
  HTML: ''
  File: nav.md
Header:
  HTML: ''
  File: header.md
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

## Branding

Metabuzz uses directory names internally when theme names are employed,
because directory names must be unique. `Branding` is an optional paragraph of text to give your theme that looks a little better. If your theme's directory were named `metabuzz-pillar` you could use `Branding` to ensure it gets
displayed as `Metabuzz Pillar` where possible.

### Example
```
Branding: 'Metabuzz Pillar'
```
## Description

When your theme gets displayed in a gallery format Metabuzz uses
the `Description` to display a short sentence or paragraph 
highlighting the theme.

### Example
```
Description: 'A flagship Metabuzz theme'
```

## Supports

Tells what elements of a [Metabuzz Theme Framework](theme-framework.html) (MTF) theme are supported, and
supplies some information about the theme itself. For example,
if a theme doesn't support footers you would add `Footer: false`.

## Supports: mode
Tells whether a [Metabuzz Theme Framework](theme-framework.html) (MTF) theme has support for dark and light themes.

### Possible values: true, false

`Mode: true` means this theme has support for both dark and light themes.
If it supports both, then `theme-light.css` will be used when `Mode: light`
is specified in the [front matter](front-matter.html), and
`theme-dark.css` will be used when `Mode: dark` is specified.

### Example
```
Supports:
  Mode: true
```


## Supports: Header 
Tells whether a [Metabuzz Theme Framework](theme-framework.html) (MTF) theme has support for .


Supports:
  Header: true
  Nav: true
  Sidebar: true
  Footer: true
Stylesheets:
  - reset.css
  - fonts.css
  - bind.css
  - sizes.css
  - theme-light.css
  - layout.css
  - pillar.css
  - responsive.css
Version: 1.2.3
Nav:
  HTML: ''
  File: nav.md
Header:
  HTML: ''
  File: header.md
Article:
  HTML: ''
  File: ''
Footer:
  HTML: ''
  File: footer.md
Sidebar:
  HTML: ''
  File: sidebar.md

