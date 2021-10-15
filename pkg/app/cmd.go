// ADD THIS
// Directory this site uses to copy themes from. If the -d option was
// set, use the global factory themes directory. Otherwise, use local copy

package app

import (
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	//"github.com/tomcam/m/pkg/defaults"
	//"github.com/spf13/viper"
)

var (

	// Declare command-line subcommand to build a project
	cmdBuild = flag.NewFlagSet("build", flag.ExitOnError)

	// Declare command-line subcommand to display config info
	//cmdInfo = flag.NewFlagSet("info", flag.ExitOnError)
)

func (app *App) addCommands() {
	var (

		/*****************************************************
		  TOP LEVEL COMMAND: new
		 *****************************************************/
		cmdNew = &cobra.Command{
			Use:   "new",
			Short: "new commands: new site|theme",
			Long:  cmdNewMsg,
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
					app.site.name = args[0]
				} else {
					// Them more likely case: it's build all by
					// itself, so go through the whole directory
					// tree and build as a complete site.
					app.site.name = promptString("Name of site to create?")
				}

				// Allocate a Site object
				var err error
				app.site, err = app.site.New()
				if err != nil {
					app.QuitError(err)
				}
				// Initialize the Site object
				err = app.site.NewSite()
				if err != nil {
					app.QuitError(err)
				} else {
					fmt.Println("Created site ", app.site.name)
				}
			},
		}

		/*****************************************************
		    Subcommand: new theme 
		*****************************************************/

		cmdNewTheme = &cobra.Command{
			Use:   "theme {themename}",
			Short: "new theme mytheme",
			Long: `new theme {themename}
      Where {themename} is a valid directory name. For example, if your site is called Autumn Leaves, you would do this:
      mb new site autumn-leaves 
`,
			Run: func(cmd *cobra.Command, args []string) {
        var themeName string
				if len(args) > 0 {
					themeName = args[0]
				} else {
					// Them more likely case: it's build all by
					// itself, so go through the whole directory
					// tree and build as a complete site.
					themeName = promptString("Name of theme to create?")
				}
        app.Note("Create the theme %v\n", themeName)
			},
		}

		/*****************************************************
		  TOP LEVEL COMMAND: build
		 *****************************************************/

		cmdBuild = &cobra.Command{
			Use:   "build",
			Short: "generates the  and copies to publish directory",
			Long:  cmdBuildLongMsg,
			Run: func(cmd *cobra.Command, args []string) {
				var err error
				if len(args) > 0 {
					err = app.build(args[0])
				} else {
					err = app.build("")
				}
				if err != nil {
					app.QuitError(err)
				}
			},
		}
	)
	/*****************************************************
	  END TOP LEVEL COMMANDS BEFORE THE ABOVE )
	 *****************************************************/

	app.Cmd.AddCommand(cmdNew)
	cmdNew.AddCommand(cmdNewSite)
	cmdNew.AddCommand(cmdNewTheme)

	app.Cmd.AddCommand(cmdBuild)



  // Handle global flags such as Info 
	app.Cmd.PersistentFlags().BoolVarP(&app.flags.Info, "info", "i", false, "show debug info")

	// After Cobra command has done its thing,
	// load configuration from config files,
	// environment, etc.
	cobra.OnInitialize(app.loadConfigs)
}


