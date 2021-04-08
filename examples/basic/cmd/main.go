package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gotuna/gotuna"
	"github.com/gotuna/gotuna/examples/basic"
)

func main() {

	port := ":8888"

	app := basic.MakeApp(gotuna.App{})

	fmt.Printf("starting server at http://localhost%s \n", port)

	if err := http.ListenAndServe(port, app.Router); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}
