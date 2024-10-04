package main

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/micahco/web/internal/httperr"
)

type handleWithError func(w http.ResponseWriter, r *http.Request) error

// Wraps handleWithError as http.HandlerFunc, with error handling.
func (app *application) handle(handler handleWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			// First, check if error is httperr.
			var httpErr *httperr.Error
			if errors.As(err, &httpErr) {
				// Log wrapped error if exists.
				if err := errors.Unwrap(httpErr); err != nil {
					app.logger.Error(
						"handled unwrapped error",
						slog.Any("error", err),
					)
				}

				http.Error(w, httpErr.Message(), httpErr.Code())
				return
			}

			// Else, send generic internal server error.
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			app.logger.Error(
				"handled unexpected error",
				slog.Any("error", err),
			)
		}
	}
}
