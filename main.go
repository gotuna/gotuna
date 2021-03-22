package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alcalbg/gotdd/app"
	"github.com/gorilla/sessions"
)

func main() {

	port := ":8888"

	logger := log.New(os.Stdout, "", 0)

	cookieStore := sessions.NewCookieStore([]byte(os.Getenv("APP_KEY")))
	srv := app.NewServer(logger, cookieStore)

	fmt.Printf("starting server at http://localhost%s \n", port)

	if err := http.ListenAndServe(port, srv.Router); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
