package main

import (
	"log/slog"
	"net/http"
)

type withError func(w http.ResponseWriter, r *http.Request) error

// Wraps handleWithError as http.HandlerFunc, with error handling
func (app *application) handle(h withError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			app.logger.Error(
				"handled unexpected error",
				slog.Any("err", err),
			)
		}
	}
}
