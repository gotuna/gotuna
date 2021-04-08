package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/examples/fullapp"
	"github.com/alcalbg/gotdd/examples/fullapp/i18n"
	"github.com/alcalbg/gotdd/examples/fullapp/static"
	"github.com/alcalbg/gotdd/examples/fullapp/views"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
)

func main() {

	port := ":8888"
	keyPairs := os.Getenv("APP_KEY")

	app := fullapp.MakeApp(gotdd.App{
		Logger:         log.New(os.Stdout, "", 0),
		UserRepository: fullapp.NewUserRepository(),
		Session:        gotdd.NewSession(sessions.NewCookieStore([]byte(keyPairs))),
		Static:         static.EmbededStatic,
		StaticPrefix:   "",
		ViewFiles:      views.EmbededViews,
		Locale:         gotdd.NewLocale(i18n.Translations),
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
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
