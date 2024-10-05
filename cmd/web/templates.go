package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	"github.com/micahco/web/ui"
)

type templateData struct {
	CSRFToken       string
	CurrentYear     int
	Flash           string
	IsAuthenticated bool
	Data            interface{}
}

// Render page template with data
func (app *application) render(w http.ResponseWriter, statusCode int, page string, data interface{}) error {
	td := templateData{
		CurrentYear: time.Now().Year(),
		//Flash:           app.popFlash(r),
		//IsAuthenticated: app.isAuthenticated(r),
		//CSRFToken: nosurf.Token(r),
		Data: data,
	}

	// In production, use template cache
	if !app.config.dev {
		return app.renderFromCache(w, statusCode, page, td)
	}

	// In development, parse files locally
	t, err := template.ParseFiles("./ui/html/base.html")
	if err != nil {
		return err
	}

	t, err = t.Funcs(functions).ParseGlob("./ui/html/partials/*.html")
	if err != nil {
		return err
	}

	t, err = t.ParseFiles("./ui/html/pages/" + page)
	if err != nil {
		return err
	}

	return writeTemplate(t, td, w, statusCode)
}

func (app *application) renderFromCache(w http.ResponseWriter, statusCode int, page string, td templateData) error {
	t, ok := app.templateCache[page]
	if !ok {
		return fmt.Errorf("template %s does not exist", page)
	}

	return writeTemplate(t, td, w, statusCode)
}

func writeTemplate(t *template.Template, td templateData, w http.ResponseWriter, statusCode int) error {
	buf := new(bytes.Buffer)

	err := t.ExecuteTemplate(buf, "base", td)
	if err != nil {
		return err
	}

	w.WriteHeader(statusCode)

	if _, err := buf.WriteTo(w); err != nil {
		return err
	}

	return nil
}

var functions = template.FuncMap{}

// Create new template cache with ui.Files embedded file system.
// Creates a template for each page in the html/pages directory
// nested with html/base.html and html/partials.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	fsys := ui.Files

	// Get list of pages
	pages, err := fs.Glob(fsys, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// Create a new template for each page and add to cache map.
	for _, page := range pages {
		name := filepath.Base(page)

		// Nest page with base template and partials
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(fsys, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
