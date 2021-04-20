package gotuna

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Authenticate middleware will redirect all guests to the destination.
// This is used to guard user-only routes and to force guests to login.
func (app App) Authenticate(destination string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if app.Session.IsGuest(r) {
				http.Redirect(w, r, destination, http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RedirectIfAuthenticated middleware will redirect authenticated
// users to the destination.
// This is used to deflect logged in users from guest-only pages like login
// or register page back to the app.
func (app App) RedirectIfAuthenticated(destination string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if !app.Session.IsGuest(r) {
				http.Redirect(w, r, destination, http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
