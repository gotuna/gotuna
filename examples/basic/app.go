package basic

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gotuna/gotuna"
)

func MakeApp(app gotuna.App) gotuna.App {

	app.Router = mux.NewRouter()
	app.Router.Handle("/", handlerHome(app)).Methods(http.MethodGet)

	return app
}

func handlerHome(app gotuna.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})
}
