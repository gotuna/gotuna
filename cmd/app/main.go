package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/test/doubles"
)

func main() {

	port := ":8888"

	app := gotdd.App{
		UserRepository: doubles.NewUserRepositoryStub(),
		Logger:         log.New(os.Stdout, "", 0),
	}

	fmt.Printf("starting server at http://localhost%s \n", port)

	if err := http.ListenAndServe(port, app.Router); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
