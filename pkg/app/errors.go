package app

// Example usage
//
//
//	if err = app.changeWorkingDir(tmpDir); err != nil {
//    msg := fmt.Sprintf("System error attempting to change to new site directory %s: %s", requested, err.Error())
//    return ErrCode("1111", msg)
//  }
//	msg := fmt.Sprintf("%s for project %s: %s", dir, pathname, err.Error())
//	  return ErrCode("0414", msg)
//   return ErrCode("1234", err.Error())
//
//   return ErrCode("PREVIOUS", err.Error())
//
//   return ErrCode("0401", err.Error(), filename)
//
//	 if err != nil {
//     return (ErrCode("1033", from, err.Error()))
//   }
//
//   Example (a very good example) from util.go
//	 if err != nil {
//     return ErrCode("1234", "from '"+source+"' to '"+dest+"'", "")
//   }
//
//   err = copyDirAll(App.themesPath, App.Site.themesPath)
//   if err != nil {
// 	   QuitError(ErrCode("0911", "from '"+App.themesPath+"' to '"+App.Site.themesPath+"'"))
//   }
//
//   if err := copyDirOnly(from, to); err != nil {
//     msg := fmt.Sprintf("Unable to copy from pageType directory %s to new pageType directory %s", from, to)
//     return ErrCode("0906", msg)
//   }
//
//    if err := RemoveAll(tmpDir); err != nil {
//      msg := fmt.Sprintf("System error attempting to delete temporary site directory %s: %s", tmpDir, err.Error())
//      return ErrCode("0601", msg)
//    }

//
import (
	"fmt"
	"github.com/tomcam/m/pkg/default"
	"os"
)

// defaults.ErrorCodePrefix is a short string (currently "mbz") used
// to make it easier to search a Metabuz error code on the web.

//	SECTIONS
//
//	0100	- Error reading file
//	0200	- Error creating file
//	0300	- Error deleting file
//	0400 	- Error creating directory
//	0500	- Error determining directory name
//	0600	- Error deleting directory
//	0700	- Error reading directory
//	0800	- Can't determine the name of something
//	0900	- Problem generating something
//  0950 - Something's already there
//	1000	- Something's missing that should be there
//	1100	- Problem changing to a directory
//  1200  - Syntax error!
//  1300  - Error writing to file
//
// The reason many of these error codes have identical text is that
// the same error occurs but in different places. Since the
// Go lib returns identical error messages for each one, tracking
// down the error code shows us where the error occurred even if the
// executable is stripped of debug info
var errMsgs = map[string]string{

	// Just print the last error
	"PREVIOUS": " ",

	// 0100	- Error reading file
	"0101": "Error reading front matter",            // filename
	"0112": "Unable to copy file",                   // filename
	"0113": "Error reading site configuration file", // filename
	"0114": "Unable to open file",
	"0115": "Unable to find starter file",                          // filename
	"0116": "Error reading site configuration file",                // filename
	"0117": "Error unmarshalling YAML for site configuration file", // filename
	// Old errors stopped at 0131
	// TODO: Get rid of the line below
	// https://github.com/tomcam/mb/blob/master/pkg/errs/errors.go
	"0118": "No site configuration file specified",
	"0132": "Unable to open theme configuration file",     // filename
	"0133": "Unable to open theme configuration file",     // filename
	"0134": "Error unmarshalling YAML for new theme file", // filename
	"0135": "Error in starter file",                       // custom message

	// 0200	- Error creating file
	"0209": "Unable to create the file", // dest, source filenames
	// TODO: Get rid of the line below
	// Old errors stopped at 0215
	"0217": "Can't publish stylesheet to same location",                  // filename
	"0218": "Can't create starter file",                                  // filename
	"0219": "Can't create site file",                                     // filename
	"0220": "Can't create site file",                                     // filename
	"0221": "Can't copy a file onto itself!",                             // filename
	"0222": "Can't create site file",                                     // filename
	"0223": "Unable to rename theme file",                                // custom message
	"0224": "Unable to rename stylesheet",                                // custom message
	"0225": "Unable to create config file for new theme",                 // filename
	"0226": "System error attempting to rename temporary site directory", // custom message
	"0227": "System error attempting to write site config file to",       // custom message
	"0228": "System error attempting to create YAML file",
	"0229": "Unable to create post file", // filename

	// 0250 - Error closing file
	// Old errors stopped at 0252

	// 0300	- Error deleting file
	"0302": "Unable to delete publish directory",
	// TODO: Get rid of the line below
	// Old errors stopped at 0302

	// 0400	- Error creating directory
	"0401": "Unable to create project directory", // filename
	"0403": "Unable to create publish directory",
	"0406": "Unable to copy site directory",
	"0409": "Error creating directory for theme", // Fully qualified pathname
	"0410": "Error creating directory",           // directory
	// TODO: Get rid of the line below
	// Old errors stopped at 0410
	"0411": "Unable to copy theme directory",
	"0412": "Error creating directory for simple page", // directory
	"0413": "Error creating directory for new theme",   // Fully qualified pathname
	"0414": "System error creating temp directory in ", // custom message
	"0415": "Error creating directory for posts",       // directory
	"0416": "Error creating directory for a post",      // directory
	"0417": "Error creating directory for starter",     // directory
	// TODO: Get rid of the line below
	// https://github.com/tomcam/mb/blob/master/pkg/errs/errors.go

	// 0500	- Error determining directory name
	"0501": "Couldn't get relative directory name", // custom message I think

	// 0600 - Error deleting directory
	"0601": "System error attempting to delete temporary site directory", // dirname

	// 0700	- Error reading directory
	"0706": "Unable to read from headtags directory", // Expected pathname of headtags directory
	// Old errors stopped at 0709
	"0709": "Unable to read from script directory " + defaults.ScriptClosePath, // Expected pathname of scripts directory
	// 0800	- Can't determine the name of something
	"0801": "",

	// 0900	- Problem generating something
	"0901": "Unable to allocate Site object",
	"0902": "Error creating new site file",
	"0915": "Error copying theme to site", // filename
	"0917": "Problem parsing",             // filename
	"0920": "Error generating Markdown",
	"0921": "Unable to build project", // filename
	"0922": "No project found at",     // message
	// TODO: Get rid of the line below
	// Old errors stopped at 0924
	"0923": "Error building",             // projectname
	"0924": "Error creating new project", // projectname
	// TODO: Get rid of the line below
	// https://github.com/tomcam/mb/blob/master/pkg/errs/errors.go
	"0925": "Error generating Markdown for page elment file", // filename
	"0926": "Unable to generate table of contents",
	"0927": "Unable to create a new page",
	"0928": "Unable to create a new post",
	"0929": "Error copying theme",         // Custom message
	"0930": "Error updating copied theme", // Custom message
	"0931": "Error copying theme",         // Custom message
	"0932": "Error populating the " + defaults.CfgDir + " directory",
	// 0950 - Something's already there
	"0951": "Site already exists at",         // sitename
	"0952": "Theme already exists at",        // sitename
	"0954": "Duplicate path for collection.", // posts directory name
	"0955": "There is already a file named",  // filename
	"0956": "There is already a post named",  // post name
	// 1000	- Something's missing that should be there
	"1001": "Missing front matter and markdown", // filename
	"1002": "This isn't a project directory:",   // directoryname
	"1004": "Trying to publish nonexistent stylesheet",
	"1005": "No publish directory specified for",
	"1013": "Please specify a site name",
	"1014": "No destination file specified when copying", // source file to copy

	// TODO: Get rid of the line below
	// Old errors stopped at 1023
	"1024": "Couldn't find stylesheet",        //filename
	"1025": "This isn't a project directory:", // directoryname
	"1026": "This isn't a project directory:", // directoryname
	//"1027": "File specified in theme configuration file is missing:", // filename
	"1027": "Theme configuration file",           // filename
	"1028": "Can't find a theme named",           // filename, theme name
	"1029": "Can't find the theme file",          // filename
	"1033": "Unable to read theme directory",     // filename
	"1034": "Unable to find layout element file", // filename
	"1035": "Missing name of theme to copy",      // filename
	"1036": "No site configuration file was specified",
	"1037": "Filename missing to create page",
	"1038": "Can't find the permalink for",       // custom message
	"1039": "Author name missing from permalink", // permalink
	// TODO: Get rid of the line below
	// https://github.com/tomcam/mb/blob/master/pkg/errs/errors.go

	//	1100	- Problem changing to a directory
	"1101": "Unable to", // chdir to bad dir name
	"1102": "Unable to", // chdir to bad dir name during showInfo()
	"1103": "Unable to", // chdir to bad dir name during newSite()
	"1104": "Missing name for starter page",
	"1105": "Unable to", // chdir to bad dir name
	// Old errors stopped at 1106
	// TODO: Get rid of the line below
	"1107": "Can't change to site directory", // project name
	"1108": "Can't change to site directory",
	"1109": "Can't change to site directory for interview",                             // directory name
	"1110": "Can't change to site directory to copy theme",                             // directory name
	"1111": "System error changing to newly create site directory",                     // directory name
	"1112": "System error changing to newly create site directory",                     // directory name
	"1113": "Can't change to site directory to create new post at",                     // directory name
	"1114": "Can't change to site directory to create new post at",                     // directory name
	"1115": ":author permalink variable specified but no author has been specified in", // collection

	// 1200 - Syntax error!
	"1204": "Unknown dot value in Go template function ", //
	"1205": "Error generating table of contents",         //
	// TODO: Get rid of the line below
	// Old errors stopped at 1206
	"1207": "Don't understand the starter type", // Name of type in a starter file
	///"1208": "Your path must start with a normal directory name, for example, `blog` or `news`, but it starts with the permalink variable", // Permalink variable
	"1208": "Path to collection starts with", // custom message
	"1209": "Starter file",                   // has Unknown permalink variable

	// 1300 -  Error writing to file
	"1301": "Unable to update site file with collection", // Collection name
	"1032": "Error creating page",                        // filename
	"1303": "Error creating page",                        // filename

}

type ErrMsg struct {
	// key to a map of error messages
	key string

	// If key is the word "PREVIOUS", this will contain an error
	// message from an earlier action, typically a return from the
	// Go runtime.
	previous string
	extra    string
	system   string
}

// Error() looks up e.key, which is an error code number
// expressed as a string (for example, "1001")
// and returns its associated map entry, which is an explanatory
// text message for that error code.
// But there's likely much more happening:
// -  If e.key is "PREVIOUS" it suggests that an error message
//    that didn't get displayed probably
//    should be displayed, and its contents
//    in e.previous are returned.
// -  If e.extra has something, say, a filename, it should be
//    used to customize the error message.
// -  If the e.key isn't recognized, it displays an
//    "error code untracked" error message as a reminder to me
//    that I screwed up.
// Why is the key a number formatted as a string?
// Because it gets appended to "mbz" in an error message,
// and I plan for Metabuzz to be so popular that people would be
// looking up error codes using search engines, e.g. mbz1001. And it's a
// ghetto way of keeping error codes unique while being kind of sorted
// in the source code.
func (e *ErrMsg) Error() string {
	var msg error
	// Make sure the error code has documentation
	if errMsgs[e.key] == "" {
		msg = fmt.Errorf("ERROR CODE %s UNTRACKED: please contact the developer\nMore info: %s\n",
			defaults.ErrorCodePrefix+e.key, e.previous)
		return msg.Error()
	}

	// Handle this special case:
	// ErrCode("1234", "PREVIOUS", something)
	// This is the case where something like a Go system call
	// returned an error, but I want to know where it occurred.
	// TODO: Revive this idea
	/*
		if e.previous == "PREVIOUS" {
			msg = fmt.Errorf("%s (error code %s%s). Previous error was '%s'",
				errMsgs[e.key], defaults.ErrorCodePrefix, e.key, e.extra)
			//fmt.Printf("HEY %v", msg.Error())
			return msg.Error()
		}
	*/
	//xxx

	// Error message from an earlier error return needs to be seen.
	if e.key == "PREVIOUS" {
		return fmt.Errorf("%s\n", e.previous).Error()
	}

	if e.previous != "" {
		msg = fmt.Errorf("%s %s (error code %s%s)\n",
			// The most common case: an error code with customization
			errMsgs[e.key], e.previous, defaults.ErrorCodePrefix, e.key)
	} else {
		msg = fmt.Errorf("%s (error code %s%s)\n",
			// The slightly unusual case of an error code with no customization
			errMsgs[e.key], defaults.ErrorCodePrefix, e.key)
	}
	return msg.Error()
}

// ErrCode() takes an error code, say "0110", and
// one or two optional strings. It adds the error code
// to the error messages map so that message can be looked
// up. The additional parameters give information such
// as whether a previous error message should be displayed,
// or something to customize the message that's already in
// the error messages map, like a filename.
// When calling a Go runtime routine that could return
// an error message, make err.Error() the second
// parameter so its contents are included, like this:
//
// Sample usages:
//
//   return ErrCode("PREVIOUS", err.Error())
//
//   return ErrCode("1234", err.Error())
//
//   return ErrCode("0401", err.Error(), filename)
//
// TODO: Make sure I implemented this or delete the comment
//   return ErrCode("0401", "PREVIOUS", errText)
//
//   Good example from util.go
//	 if err != nil {
//     return ErrCode("1234", "from '"+source+"' to '"+dest+"'", "")
//   }

//
//   err = copyDirAll(App.themesPath, App.Site.themesPath)
//   if err != nil {
// 	   QuitError(ErrCode("0911", "from '"+App.themesPath+"' to '"+App.Site.themesPath+"'"))
//   }
//
//
//   if err := copyDirOnly(from, to); err != nil {
//     msg := fmt.Sprintf("Unable to copy from pageType directory %s to new pageType directory %s", from, to)
//     return ErrCode("0906", msg)
//   }
//
//	 msg := fmt.Errorf("Error attempting to create project file %s: %v", projectFile, err.Error()).Error()
//
func ErrCode(key string, previous string, extra ...string) error {
	var e error
	if len(extra) > 0 {
		e = add(key, previous, extra[0])
	} else {
		e = add(key, previous)
	}
	return e
}

// add() allocates a map entry for errMsgs.
func add(key string, previous string, extra ...string) error {
	switch len(extra) {
	case 1:
		return &ErrMsg{key, previous, extra[0], ""}
	case 2:
		return &ErrMsg{key, previous, extra[0], extra[1]}
	default:
		return &ErrMsg{key, previous, "", ""}
	}

}

// QuitError() displays the error passed to it and exits
// to the operating system, returning a 1 (any nonzero
// return means an error occurred).
// Normally functions that can generate a runtime error
// do so by returning an error. But sometimes there's a
// constraint, for example, fulfilling an interface method
// that doesn't support this practice.
func (a *App) QuitError(e error) {
	// Precede error message with name of file
	if a.Page.filePath != "" {
		fmt.Print(a.Page.filePath + ": ")
	}
	// Error message from an earlier error return needs to be seen.
	displayError(e)
	if e == nil {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// displayError() shows the specified error message
// without exiting to the OS.
func displayError(e error) {
	fmt.Println(e.Error())
}
