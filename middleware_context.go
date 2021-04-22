package gotuna

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

// StoreParamsToContext middleware will add all parameters from the current
// request to the context, this includes query, form, and route params
func (app App) StoreParamsToContext() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()

			params := url.Values{}

			// Parse route vars
			vars := mux.Vars(r)
			for k, v := range vars {
				params.Add(k, v)
			}

			// Parse form and add to params
			if err := r.ParseForm(); err == nil {
				for k, v := range r.Form {
					for _, vv := range v {
						params.Add(k, vv)
					}
				}
			}

			next.ServeHTTP(w, r.WithContext(ContextWithParams(ctx, params)))
		})
	}
}

// StoreUserToContext middleware will add the current logged in
// user object (if any) to the request context for further use.
func (app App) StoreUserToContext() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()

			// Next, get the user ID if logged in
			if app.Session == nil {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			userID, err := app.Session.GetUserID(r)

			if err != nil {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Get the full user object from the UserRepository
			user, err := app.UserRepository.GetUserByID(userID)

			if err != nil {
				// Something went wrong, authenticated user object cannot be retrieved
				// from the repo. Destroy this session and logout the user.
				app.Session.Destroy(w, r)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ctx = ContextWithUser(ctx, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
