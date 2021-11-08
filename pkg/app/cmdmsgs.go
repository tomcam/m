package app

var (
	cmdBuildLongMsg = `
			build: Generates the site HTML and copies to publish directory 

      Typical usage:
      : Create the project named mysite in its own directory.
      : (Generates a tiny file named index.md)
      mb new site mysite
      : Make that the current directory. 
      cd mysite
      : Optional step: Write your Markdown here!
      : Find all .md files and convert to HTML
      : Copy them into the publish directory named .pub
      mb build
      : Load the site's home page into a browser.
      : Windows users, omit the open
      open .pub/index.html
`

	CmdNewMsg = `
site: Use new site to start a new project. Use new theme to 
create theme based on an existing one. 
      Typical usage of new site:
      : Create the project named mysite in its own directory.
      : (Generates a tiny file named index.md)
      mb new site mysite
      : Make that the current directory. 
      cd mysite
      : Optional step: Write your Markdown here!
      : Find all .md files and convert to HTML
      : Copy them into the publish directory named .pub
      mb build
      : Load the site's home page into a browser.
      : Windows users, omit the open
      open .pub/index.html
`
)
