# Internals: How Metabuzz builds, starts, and runs

## Utilities
* errdoc 1234 creates a help file for error code 1234. Must be 4 digits.
* mbtest deletes and recreates the directory .theme-test, and populates
it with automatically generated pages to test themes visually
* mbtodo brings up the todo list

## REMINDERS FOR TC:
* mb filename.md should generate an HTML file in .mb/pub
* mb run should start a web server in the current directory 
or in the specified directory if a directory name is supplied 
on the command line
* `/cmd/mb` CLI app to compile projects in bulk
* `/cmd/web` Web app to create pages interactively

## Testing
* Ideas for defective themename.yaml files
```
Branding: Debut theme"
Stylesheets:
fred.css
```


## Code reminders
* Base name of theme is in app.Page.Theme.Name
* Publish means copying from somewhere such as the source Markdown
file or a theme in .mb/themes to the generated project directory,
where everything is expected to have read permissions to the world
upon publication


# Theme
* Site.FactoryThemesPath is where the source themes are
* loadThemeConfig() forces Page.FrontMatter.Theme to lowercase
because it will be the name of a file. Could be like "debut" or
also "debut/gallery/item". DOCUMENT that the covnention is always
forward slash even on Windows
Put CSS files sin this order
- reset.css
- fonts.css
- bind.css
- sizes.css
- theme-light.css
- layout.css
- [theme name such as wide].css
- sidebar-left.css (or sidebar-right.css)
- responsive.css

REMEMBER responsive*.css has to follow sidebar*.css

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

* Some cobra stuff:
  - https://stackoverflow.com/questions/43847791/why-is-cobra-not-reading-my-configuration-file

Here's what happens in `main.main`:

```
func main() {
	app := app.NewApp()
	app.Execute()
}
```
* In the `main` function, NewApp() allocates the App object, then calls its
`App.Execute()` method.
* To summarize, `App.Execute()`, obtains commands
and flags from the command line, obtains other configuration from 
variables files and the environment, and then calls code based on the results.

In more detail:
* `app.initCobra()` is used the same way they use  the `init`
function of any package (see https://github.com/spf13/cobra/blob/master/user_guide.md). It's used to create flags, like this:

```
app.RootCmd.PersistentFlags().StringVar(&app.cfgFile, "config", "", "config file (default is $HOME/.mb.yaml)")
```

## Different way of looking at execution flow
* `cobra.OnInitialize(app.initConfig)`
* `App.initConfig()` # Same as init() in Cobra
* App.Execute() 
* `app.RootCmd.Execute()` is the code that 
actually gets command line arguments and flags, then executes whatever
code they call. 
* Before that `app.RootCmd.Execute()` calls `app.initCobra()` 
*	At the end of `App.initCobra()`,  `cobra.OnInitialize(app.initConfig)` 
gets called. The argument to cobra.OnInitialize() is optional, 
but Metabuzz calls `App.initConfig()`.
* `App.initConfig()` is where Viper starts looking for configuration files,
and where the Application object can finish initialization (because
some of its state, such as which Goldmark extensions to apply), 
depends on config values set on the command line, in the environment,
or whatever.
After it runs `app.RootCmd.Execute()` finally runs.
* `app.RootCmd.Execute()`, actually does the command line parsing. First it calls initCobra(), which is explained just after this. It also uses Viper to obtain configuration information from config files in
the user's documents directory, config files in the
project directory, and the environment.
* App.Execute() returns its value to cobra.CheckErr(), which probably
does something useful but I don't know what that is



## Project structure

* Metabuzz assumes the tree of files that make up a
project are not in the root directory, but in the
/docs directory. If your site is named example.com
and it lives in a directory named `example`,
Metabuzz looks in `example/docs` for the site
file, source files, and so on.
* Metabuzz first looks for an `/.mb` subirectory
under the `/docs` directory. Without that subdirectory
Metabuzz won't generate a website. Otherwise Metabuzz would create an .mb directory and all its subdirectories if you accidentally ran `mb build`
in a directory that wasn't supposed to be your website.


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
## Reading configuration (or cfg) after startup

* Front matter: call `app.Page.frontMatterMust()` or
`app.Page.frontMatterMustLower()` (forces return
value to loser case),
e.g., `app.Page.frontMatterMust("Description")`
to get some value without throwing an error if
for example that key doesn't exist.

