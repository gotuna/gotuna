package basic

import (
	"fmt"
	"net/http"

	"github.com/alcalbg/gotdd"
	"github.com/gorilla/mux"
)

func MakeApp(app gotdd.App) gotdd.App {

	app.Router = mux.NewRouter()
	app.Router.Handle("/", handlerHome(app)).Methods(http.MethodGet)

	return app
}

func handlerHome(app gotdd.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})
}
