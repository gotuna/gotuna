package gotdd

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app App) StoreUserToContext() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userID, err := app.Session.GetUserID(r)

			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()

			if user, err := app.UserRepository.GetUserByID(userID); err == nil {
				ctx = ContextWithUser(ctx, user)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
