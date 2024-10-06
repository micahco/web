package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/micahco/web/ui"
)

// App router
func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(app.recovery)
	r.Use(secureHeaders)

	// Static files
	r.Handle("/static/*", app.handleStatic())
	r.Get("/favicon.ico", app.handleFavicon)

	r.Route("/", func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave)
		r.Use(app.noSurf)
		r.Use(app.authenticate)

		r.Get("/", app.handle(app.getIndex))

		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", app.handle(app.handleAuthLoginPost))
			r.Post("/logout", app.handle(app.handleAuthLogoutPost))
			r.Post("/signup", app.handle(app.handleAuthSignupPost))
			r.Get("/register", app.handle(app.handleAuthRegisterGet))
			r.Post("/register", app.handle(app.handleAuthRegisterPost))
			r.Get("/reset", app.handle(app.handleAuthResetGet))
			r.Post("/reset", app.handle(app.handleAuthResetPost))
			r.Get("/reset/update", app.handle(app.handleAuthResetUpdateGet))
			r.Post("/reset/update", app.handle(app.handleAuthResetUpdatePost))
		})

		r.Route("/articles", func(r chi.Router) {
			r.Use(app.requireAuthentication)

			r.Get("/{id}", app.handle(app.getArticleID))
		})
	})

	return r
}

func (app *application) refresh(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func (app *application) handleStatic() http.Handler {
	if app.config.dev {
		fs := http.FileServer(http.Dir("./ui/static/"))

		return http.StripPrefix("/static", fs)
	}

	return http.FileServer(http.FS(ui.Files))
}

func (app *application) handleFavicon(w http.ResponseWriter, r *http.Request) {
	if app.config.dev {
		http.ServeFile(w, r, "./ui/static/favicon.ico")

		return
	}
	http.ServeFileFS(w, r, ui.Files, "static/favicon.ico")
}

func (app *application) getIndex(w http.ResponseWriter, r *http.Request) error {
	if app.isAuthenticated(r) {
		return app.render(w, r, http.StatusOK, "dashboard.tmpl", nil)
	}

	return app.render(w, r, http.StatusOK, "login.tmpl", nil)
}

func (app *application) getArticleID(w http.ResponseWriter, r *http.Request) error {
	p := chi.URLParam(r, "id")

	id, err := strconv.Atoi(p)
	if err != nil {
		return respErr{nil, http.StatusBadRequest, "invalid id param"}
	}

	fmt.Fprintf(w, "Article ID: %d", id)

	return nil
}
