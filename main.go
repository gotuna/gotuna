package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alcalbg/gotdd/app"
	"github.com/alcalbg/gotdd/models"
	"github.com/alcalbg/gotdd/session"
	"github.com/alcalbg/gotdd/static"
	"github.com/gorilla/sessions"
)

func main() {

	port := ":8888"

	gorillaSessionStore := sessions.NewCookieStore([]byte(os.Getenv("APP_KEY")))
	fs := static.EmbededStatic

	srv := app.NewServer(
		log.New(os.Stdout, "", 0),
		fs,
		session.NewSession(gorillaSessionStore),
		models.NewInMemoryUserRepository(),
	)

	fmt.Printf("starting server at http://localhost%s \n", port)

	if err := http.ListenAndServe(port, srv); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
