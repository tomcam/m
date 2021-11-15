package main

import (
	"html/template"
	"net/http"
)

func (app *app) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	paths := []string{
		"./ui/html/body.tmpl",
		"./ui/html/footer.tmpl",
		"./ui/html/header.tmpl",
	}

	tmpl, err := template.ParseFiles(paths...)
	if err != nil {
		// TODO: Use my error handling style
		app.errLog.Println(err.Error())
		http.Error(w, "Internal server error parsing templates", 500)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		// TODO: Use my error handling style
		app.errLog.Println(err.Error())
		http.Error(w, "Internal server error executing templates", 500)
		return
	}
}

func (app *app) newPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// TC: use my error conventions
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "newPage(): "+r.Method+" Method not allowed", 405)
		return
	}
}
