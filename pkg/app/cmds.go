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
		  TOP LEVEL COMMAND:build
		 *****************************************************/
		cmdBuild = &cobra.Command{
			Use:   "build",
			Short: "build: Generates the site HTML and copies to publish directory",
			Long:  cmdBuildLongMsg,
			/*
							Long: `"build: Generates the site HTML and copies to publish directory
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
				`,
			*/
			Run: func(cmd *cobra.Command, args []string) {
				err := app.build()
				if err != nil {
					app.QuitError(err)
				}
			},
		}
	)
	app.Cmd.AddCommand(cmdBuild)

	// After Cobra command has done its thing,
	// load configuration from config files,
	// environment, etc.
	cobra.OnInitialize(app.loadConfigs)
}
