package gotdd

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gorilla/mux"
)

const CORSAllowedOrigin = "*"
const CORSAllowedMethods = "OPTIONS,HEAD,GET,POST,PUT,PATCH,DELETE"

func Cors() mux.MiddlewareFunc {
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

func Logger(options Options) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			next.ServeHTTP(w, r)

			options.Logger.Printf("%s %s %s %s", start.Format(time.RFC3339), r.Method, r.URL.Path, time.Since(start))
		})
	}
}

func Recoverer(options Options) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				if err := recover(); err != nil {
					stacktrace := string(debug.Stack())

					options.Logger.Printf("PANIC RECOVERED: %v", err)
					options.Logger.Println(stacktrace)

					//fmt.Println(err, stacktrace)

					w.WriteHeader(http.StatusInternalServerError)
					GetEngine(options). // TODO lang per user
								Set("error", err).
								Set("stacktrace", string(debug.Stack())).
								Render(w, r, "app.html", "error.html")
				}
			}()

			next.ServeHTTP(w, r)

		})
	}
}

func Authenticate(options Options, destination string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if options.Session.IsGuest(r) {
				http.Redirect(w, r, destination, http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RedirectIfAuthenticated(options Options, destination string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if !options.Session.IsGuest(r) {
				http.Redirect(w, r, destination, http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
