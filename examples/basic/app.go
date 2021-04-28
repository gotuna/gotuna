package basic

import (
	"fmt"
	"net/http"

	"github.com/gotuna/gotuna"
)

// MakeApp with sample dependencies.
func MakeApp(app gotuna.App) gotuna.App {

	app.Router = gotuna.NewMuxRouter()
	app.Router.Handle("/", handlerHome(app)).Methods(http.MethodGet)

	return app
}

func handlerHome(app gotuna.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Hello World!")
	})
}
