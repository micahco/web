package main

import (
	"log/slog"
	"net/http"

	"github.com/justinas/nosurf"
)

func (app *application) recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.logger.Error("recovered from panic", slog.Any("error", err))

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; frame-ancestors 'self'; form-action 'self';")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) csrfFailureHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Error("csrf failuie",
			slog.String("method", r.Method),
			slog.String("uri", r.URL.RequestURI()),
		)

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	})
}

func (app *application) noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	csrfHandler.SetFailureHandler(app.csrfFailureHandler())

	return csrfHandler
}

// TODO:
// func (app *application) authenticate(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		id := app.sessionManager.GetInt(r.Context(), authenticatedUserIDSessionKey)
// 		if id == 0 {
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		exists, err := app.models.User.Exists(id)
// 		if err != nil {
// 			app.serverError(w, r, err)
// 			return
// 		}

// 		if exists {
// 			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
// 			r = r.WithContext(ctx)
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// func (app *application) requireAuthentication(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if !app.isAuthenticated(r) {
// 			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
// 			return
// 		}

// 		w.Header().Add("Cache-Control", "no-store")

// 		next.ServeHTTP(w, r)
// 	})
// }
