package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alcalbg/gotdd/server"
)

func main() {

	port := ":8888"
	server := server.NewServer()

	fmt.Printf("starting server at http://localhost%s \n", port)

	if err := http.ListenAndServe(port, server.Router); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
