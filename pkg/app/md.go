package app

import (
	"bytes"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"text/template"
)

// mdFiletoHTML converts a Markdown document
// WITHOUT YAML front matter to HTML.
// Returns a byte slice containing the HTML source.
func (app *App) mdFileToHTML(filename string) ([]byte, error) {
	source := fileToBuf(filename)
	return app.mdToHTML(source)
}

// mdYAMLFiletoHTML converts a Markdown document
// with YAML front matter to HTML.
// Returns a byte slice containing the HTML source.
func (app *App) mdYAMLFileToHTML(filename string) ([]byte, error) {
	source := fileToBuf(filename)
	return app.mdYAMLToHTML(source)
}


// mdYAMLtoHTML converts a Markdown document with optional
// YAML front matter to HTML. YAML is written to app.metaData
// Returns a byte slice containing the HTML source.
// Pre: parser.NewContext() has already been called on app.parserCtx
func (app *App) mdYAMLToHTML(source []byte) ([]byte, error) {
	var buf bytes.Buffer
	// Convert Markdown source to HTML and deposit in buf.Bytes().
	if err := app.mdParser.Convert(source, &buf, parser.WithContext(app.mdParserCtx)); err != nil {
		return []byte{}, err
	}
	// Obtain YAML front matter from document.
	app.metaData = meta.Get(app.mdParserCtx)
	return buf.Bytes(), nil
}

// mdYAMLtoHTMLStr converts a Markdown document with optional YAML front matter to HTML. YAML is written to app.metaData
// Returns a string containing the HTML source.
// Pre: parser.NewContext() has already been called on app.parserCtx
func (app *App) mdYAMLToHTMLStr(source []byte) (string, error) {
	var buf bytes.Buffer
	// Convert Markdown source to HTML and deposit in buf.Bytes().
	if err := app.mdParser.Convert(source, &buf, parser.WithContext(app.mdParserCtx)); err != nil {
		return "", err
	}
	// Obtain YAML front matter from document.
	app.metaData = meta.Get(app.mdParserCtx)
	return string(buf.Bytes()), nil
}

// mdtoHTML converts a Markdown document to HTML.
// YAML front matter should not be present.
// Returns a byte slice containing the HTML source.
// Pre: parser.NewContext() has already been called on app.parserCtx
func (app *App) mdToHTML(source []byte) ([]byte, error) {
	var buf bytes.Buffer
	// Convert Markdown source to HTML and deposit in buf.Bytes().
	if err := app.mdParser.Convert(source, &buf, parser.WithContext(app.mdParserCtx)); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

// doTemplate takes HTML in source, expects parsed front
// matter in app.metaData, and executes Go templates
// against the source.
// Returns a string containing the HTML with the
// template values embedded.
func (app *App) doTemplate(templateName string, source string) (string, error) {
	if templateName == "" {
		templateName = "Metabuzz"
	}
	tmpl, err := template.New(templateName).Parse(source)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, app.metaData)

	if err != nil {
		return "", err
	}
	return buf.String(), err

}

// doTemplateFuncs takes HTML in source, expects parsed front
// matter in app.metaData, and executes Go templates
// against the source. It also handles user-defined
// functions, expected in funcMap
// Returns a string containing the HTML with the
// template values embedded.
func (app *App) doTemplateFuncs(templateName string, source string) (string, error) {
	if templateName == "" {
		templateName = "Metabuzz"
	}
	var tmpl *template.Template
	var err error
	if tmpl, err = template.New(templateName).Funcs(app.funcs).Parse(source); err != nil {
		// TODO: Function should return error
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, templateName, app.metaData)

	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
