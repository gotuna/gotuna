package middleware

import (
	"net/http"

	"github.com/alcalbg/gotdd/util"
	"github.com/gorilla/mux"
)

func RedirectIfAuthenticated(options util.Options, destination string) mux.MiddlewareFunc {
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
