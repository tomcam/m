# Collections: News, blogs, and gallery sections

# TODO
* Add illustrations

Suppose you want to create a news section, blog, or gallery. 
All of these entities are termed [collections](collections.html) and have similar
characteristics. They are:

* Created once, then populated by zero or more *posts* that get
displayed using your sort criteria without you having to specify
each item explicitly.
* Designed to hold many related things. A collection differs from a page
because a page is meant to stand alone, but a collection has many
of the same kind of items in a directory structure. Examples of pages: http://metabuzz.com/about.html
or http://metabuzz.com/contact.html. An example of a collection would be the
[Metabuzz theme gallery](http://metabuzz.com/gallery/index.html).
* Sorted by attributes such as file created, last edited, or publish date
* Stored in directory structures that often reflect dates, for example,
`/blog/2023/December` or `bugs/tracker/2023/12/01`.
* Designed so that creating a new element in that collection automatically places
that item in the defined directory structure. The directory structure is a 
template that can look something like `/blog/:year/:monthnum` 
* Designed for SEO, by neatly categorizing your site in sections that are clear
and easily found.

## Collections and posts

The things stored in a collection are called *posts*. A post is a standard
Markdown file in a directory structure you specify. That directory holds
multiple pages that can be sorted and displayed automatically using
those sort criteria. If you don't specify a sort Metabuzz displays
them in descending chronological order, because that's how most
people expect news items to be presented and because search engines give
preference to new content.

## Creating a collection

Metabuzz lets you create collections using expected defaults so you don't have
to remember a lot of options. Suppose you'd like to add a blog to your site
at `example.com`. The most likely format would be this,
assuming an article was named `welcome.md`:

```
https://example.com/blog/2022/03/11/welcome.html
```

You only need to create collection once. Thereafter you add to it with the `new post` command.
That creates Markdown files and generates a directory path for them automatically.

### Creating a collection on the command line using default values

The easiest case for a collection is naming a single directory.
You could a section named `blog` (a.k.a. Metabuzz *collection*) on the command line as easily as this:

```
mb new collection blog
```

This creates a template in your site configuration file as follows:

```
blog:
  Type: collection 
  Permalink: "blog/:year/:monthnum/:day/:postname"
```

### How the permalink for a collection works

As you might imagine, `:year` gets replaced with the 4-digit year, `:monthnum` gets
replaced by the month in 2-digit form, and `:postname` is a string that
gets converted to a filename.

### Creating a collection on the command line using explicit values

Remember that the following two examples result in identical collections:

```
mb new collection blog
mb new collection blog/:year/:monthnum/:daynum/:postname
```

Let's explore other permalink options.
Perhaps you need a product update section but 
also an area for product press releases. For the news section you could do something like this:

### Creating a section for product updates using :month instead of :monthnum and :day instead of :daynum
```
mb new collection /updates/:year/:month/:day
```
This would create the directory `example.com/updates/` and
when you created a new post you'd do this:

```
mb new post updates "Announcing Metabuzz 1.1"
```

Assuming today was March 12, 2022, Metabuzz would generate the file
`https://example.com/updates/2022/december/march/saturday/12/announcing-metatbuzz-1-1.md`
because `:month` is replaced with the month name as a word, and `:day` is
replaced with the day name as a word.

### Creating a press release section with a nested directory

Your permalinks don't have to describe a directory right off the root.
This creates one at `/products/newsroom`:

```
mb new collection /products/newsroom/:year/:month/
```

* To add a press release called `metabuzz-mobile.md`:

```
mb new post products/newsroom "Metabuzz Mobile"
```
Assuming today was March 12, 2022, Metabuzz would generate the file
`https://example.com/products/newsroom/2022/december/metatbuzz-mobile.md`

## Permalink reference

This section details what you can put in a permalink and in what order.

### Permalink format

A permalink is a template for a collection's directory
and it takes this form:

```
[directory path]/{variables...}/:postname
```

* A permalink must include a directory designation. 
* A permalink always ends in `:postname` and it's implied if you don't supply it.

### Behavior when creating a collection

Here's what happens when you create a collection.

* If a permalink consists only of a directory path when you create a collection, 
Metabuzz appends `:year/:monthnum/:daynum/:postname`.

Where `:postname` is implied if not explicitly added.

### Permalink variables

Metabuzz supports the following permalink variables:

| Variable  | Purpose                                      | Sample input      | Sample output  |
| --------- | -------------------------------------------- | ----------------- | -------------- |
| :postname | Convert post name to filename with no spaces | "my article"      | my-article.md  |
| :year     | Display current year in 4 digit form         |                   | 2022           | 
| :month    | Display month name as word                   |                   | March          |
| :monthnum | Display month in 2 digit form                |                   | 03             |
| :daynum   | Display day name as word                     |                   | Friday         |
| :day      | Display date in 2 digit form                 |                   | 09             |
| :hour     | Display current hour in 2 digit form         |                   | 23             |
| :minute   | Display current minute in 2 digit form       |                   | 59             |
| :second   | Display current second in 2 digit form       |                   | 30             |
| :author   | Display name of author from front matter     | Author: Jake      | Jake           |

## Example new collection permalinks

Here are some examples of how Metabuzz handles default cases when
you specify a new collection. The full syntax for the `new collection`
example in the table below would be this:

```
new collection news
```
And as shown below, the following permalink would be generated:

#### generates permalink:
```
/news/:year/:monthnum/:daynum/:postname 
````

| new collection               | generates permalink                             | 
| ---------------------------- | ----------------------------------------------- | 
| news                         | /news/:year/:monthnum/:daynum/:postname         |
| products/news                | products/news/:year/:monthnum/:daynum/:postname |
| /products/news/:postname     | /products/news/:postname                        |
| products/news/:postname      | /products/news/:postname                        |
| news/:year/:monthnum/:daynum | /news/:year/:monthnum/:daynum/:postname         |
| articles/:postname/tech      | /articles/tech/:postnamae                       |

Suppose the current date is January 17, 2023, and the time is 11:53:22am.
Here's what some example permalinks would convert to when you added a post.


