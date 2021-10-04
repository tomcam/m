package errs
import (
	"fmt"
	"github.com/tomcam/m/pkg/default"
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
//	1000	- Something's missing that should be there
//	1100	- Problem changing to a directory
//  1200  - Syntax error!
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
	"0101": "Error reading front matter",           // filename
	// 0200	- Error creating file
	// 0250 - Error closing file
	// 0300	- Error deleting file
	// 0400	- Error creating directory
	// 0500	- Error determining directory name
	// 0600 - Error deleting directory
	// 0700	- Error reading directory
	// 0800	- Can't determine the name of something
	"0801": "",
	// 0900	- Problem generating something
	// 0950 - Something's already there
	// 1000	- Something's missing that should be there
	"1001": "Missing front matter and markdown", // filename
	// 1100 - Problem changing to a directory
	// 1200 - Syntax error!
}

type ErrMsg struct {
	// Key to a map of error messages
	Key string

	// If Key is the word "PREVIOUS", this will contain an error
	// message from an earlier action, typically a return from the
	// Go runtime.
	Previous string
	Extra    []string
}

// Error() looks up e.Key, which is an error code number 
// expressed as a string (for example, "1001") 
// and returns its associated map entry, which is an explanatory
// text message for that error code.
// But there's likely much more happening:
// -  If e.Key is "PREVIOUS" it suggests that an error message
//    that didn't get displayed probably
//    should be displayed, and its contents 
//    in e.previous are returned.
// -  If e.Extra has something, say, a filename, it should be
//    used to customize the error message.
// -  If the e.Key isn't recognized, it displays an 
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
	if errMsgs[e.Key] == "" {
		msg = fmt.Errorf("ERROR CODE %s UNTRACKED: please contact the developer\nMore info: %s\n",
			defaults.ErrorCodePrefix+e.Key, e.Previous)
		return msg.Error()
	}

	// Error message from an earlier error return needs to be seen.
	if e.Key == "PREVIOUS" {
		return fmt.Errorf("%s\n", e.Previous).Error()
	}

	if e.Previous != "" {
		msg = fmt.Errorf("%s %s (error code %s%s)\n",
			// The most common case: an error code with customization
			errMsgs[e.Key], e.Previous, defaults.ErrorCodePrefix, e.Key)
	} else {
		msg = fmt.Errorf("%s (error code %s%s)\n",
			// The slightly unusual case of an error code with no customization
			errMsgs[e.Key], defaults.ErrorCodePrefix, e.Key)
	}
	return msg.Error()
}

// new() allocates a map entry for errMsgs.
func new(key string, previous string, extra ...string) error {
	return &ErrMsg{key, previous, extra}
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
//   return ErrCode("0401", err.Error())
//
//   return ErrCode("0401", err.Error(), filename)
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
		e = new(key, previous, extra[0])
	} else {
		e = new(key, previous)
	}
	return e
}
