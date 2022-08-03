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

// mdFileToTemplatedHTMLString converts a Markdown document
// to HTML. Executes any templates.
// Returns a string containing the HTML source.
func (app *App) mdFileToTemplatedHTMLString(filename string) (string, error) {
	source := fileToBuf(filename)
  var b []byte
  var err error
  if b, err  = app.mdToHTML(source); err != nil {
   // TODO: Real error code
	  return "",err
  }
  var s string
	if s, err = app.doTemplateFuncs(filename, string(b)); err != nil {
    // TODO: Real error code
	  return "",err
  }
  return s, nil
}

// mdFiletoTemplatedFuncsHTMLString converts a Markdown document
// to HTML. Executes any templates using the funcmap specified in funcs.
// Returns a string containing the HTML source.
// xxx
func (app *App) mdFileToTemplatedFuncsHTMLString(filename string, funcs map[string]interface{}) (string, error) {
  var s string
  var err error
  if s, err  = app.mdFileToTemplatedHTMLString(filename); err != nil {
   // TODO: Real error code
	  return "",err
  }
	var tmpl *template.Template

  // TODO: Refactor
	if tmpl, err = template.New(filename).Funcs(funcs).Parse(s); err != nil {
		// TODO: Function should return proper error code
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, filename, app)


	if err != nil {
		// TODO: Function should return proper error code
		return "", err
	}
	return buf.String(), nil
}

// mdYAMLFiletoTemplatedFuncsHTMLString converts a Markdown document
// to HTML. Executes any templates using the funcmap specified in funcs.
// Returns a string containing the HTML source.
// xxx
func (app *App)mdYAMLFileToTemplatedFuncsHTMLString(filename string, funcs map[string]interface{}) (string, error) {
  var s string
  var err error
  if s, err  = app.mdYAMLFileToTemplatedHTMLString(filename); err != nil {
   // TODO: Real error code
	  return "",err
  }
	var tmpl *template.Template

  // TODO: Refactor
	if tmpl, err = template.New(filename).Funcs(funcs).Parse(s); err != nil {
		// TODO: Function should return proper error code
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, filename, app)


	if err != nil {
		// TODO: Function should return proper error code
		return "", err
	}
	return buf.String(), nil
}



// mdYAMLFiletoTemplatedHTMLString converts a Markdown document
// to HTML. Executes any templates.
// Returns a string containing the HTML source.
func (app *App) mdYAMLFileToTemplatedHTMLString(filename string) (string, error) {
	source := fileToBuf(filename)
  var b []byte
  var err error
  if b, err  = app.mdYAMLToHTML(source); err != nil {
   // TODO: Real error code
	  return "",err
  }
  var s string
	if s, err = app.doTemplateFuncs(filename, string(b)); err != nil {
    // TODO: Real error code
	  return "",err
  }
  return s, nil
}



// mdYAMLtoHTML converts a Markdown document with optional
// YAML front matter to HTML. YAML is written to app.metaData
// and it populates app.Page.FrontMatter.
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
  // Copy populated fields to app.Page.FrontMatter
  app.frontMatterRawToStruct()
  app.Print("mdYAMLToHTML(): %#v", app.Page.FrontMatter)
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
	err = tmpl.Execute(buf, app)

	if err != nil {
		return "", err
	}
	return buf.String(), err

}

// doTemplateFuncs takes HTML in source, expects parsed front
// matter in app.metaData, and executes Go templates
// against the source. It also handles user-defined
// functions, expected in app.funcs
// templateName is expected to be the filename of the Markdown source file
// Returns a string containing the HTML with the
// template values embedded.
// TODO: Refactor to call doTemplate()?
func (app *App) doTemplateFuncs(templateName string, source string) (string, error) {
	if templateName == "" {
		templateName = "Metabuzz"
	}
	var tmpl *template.Template
	var err error
  // xxx
	if tmpl, err = template.New(templateName).Funcs(app.funcs).Parse(source); err != nil {
		// TODO: Function should return proper error code
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, templateName, app)

	if err != nil {
		// TODO: Function should return proper error code
		return "", err
	}
	return buf.String(), nil
}
