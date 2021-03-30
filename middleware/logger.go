package middleware

import (
	"net/http"
	"time"

	"github.com/alcalbg/gotdd/util"
	"github.com/gorilla/mux"
)

func Logger(options util.Options) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			next.ServeHTTP(w, r)

			options.Logger.Printf("%s %s %s %s", start.Format(time.RFC3339), r.Method, r.URL.Path, time.Since(start))
		})
	}
}
