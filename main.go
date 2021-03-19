package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alcalbg/gotdd/server"
)

func main() {

	port := ":8888"

	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	srv := server.NewServer(logger)

	fmt.Printf("starting server at http://localhost%s \n", port)

	if err := http.ListenAndServe(port, srv.Router); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
