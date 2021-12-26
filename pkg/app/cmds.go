package app

import (
	"github.com/spf13/cobra"
	"github.com/tomcam/m/pkg/default"
	"path/filepath"
)

func (app *App) addCommands() {
	// Initialize paths to current directory in case this
	// is something like a `mb -i` and nothing else
	var (
		pathname string
		err      error

		/*****************************************************
		  TOP LEVEL COMMAND: interview
		 *****************************************************/

		cmdInterview = &cobra.Command{
			Use:   "interview",
			Short: "Fills in questions about your site",
			Long:  "Allows you to fill get links and directory names so you don't need to edit a configuration file",
			Run: func(cmd *cobra.Command, args []string) {
				var err error
				if len(args) > 0 {
					err = app.interviewSiteBrief()
				} else {
					err = app.interviewSiteBrief()
				}
				// TODO: Return standard error?
				if err != nil {
					app.QuitError(err)
				}
			},
		}

		/*****************************************************
		  TOP LEVEL COMMAND: kitchen
		 *****************************************************/

		cmdKitchenSink = &cobra.Command{
			Use:   "kitchen",
			Short: "generates a test site",
			Long:  "creates a test site called ./kitchensink",
			Run: func(cmd *cobra.Command, args []string) {
				var err error
				if len(args) > 0 {
					err = app.kitchenSink(args[0])
				} else {
					err = app.kitchenSink("")
				}
				if err != nil {
					app.QuitError(err)
				}
			},
		}

		/*****************************************************
		TOP LEVEL COMMAND: info
		*****************************************************/

		cmdInfo = &cobra.Command{
			Use:   "info",
			Short: "Get information about this project",
			Long:  "Get information about this project",
			Run: func(cmd *cobra.Command, args []string) {
				pathname = ""
				if len(args) > 0 {
					pathname = args[0]
				}
				err = app.ShowInfo(pathname)
				if err != nil {
					// TODO: Use pathname in error message
					app.QuitError(ErrCode("PREVIOUS", pathname))
				}
			},
		}

		/*****************************************************
		TOP LEVEL COMMAND: build
		*****************************************************/

		cmdBuild = &cobra.Command{
			Use:   "build",
			Short: "Generate a website from your Markdown files",
			Long:  "Generate a website from your Markdown files",
			Run: func(cmd *cobra.Command, args []string) {
				if len(args) > 0 {
					pathname = args[0]
					err = app.build(args[0])
				} else {
					err = app.build("")
				}
				if err != nil {
					app.QuitError(ErrCode("0923", err.Error()))
				}
			},
		}

		/*****************************************************
		TOP LEVEL COMMAND: new
		*****************************************************/
		// TODO: Do/correct all the help text on all of these
		CmdNew = &cobra.Command{
			Use:   "new",
			Short: "new commands: new site|theme|page",
			Long: `site: Use new site to start a new project. Use new theme to 
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
`,
		}

		/*****************************************************
		    Subcommand: new site
		*****************************************************/

		CmdNewSite = &cobra.Command{
			Use:   "site {sitename}",
			Short: "new site mycoolsite",
			Long: `new site {sitename}
      Where {sitename} is a valid directory name. For example, if your site is called basiclaptop.com, you would do this:
      mb new site basiclaptop
`,
			Run: func(cmd *cobra.Command, args []string) {
				var pathname string
				var err error
				// See if the user specfied a directory name.
				if len(args) > 0 {
					pathname = args[0]
				} else {
					// Them more likely case: it's build all by
					// itself, so go through the whole directory
					// tree and build as a complete site.
					pathname = promptString("Name of site to create?")
				}
				// Site.new() requires a fully qualified filename.
				/*
					if pathname == "" || pathname == "." {
						pathname = currDir()
					}
				*/
				if pathname, err = filepath.Abs(pathname); err != nil {
					app.QuitError(ErrCode("PREVIOUS", err.Error()))
				}
				if err = app.newSite(pathname); err != nil {
					//app.QuitError(ErrCode("0924", pathname))
					//app.QuitError(ErrCode("0924", "PREVIOUS", err.Error()))
					app.QuitError(ErrCode("PREVIOUS", err.Error()))
				}
				app.Debug("Created site %v", app.Site.path)
				if app.Flags.Info == true {
					app.ShowInfo(pathname)
				}
			},
		}

		/*****************************************************
		    Subcommand: new theme
		*****************************************************/

		CmdNewTheme = &cobra.Command{
			//Use:   "new theme {from} {to}",
			Use:   "theme",
			Short: "theme",
			// TODO: Copy text from old site
			Long: `new theme {from} {to}
      Creates a new theme based on an existing one.
      {to} must be a valid directory name. For example, 
      if you want create a theme called 
      itemdetail from the theme named item:

      mb new theme item-detail item
`,
			Run: func(cmd *cobra.Command, args []string) {
				if !isProject(".") {
					app.QuitError(ErrCode("1025", currDir()))
				}
				var from, to string
				// See if the user specfied a theme name.
				switch len(args) {
				case 0:
					from = promptString("Theme to copy from?")
					to = promptString("Name of new theme?")
				case 1:
					from = args[0]
					to = promptString("Name of theme to create from theme " +
						from + "?")
				case 2:
					from = args[0]
					to = args[1]
				}

				// TODO: Difference between this and copyTheme?
				//err := app.newTheme(from, to, app.Flags.Factory)
				err := app.newTheme(from, to)
				//err := app.copyTheme(from, to)
				if err != nil {
					app.QuitError(err)
				}
				app.Print("Created theme %v", to)
			},
		}

		/*****************************************************
		    Subcommand: new page
		*****************************************************/

		CmdNewPost = &cobra.Command{
			Use:   "post {postname}",
			Short: "post",
			Long: `new post "{collection}" "{postname}"
      Where both {postname} and {collection} are either a single word, or a quoted string value.
      mb new post "
`,
			Run: func(cmd *cobra.Command, args []string) {
				var collection, postname string

				switch len(args) {
				case 0:
					collection = promptString("Directory (collection) to post to?")
					postname = promptString("Title of post?")
				case 1:
					collection = args[0]
					postname = promptString("Title of post for %v?", collection)
				case 2:
					collection = args[0]
					postname = args[1]
				}

				err := app.newPost(collection, postname)
				if err != nil {
					//app.QuitError(ErrCode("0928", postname))
					app.QuitError(ErrCode("PREVIOUS", err.Error()))
				}
				app.Debug("Created post %v in %v", postname, collection)
			},
		}

		/*****************************************************
		END TOP LEVEL COMMANDS
		*****************************************************/

	)

	/*****************************************************
	THE PREVIOUS ) IS THE END OF TOP LEVEL COMMANDS
	*****************************************************/

	/*****************************************************
	  GLOBAL FLAGS CREATED HERE
		*****************************************************/

	// See also flgs.go
	CmdNewTheme.PersistentFlags().BoolVarP(&app.Flags.Factory, "factory", "", false, "use factory theme, not from local project")
	CmdNewSite.PersistentFlags().StringVar(&app.Flags.Starters, "starter", "", "config file (default is "+defaults.ConfigStartersFilename+")")
	CmdNewSite.PersistentFlags().StringVar(&app.Flags.Site, "site", "", "site config file (default is "+defaults.SiteConfigFilename+")")

	/*****************************************************
	  AddCommand()
		*****************************************************/
	app.RootCmd.AddCommand(CmdNew)
	app.RootCmd.AddCommand(cmdKitchenSink)
	app.RootCmd.AddCommand(cmdInterview)
	CmdNew.AddCommand(CmdNewSite)
	CmdNew.AddCommand(CmdNewTheme)
	CmdNew.AddCommand(CmdNewPost)
	app.RootCmd.AddCommand(cmdInfo)
	app.RootCmd.AddCommand(cmdBuild)
}
