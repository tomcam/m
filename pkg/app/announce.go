package app

import (
	"fmt"
)

// App.Verbose() displays a message followed
// by a newline to stdout
// if the verbose flag was used. Formats it like Fprintf.
func (a *App) Verbose(format string, ss ...interface{}) {
	if a.flags.Verbose {
		fmt.Println(a.fmtMsg(format, ss...))
	}
}

// App.Note() displays a message followed by a newline
// to stdout, preceded by the text "NOTE: "
// For temporary use
// Overrides the verbose flag. Formats it like Fprintf.
func (a *App) Note(format string, ss ...interface{}) {
	fmt.Println("NOTE: " + a.fmtMsg(format, ss...))
}

// App.Warning() displays a message followed by a newline
// to stdout, preceded by the text "Warning: "
// Overrides the verbose flag. Formats it like Fprintf.
func (a *App) Warning(format string, ss ...interface{}) {
	fmt.Println("Warning: " + a.fmtMsg(format, ss...))
}

// fmtMsg() formats string like Fprintf and writes to a string
func (a *App) fmtMsg(format string, ss ...interface{}) string {
	return fmt.Sprintf(format, ss...)
}
