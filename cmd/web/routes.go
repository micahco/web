package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

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
	})

	return r
}

type handler func(w http.ResponseWriter, r *http.Request) error

func (app *application) handle(handler handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			app.errorLog.Printf("%v", err)

			http.Error(w, err.Error(), 500)
		}
	}
}

func (app *application) getIndex(w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintf(w, "hello world!")

	return nil
}
