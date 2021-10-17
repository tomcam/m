package app

import (
	"github.com/spf13/cobra"
)

func (app *App) addCommands() {
	var (
		pathname string
		err      error
		/*****************************************************
		TOP LEVEL COMMAND: info
		*****************************************************/

		CmdInfo = &cobra.Command{
			Use:   "info",
			Short: "short about info",
			Long:  `long about info`,
			Run: func(cmd *cobra.Command, args []string) {
				app.info()
			},
		}

		/*****************************************************
		TOP LEVEL COMMAND: build
		*****************************************************/

		CmdBuild = &cobra.Command{
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
					// TODO: Alter to include pathname?
					app.QuitError(ErrCode("1002", pathname))
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
	  GLOBAL FLAGS
		*****************************************************/
	app.RootCmd.PersistentFlags().BoolVarP(&app.flags.QTest, "q", "q", false, "something q")
	app.RootCmd.PersistentFlags().BoolVarP(&app.flags.RTest, "r", "r", false, "something r")

	/*****************************************************
	  AddCommand()
		*****************************************************/
	app.RootCmd.AddCommand(CmdInfo)
	app.RootCmd.AddCommand(CmdBuild)
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
	app.Note("updateConfig()")
}
