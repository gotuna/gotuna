package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/alcalbg/gotdd/i18n"
	"github.com/alcalbg/gotdd/templating"
	"github.com/gorilla/mux"
)

func Logger(logger *log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			defer func() {
				if err := recover(); err != nil {
					stacktrace := string(debug.Stack())

					logger.Printf("PANIC RECOVERED: %v", err)
					logger.Println(stacktrace)

					//fmt.Println(err, stacktrace)

					w.WriteHeader(http.StatusInternalServerError)
					templating.GetEngine(i18n.NewTranslator(i18n.En), nil). // TODO lang
												Set("error", err).
												Set("stacktrace", string(debug.Stack())).
												Render(w, r, "app.html", "error.html")
				}
			}()

			next.ServeHTTP(w, r)

			logger.Printf("%s %s %s %s", start.Format(time.RFC3339), r.Method, r.URL.Path, time.Since(start))
		})
	}
}
