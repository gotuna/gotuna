package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/static"
	"github.com/alcalbg/gotdd/test/doubles"
	"github.com/gorilla/sessions"
)

func main() {

	port := ":8888"
	keyPairs := os.Getenv("APP_KEY")

	app := gotdd.NewApp(gotdd.App{
		Logger:         log.New(os.Stdout, "", 0),
		UserRepository: doubles.NewUserRepositoryStub(),
		Session:        gotdd.NewSession(sessions.NewCookieStore([]byte(keyPairs))),
		FS:             static.EmbededStatic,
		StaticPrefix:   "",
	})

	fmt.Printf("starting server at http://localhost%s \n", port)

	if err := http.ListenAndServe(port, app.Router); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
