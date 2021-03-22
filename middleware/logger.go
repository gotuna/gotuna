package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/alcalbg/gotdd/util"
	"github.com/gorilla/mux"
)

func Logger(logger *log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// log every request
			logger.Println(time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

			defer func() {
				if err := recover(); err != nil {
					// log error and stack trace to console
					logger.Printf("PANIC RECOVERED: %v", err)
					logger.Println(string(debug.Stack()))

					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(util.DefaultError))
					return
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
