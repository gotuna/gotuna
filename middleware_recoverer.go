package gotuna

import (
	"net/http"
	"runtime/debug"
)

// Recoverer middleware is used to recover the app from panics, to log the
// incident, and to redirect user to the error page.
func (app App) Recoverer(destination string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				if err := recover(); err != nil {
					stacktrace := string(debug.Stack())

					app.Logger.Printf("PANIC RECOVERED:\n%v", err)
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
