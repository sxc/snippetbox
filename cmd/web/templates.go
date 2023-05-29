package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/sxc/snippetbox/internal/models"
	"github.com/sxc/snippetbox/ui"
)

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")

}

var functions = template.FuncMap{

	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Loop through the page filepaths one-by-one.
	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil

	// Rigistered with the template set
	// 	ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
	// 	if err != nil {
	// 		return nil, err

	// 	}

	// 	ts, err = ts.ParseFiles(page)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	cache[name] = ts
	// }
	// return cache, nil
}

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}
