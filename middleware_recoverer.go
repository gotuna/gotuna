package gotdd

import (
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
)

func (app App) Recoverer(destination string) mux.MiddlewareFunc {
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
