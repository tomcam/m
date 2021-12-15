# Metabuzz documentation conventions

This section is mostly for internal use, but if you want
to provide complete documentation for a Metabuzz theme it 
may save you from too much experimenting. It can get a 
little tricky when you're writing documentation about Markdown
and especially templates, from within a Markdown document.

It also includes information on tooling required to preprocess
documents for this site, take screenshots, illustrating themes,
and so on.

See also [Documentation conventions](documentation-conventions.html)


## Conventions

Some of these issues still need resolving.
* A note or warning should use an H3:

```
### Note
Directories with names preceded by a ". are excluded
```
  
* File listings should be preceded with an H4 with the filename:

```
#### file: index.md
```

## Creating directory tree listings on MacOS

The Metabuzz directory tree listings are created using the [tree](https://rschu.me/list-a-directory-with-tree-command-on-mac-os-x-3b2d4c4a4827) command-line utility by Robin Schulz.

# Internal documentation notes

Most Metabuzz documentation is standard Markdown. Documenting the [template language](glossary.html#template-language) complicates things.


# Using s.sh to format for code fences


# Markdown notes

## Look at Syntax.text for Markdown examples

Good source of Markdown info at [Syntax.text](https://github.com/russross/blackfriday/blob/master/testdata/Markdown%20Documentation%20-%20Syntax.text)

## Adding anchors the easy way

You can add a custom anchor tag like this:

```
<a id="gallery"></a>
```

## Anchor ID tags the hard way

You can also do this:

```
<h2 id="block">Block Elements</h2>
```

## End a line with two spaces for a `<br/>` effect

## Links

### Titles in links

They appear as tooltips:

```
This is [an example](http://example.com/ "Title") inline link.

[This link](http://example.net/) has no title attribute.
```

Result:

This is [an example](http://example.com/ "Title") inline link.

[This link](http://example.net/) has no title attribute.

### Reference-style links

```
See my [About](/about/) page for details.
```

From https://github.com/russross/blackfriday/blob/master/testdata/Markdown%20Documentation%20-%20Syntax.text 

Reference-style links use a second set of square brackets, inside
which you place a label of your choosing to identify the link:

```
This is [an example][id] reference-style link.
```

You can optionally use a space to separate the sets of brackets:

```
This is [an example] [id] reference-style link.
```

## Notes on escaping things

To include a literal backtick character within a code span, you can use
multiple backticks as the opening and closing delimiters:

```
    ``There is a literal backtick (`) here.``
```


```
It's worth noting that it's possible 
to trigger an ordered list by accident, 
by writing something like this:

    1986. What a great season.

In other words, a *number-period-space* 
sequence at the beginning of a line. 
To avoid this, you can 
backslash-escape the period:

    1986\. What a great season.
```

### Filenames

## Screenshots

### Screenshot dimensions

* Landscape mode screenshot dimensions: 1280x1024
* Portrait mode screenshot dimensions: 1024x1280
* Smaller versions for galleries are currently 200x250

### Screenshot file naming conventions

* Append "-" and the dimensions after the screenshot.
* For galleries, use "theme-", theme name, "-", and dimensions. Examples:
```
theme-wide-dark-1024x1280.png
theme-wide-dark-200x250.png
theme-wide-dark-nosidebar-1280x1024.png
```

