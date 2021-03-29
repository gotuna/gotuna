package middleware

import (
	"net/http"

	"github.com/alcalbg/gotdd/session"
	"github.com/gorilla/mux"
)

func AuthRedirector(s *session.Session, guestRoutes map[string]string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			path := r.URL.Path
			_, isGuestRoute := guestRoutes[path]

			// serving static file?
			if path != "/" && !isGuestRoute {
				next.ServeHTTP(w, r)
				return
			}

			guest := s.IsGuest(r)

			if guest && !isGuestRoute {
				http.Redirect(w, r, "/login", http.StatusFound)
			}

			if !guest && isGuestRoute {
				http.Redirect(w, r, "/", http.StatusFound)
			}

			next.ServeHTTP(w, r)
		})
	}
}
