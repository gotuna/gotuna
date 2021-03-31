package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alcalbg/gotdd/app"
	"github.com/alcalbg/gotdd/test/doubles"
	"github.com/alcalbg/gotdd/util"
)

func main() {

	port := ":8888"

	app := app.NewApp(util.Options{
		UserRepository: doubles.NewUserRepositoryStub(),
	})

	fmt.Printf("starting server at http://localhost%s \n", port)

	if err := http.ListenAndServe(port, app); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
