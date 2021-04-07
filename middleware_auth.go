package gotdd

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app App) Authenticate(destination string) mux.MiddlewareFunc {
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

func (app App) RedirectIfAuthenticated(destination string) mux.MiddlewareFunc {
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

func (app App) StoreUserToContext() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID, err := app.Session.GetUserID(r)

			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if user, err := app.UserRepository.GetUserByID(userID); err == nil {
				next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), user)))
				return
			}

			next.ServeHTTP(w, r)
			return
		})
	}
}
