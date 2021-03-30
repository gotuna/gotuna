package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

const CORSAllowedOrigin = "*"
const CORSAllowedMethods = "OPTIONS,HEAD,GET,POST,PUT,PATCH,DELETE"

func Cors() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", CORSAllowedMethods)
				w.Header().Set("Access-Control-Allow-Origin", CORSAllowedOrigin)
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
