package app

import (
	"bytes"
	"text/template"
)

// Resolve template variables
// input is an HTML file that includes entities like {{.FrontMatter.Description}}
// Replace with the appropriate values in generated output.
// The filename is passed in because it
// produces an accurate location of any
// source file parsing errors that occur.
// Skips if frontmatter.Templates is set to "off"
func (app *App) interps(filename string, input string) string {
	//if strings.ToLower(app.Page.FrontMatter.Templates) != "off" {
	if !app.Page.FrontMatter.Templates {
		return app.execute(filename, input, app.funcs)
	}
	return input
}

// execute() parses a Go template, then executes it against HTML/template source.
// Returns a string containing the result.
func (a *App) execute(templateName string, tpl string, funcMap template.FuncMap) string {
	var t *template.Template
	var err error
	if t, err = template.New(templateName).Funcs(funcMap).Parse(tpl); err != nil {
		// TODO: Function should return error
		a.QuitError(ErrCode("0917", err.Error()))
	}
	var b bytes.Buffer
	err = t.ExecuteTemplate(&b, templateName, a)
	if err != nil {
		// TODO: Handle error properly & and document error code
		a.QuitError(ErrCode("1204", err.Error()))
	}
	return b.String()
}
