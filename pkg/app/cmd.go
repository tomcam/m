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
)

func (app *App) addCommands() {
	var (

		/*****************************************************
		  TOP LEVEL COMMAND: new
		 *****************************************************/
		cmdNew = &cobra.Command{
			Use:   "new",
			Short: "new commands: new site|theme",
			Long:  cmdBuildNewMsg,
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
					app.Site.Name = args[0]
				} else {
					// Them more likely case: it's build all by
					// itself, so go through the whole directory
					// tree and build as a complete site.
					app.Site.Name = promptString("Name of site to create?")
				}
				err := app.NewSite(app.Site.Name)
				if err != nil {
					app.QuitError(err)
				} else {
					fmt.Println("Created site ", app.Site.Name)
				}
			},
		}

		/*****************************************************
		  TOP LEVEL COMMAND: build
		 *****************************************************/

		cmdBuild = &cobra.Command{
			Use:   "build",
			Short: "build: Generates the site HTML and copies to publish directory",
			Long:  cmdBuildLongMsg,
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("Runnning app.build()")
				err := app.build()
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

	app.Cmd.AddCommand(cmdBuild)
	// After Cobra command has done its thing,
	// load configuration from config files,
	// environment, etc.
	cobra.OnInitialize(app.loadConfigs)
}
