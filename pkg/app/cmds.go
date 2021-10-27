package app

import (
	"github.com/spf13/cobra"
	//"os"
)

func (app *App) addCommands() {
	// Initialize paths to current directory in case this
	// is something like a `mb -i` and nothing else
	var (
		pathname string
		err      error

		// Declare command to build a hardcoded test site
		//cmdKitchenSink = flag.NewFlagSet("kitchen", flag.ExitOnError)

		/*****************************************************
		  TOP LEVEL COMMAND: kitchen
		 *****************************************************/

		cmdKitchenSink = &cobra.Command{
			Use:   "kitchen",
			Short: "generates a test site",
			Long:  "creates a test site called ./kitchensink",
			Run: func(cmd *cobra.Command, args []string) {
				app.Note("kitchensink")
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
					// TODO: Use pathname in error message
					app.QuitError(ErrCode("1002", pathname))
				}
			},
		}

		/*****************************************************
		TOP LEVEL COMMAND: new
		*****************************************************/

		cmdNew = &cobra.Command{
			Use:   "new",
			Short: "create a new website or theme",
			// TODO: Flesh this out
			Long: "new site {sitename} or new theme {themename}",
			Run: func(cmd *cobra.Command, args []string) {
				if len(args) > 0 {
					pathname = args[0]
					err = app.build(args[0])
				} else {
					err = app.build("")
				}
				if err != nil {
					// TODO: Use pathname in error message
					// TODO: Generate new error code. This is used elsewhere
					app.QuitError(ErrCode("1002", pathname))
				}
			},
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
				var pathname string
				// See if the user specfied a directory name.
				if len(args) > 0 {
					pathname = args[0]
				} else {
					// Them more likely case: it's build all by
					// itself, so go through the whole directory
					// tree and build as a complete site.
					pathname = promptString("Name of site to create?")
				}
				// site.new() requires a fully qualified filename.
				if pathname == "" || pathname == "." {
					pathname = currDir()
				}
				err := app.createSite(pathname)
				if err != nil {
					app.QuitError(err)
				}
				app.Note("Created site %v", app.site.path)
				if app.Flags.Info == true {
					app.ShowInfo(pathname)
				}
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
	  GLOBAL FLAGS COULD BE CREATED HERE
		*****************************************************/

	// Example:
	// app.RootCmd.PersistentFlags().BoolVarP(&app.Flags.Verbose, "verbose", "v", false, "verbose output")

	/*****************************************************
	  AddCommand()
		*****************************************************/
	app.RootCmd.AddCommand(cmdNew)
	app.RootCmd.AddCommand(cmdKitchenSink)
	cmdNew.AddCommand(cmdNewSite)
	app.RootCmd.AddCommand(cmdInfo)
	app.RootCmd.AddCommand(cmdBuild)
}

// updateConfig() determines where configuration file (and other
// forms of configuration info, such as
// environment variables) can be found, then reads in
// all that info. It overrides defaults established
// in NewApp(). It isn't necessary. That us, NewApp()
// will have initialized the App data structure sufficiently
// to create a new project in the absence of any
// overriding config information.
func (app *App) updateConfig() {
	app.Note("\nupdateConfig()")
}
