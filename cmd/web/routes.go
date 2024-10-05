package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// App router
func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(app.recovery)

	r.Route("/", func(r chi.Router) {
		r.Get("/", app.handle(app.getIndex))

		r.Route("/articles", func(r chi.Router) {
			r.Get("/{id}", app.handle(app.getArticleID))
		})
	})

	return r
}

func (app *application) getIndex(w http.ResponseWriter, r *http.Request) error {
	return app.render(w, http.StatusOK, "welcome.html", nil)
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
