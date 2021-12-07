package app

import (
	"fmt"
	/*
		"bufio"
		"os"
		"strings"
	*/)

///func promptString(prompt string) string {
//func // promptYes() displays a prompt, then awaits
//func promptYes(prompt string) bool {
//promptStringDefault(prompt string, defaultValue string) string {
// See also inputString(), promptYes()
//func inputString() string {

func (app *App) interviewSiteBrief() error {
  app.Print("interviewSiteBrief()")
	path := currDir()
	if err := app.changeWorkingDir(currDir()); err != nil {
		app.Debug("\tUnable to change to directory (%v)", path)
		return ErrCode("1109", path)
	}
	if err := app.readSiteConfig(); err != nil {
		return ErrCode("PREVIOUS", err.Error())
	}
  site := app.Site
	fmt.Println("Let's get a few pieces information about your site and fill it in for you. You can run this interview as often as you wish. This is the brief interview. You can also run the full site interview to enter even more information.")
	site.Company.Name = promptStringDefault("Company name?", site.Company.Name)
	site.Theme = promptStringDefault("Default theme for new pages?", site.Theme)
	site.Author.FullName = promptStringDefault("Author?", site.Author.FullName)
	site.Branding = promptStringDefault("Branding?", site.Branding)
	site.Social.Twitter = promptStringDefault("Twitter URL?", site.Social.Twitter)
	site.Social.Facebook = promptStringDefault("Facebook URL?", site.Social.Facebook)
	site.Social.YouTube = promptStringDefault("YouTube channel URL?", site.Social.YouTube)
	site.Social.Instagram = promptStringDefault("Instagram URL?", site.Social.Instagram)
	if promptYes("Keep these changes to your site configuration?") {
		app.Site = site
    app.Print("Your new site info includes:\n\tcompany name: %s\n\tdefault theme: %s\n\tauthor: %s\n", app.Site.Company.Name, app.Site.Theme, app.Site.Author.FullName)
		if err := app.writeSiteConfig(); err != nil {
      app.Print("Error writing site config")
			return ErrCode("0222", app.Site.siteFilePath)
		}
	}
	return nil
}
