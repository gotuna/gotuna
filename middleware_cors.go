package gotuna

import (
	"net/http"
)

// Cors middleware will add CORS headers to OPTIONS request and respond with 204 status.
func (app App) Cors() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,HEAD,GET,POST,PUT,PATCH,DELETE")
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
