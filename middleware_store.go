package gotuna

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

// StoreToContext middleware will add common values to the context for further use
// this includes all of the parameters for the current request query/form/route
// and the current logged in user (if any)
func (app App) StoreToContext() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()

			params := url.Values{}

			// parse route vars
			vars := mux.Vars(r)
			for k, v := range vars {
				params.Add(k, v)
			}

			// parse form and add to params
			if err := r.ParseForm(); err == nil {
				for k, v := range r.Form {
					for _, vv := range v {
						params.Add(k, vv)
					}
				}
			}

			ctx = ContextWithParams(ctx, params)

			// next, get the user if logged in
			if app.Session == nil {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			userID, err := app.Session.GetUserID(r)

			if err != nil {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if user, err := app.UserRepository.GetUserByID(userID); err == nil {
				ctx = ContextWithUser(ctx, user)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
