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

* app.NewDefaultApp() allocates the App object 
* App.Cmd.Execute() gets commands from the command line
using Cobra Command (https://github.com/spf13/cobra), 
but also configuration information from config files in
the user's documents directory, config files in the
project directory, and the environment



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

## Adding to the CLI

### Adding a compound command like new site


Add this to  var( declaration under func (app *App) addCommands(). 
Note that cmdBuildNewMsg is just a string variable declared
in cmdmsgs.go, described later.

```
    /*****************************************************
      TOP LEVEL COMMAND: new
     *****************************************************/
    cmdNew = &cobra.Command{
      Use:   "new",
      Short: "new commands: new site|theme",
      Long: cmdBuildNewMsg, 
    }
 
		/*****************************************************
		    Subcommand: new site
		*****************************************************/

		cmdNewSite = &cobra.Command{
			Use:   "site {sitename}",
			Short: "new site mycoolsite",
			Long: `new site {sitename}
      Where {sitename} is a valid directory name. For example, if your site is called basiclaptop.com, you would do this:
      mb new site basiclaptop
`,
			Run: func(cmd *cobra.Command, args []string) {
				// If there are arguments after build, then
				// just convert these files one at at time.
				if len(args) > 0 {
					a.Site.Name = args[0]
				} else {
					// Them more likely case: it's build all by
					// itself, so go through the whole directory
					// tree and build as a complete site.
					a.Site.Name = promptString("Name of site to create?")
				}
				err := a.NewSite(a.Site.Name)
				if err != nil {
					a.QuitError(err)
				} else {
					fmt.Println("Created site ", a.Site.Name)
				}
			},
		}
```

At the end of the `(a *App) addCommands()` functions add these under the comment shown:

```
  /*****************************************************
    END TOP LEVEL COMMANDS BEFORE THE ABOVE )                    
   *****************************************************/    
                                                          
  app.Cmd.AddCommand(cmdNew)                           
  cmdNew.AddCommand(cmdNewSite)                              
```

### Add the variable `cmdBuildNewMsg` to `cmdmsgs.go`

The text for `cmdBuildNewMsg` is long.If it were shorter
it could be inline with the rest of the code, but long
messages get declared as variables in `cmdmsgs.go`. If
it were the only variable in `cmdmsgs.go` the while
file would look as show below, but of course the
file actually contains many variables separated by commas.

```
package app

var (
	cmdBuildNewMsg = `
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
```

