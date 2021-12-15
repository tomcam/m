# About Metabuzz

Metabuzz is a *static site generator*. It lets you create professional, compatible, safe, high-performance websites using simple text files using Markdown, an easy-to-learn typing convention. Metabuzz is free to use, and because Markdown is an [industry standard](https://commonmark.org) if you end up choosing another static site generator all of the content you write can be used unchanged.

## Metabuzz design goals

Here's where Metabuzz fits into the website publishing landscape. The Metabuzz design philosophy can be stated like this:

* You should be able to install Metabuzz and create a complete website within 60 seconds of installation.
* You should be able to build an informational website immediately, knowing only Markdown. You don't even have to know Markdown because there's a full tutorial. In other words, you don't also need to learn a specialized template language or front matter conventions just to put your first site up.
* You should be able to create Metabuzz themes without knowing another computer language like Ruby, a CSS preprocessor, or a specialized template language. (Metabuzz does use the sample template language as Go and it's very useful, but you can create a theme without knowing it.)
* If you live in a developing country or rural area where bandwidth is limited you should be able to generate attractive, professional-loking sites that use minimal resources even if you don't have design skills or know CSS. Most Metabuzz themes are exceptionally small and are tested with browsers going back to Android version 4. Conversely if you don't care about page size you can easily create more resource-intensive themes.
* Metabuzz has no external dependencies. That means it doesn't require that you install anything else to create your site, and you don't need to learn about CSS preprocessors or install any Node.js packages.
* Metabuzz comes with several production-ready themes, and they are batteries-included: with a single line of text your search engine-optimized site will be built with a sample header, navigation bar, sidebar, article area, and footer.
* You should be able to test the vast majority of your site locally, without having to spin up a local webserver. Just open the HTML files with your browser. (Some themes may have features that require a server, but all sites are at least viewable without one.)
* The generated site should be GDPR-compliant (following European privacy standards) and ADA-compliant (offering maximum support for disabled users)
* If you decide to use another CMS all of your Markdown content should be usable without any changes.
* Metabuzz takes a "words first" approach in that sites are semantically driven using modern web standards. Headers 1-3 are used purely for subject matter because that's the deepest nesting level humans and websites understand easily. Header, footer, and nav bar are used for their intended purposes with the right tags. Tables are seldom used for formatting purposes. Your text is automatically enclosed in an `<article>` tag to give it the appropriate importance to search engines, and so on.
* The site creation experience is meant to resemble the early days of web page creation, but a little better.  It's similar in that you can just create a page and publish it simply by copying files to the publication directory (sometimes called the WWW directory). What makes it more fun is that while you can use HTML if you wish, you can stop worrying about hard-to-read source text and having to keep track of angle brackets. This requires some magic behind the scenes because Metabuzz sites use themes, but it is normally a seamless experience.

## What makes Metabuzz different from its competition?

**The shortest learning curve** When you actually try building a site with other static site generators, you invariably discover that it's not just a matter of knowing Markdown and being able to get going immediately. Some require to build your own theme from scratch before you can start writing. Some require that you know Ruby. All the others have more complicated requirements to create a theme, whereas you can create a full Metabuzz theme using only Markdown. Most require that you precede your Markdown with specialized, non-Markdown front matter before you can build your site. Metabuzz has front matter too, but it's all optional. Your first site can be published in seconds using only Markdown--or even just plain text with no Markdown at all.

**Unmatched feature set** Metabuzz builds complete websites with a much richer feature set than most static site generators:

* Most static site generators use a simplified page structure that's pretty much just an article or group of articles. Metabuzz thinks of a page as having a header, navigation bar, main article area, sidebar, and footer. Its themes do not require all these elements but they are supported fully if your site needs them
* Step by step tutorials for creating a simple site, creating a production site, and designing your own themes
* Themes can also have unlimited page types which inherit from the main theme. They act as "subthemes" 
* The built-in Metabuzz themes have support for light and dark themes
* The Metabuzz theme architecture is highly search engine optimized (SEO features are built in)

**Simple GitHub and Netlify integration** You can get free hosting of static websites through GitHub and Netlify. Metabuzz supports both fully, which means you can get an information site up for free, paying only for a domain name if you feel it's necessary.

## Metabuzz pros

Here's what makes Metabuzz your best choice for creating static websites:

* Metabuzz requires only knowledge of the command line and a vague idea of what Markdown is used for. If you can edit a text file, you can create a commercial-quality informational site without ever having to learn template languages, or CSS. However:
* When you choose to go deeper, Metabuzz rewards knowledge of existing skills. Its style sheets are standard CSS. You can design themes using Markdown or HTML. You can use a CSS preprocessor if you want to. 
* You don't need to learn Ruby to create a theme. You don't need to install a preprocessor of NodeJS. 

## Metabuzz cons

Here's why you might not want to use Metabuzz
