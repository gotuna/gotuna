package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/examples/fullapp"
	"github.com/gotuna/gotuna/examples/fullapp/i18n"
	"github.com/gotuna/gotuna/examples/fullapp/static"
	"github.com/gotuna/gotuna/examples/fullapp/views"
)

func main() {

	port := ":8888"
	keyPairs := os.Getenv("APP_KEY")

	app := fullapp.MakeApp(gotuna.App{
		Logger:         log.New(os.Stdout, "", 0),
		UserRepository: fullapp.NewUserRepository(),
		Session:        gotuna.NewSession(sessions.NewCookieStore([]byte(keyPairs))),
		Static:         static.EmbededStatic,
		StaticPrefix:   "",
		ViewFiles:      views.EmbededViews,
		Locale:         gotuna.NewLocale(i18n.Translations),
	})

	// production only, do not use in tests
	app.Router.Use(
		csrf.Protect(
			[]byte(keyPairs),
			csrf.FieldName("csrf_token"),
			csrf.CookieName("csrf_token"),
		))

	fmt.Printf("starting server at http://localhost%s \n", port)

	if err := http.ListenAndServe(port, app.Router); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}
