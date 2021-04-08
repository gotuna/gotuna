package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alcalbg/gotdd"
	"github.com/alcalbg/gotdd/examples/basic"
)

func main() {

	port := ":8888"

	app := basic.MakeApp(gotdd.App{})

	fmt.Printf("starting server at http://localhost%s \n", port)

	if err := http.ListenAndServe(port, app.Router); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}
