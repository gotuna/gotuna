package gotdd

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (app App) Logging() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			next.ServeHTTP(w, r)

			if app.Logger != nil {
				app.Logger.Printf("%s %s %s %s", start.Format(time.RFC3339), r.Method, r.URL.Path, time.Since(start))
			}
		})
	}
}
