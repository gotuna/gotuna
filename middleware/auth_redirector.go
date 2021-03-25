package middleware

import (
	"net/http"
	"strings"

	"github.com/alcalbg/gotdd/session"
	"github.com/alcalbg/gotdd/util"
	"github.com/gorilla/mux"
)

func AuthRedirector(s *session.Session) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			path := r.URL.Path

			if strings.HasPrefix(path, util.StaticPath) {
				next.ServeHTTP(w, r)
				return
			}

			guest := s.IsGuest(r)

			if guest && path != "/login" && path != "/register" {
				http.Redirect(w, r, "/login", http.StatusFound)
			}

			if !guest && (path == "/login" || path == "/register") {
				http.Redirect(w, r, "/", http.StatusFound)
			}

			next.ServeHTTP(w, r)
		})
	}
}
