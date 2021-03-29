package middleware

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/alcalbg/gotdd/templating"
	"github.com/alcalbg/gotdd/util"
	"github.com/gorilla/mux"
)

func Logger(options util.Options) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

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

			options.Logger.Printf("%s %s %s %s", start.Format(time.RFC3339), r.Method, r.URL.Path, time.Since(start))
		})
	}
}
