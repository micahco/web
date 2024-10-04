package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/micahco/web/internal/httperr"
)

// App router
func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.NotFound(http.NotFound)

	r.Route("/", func(r chi.Router) {
		r.Get("/", app.handle(app.getIndex))

		r.Route("/articles", func(r chi.Router) {
			r.Get("/{id}", app.handle(app.getArticleID))
		})
	})

	return r
}

func (app *application) getIndex(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "hello world!")

	return nil
}

func (app *application) getArticleID(w http.ResponseWriter, r *http.Request) error {
	p := chi.URLParam(r, "id")

	id, err := strconv.Atoi(p)
	if err != nil {
		return httperr.New(http.StatusBadRequest, "invalid id param")
	}

	fmt.Fprintf(w, "Article ID: %d", id)

	return nil
}
