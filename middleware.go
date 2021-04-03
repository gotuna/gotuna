package gotdd

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gorilla/mux"
)

const CORSAllowedOrigin = "*"
const CORSAllowedMethods = "OPTIONS,HEAD,GET,POST,PUT,PATCH,DELETE"

func (app App) Cors() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", CORSAllowedMethods)
				w.Header().Set("Access-Control-Allow-Origin", CORSAllowedOrigin)
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (app App) Logging() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			next.ServeHTTP(w, r)

			app.Logger.Printf("%s %s %s %s", start.Format(time.RFC3339), r.Method, r.URL.Path, time.Since(start))
		})
	}
}

func (app App) Recoverer(destination string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				if err := recover(); err != nil {
					stacktrace := string(debug.Stack())

					app.Logger.Printf("PANIC RECOVERED: %v", err)
					app.Logger.Println(stacktrace)

					// TODO: when templates are broken, redirecton can't work properly
					http.Redirect(w, r, destination, http.StatusInternalServerError)
					return
				}
			}()

			next.ServeHTTP(w, r)

		})
	}
}

func (app App) Authenticate(destination string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if app.Session.IsGuest(r) {
				http.Redirect(w, r, destination, http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app App) RedirectIfAuthenticated(destination string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if !app.Session.IsGuest(r) {
				http.Redirect(w, r, destination, http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
