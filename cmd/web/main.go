package main

import (
	"log"
	"net/http"

	"github.com/SolBaa/websockets/internal/handlers"
)

func main() {
	routes := routes()
	log.Println("Starter channel listener")
	go handlers.ListenforWSChannel()
	http.ListenAndServe(":8000", routes)
}
