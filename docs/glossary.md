## cfg value or config value

A config value is one that can come from any of several
places. They are obtained using functions that
start with `App.Cfg`. For example, the theme name might normally
come from the individual `Page.FrontMatter` setting. 
Or you might prefer to use `Site.FrontMatter` to 
set a default theme for the entire site, then change it
only for specific pages in `Page.FrontMatter.`
Or it could come from a config file in the user's
application data directory, or even the environment.

## project

The set of files needed to generate a website.
That includes one or more Markdown source
files, along with other assets such as themes, graphics,
media files, schema.org files, and HTML fragments. 
One project generates
one website. See also [Creating a project](internals.html#creating-a-project)


