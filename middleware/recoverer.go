package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/alcalbg/gotdd/templating"
	"github.com/alcalbg/gotdd/util"
	"github.com/gorilla/mux"
)

func Recoverer(options util.Options) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				if err := recover(); err != nil {
					stacktrace := string(debug.Stack())

					options.Logger.Printf("PANIC RECOVERED: %v", err)
					options.Logger.Println(stacktrace)

					//fmt.Println(err, stacktrace)

					w.WriteHeader(http.StatusInternalServerError)
					templating.GetEngine(options). // TODO lang per user
									Set("error", err).
									Set("stacktrace", string(debug.Stack())).
									Render(w, r, "app.html", "error.html")
				}
			}()

			next.ServeHTTP(w, r)

		})
	}
}
