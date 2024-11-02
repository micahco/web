package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

type withError func(w http.ResponseWriter, r *http.Request) error

// Wraps handleWithError as http.HandlerFunc, with error handling
func (app *application) handle(h withError) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			app.logger.Error("handled unexpected error", slog.Any("err", err), slog.String("type", fmt.Sprintf("%T", err)))

			http.Error(w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
		}
	}
}

func (app *application) parseForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.Form)
	if err != nil {
		return err
	}

	err = app.validate.Struct(dst)
	if err != nil {
		return err
	}

	return nil
}
