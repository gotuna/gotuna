package middleware

import (
	"net/http"

	"github.com/alcalbg/gotdd/util"
	"github.com/gorilla/mux"
)

func AuthRedirector(options util.Options) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			path := r.URL.Path
			_, isGuestRoute := options.GuestRoutes[path]

			// serving a static file?
			if path != "/" && !isGuestRoute {
				next.ServeHTTP(w, r)
				return
			}

			guest := options.Session.IsGuest(r)

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
