package main

import (
	"SolBaa/websockets/internal/handlers"
	"net/http"

	"github.com/bmizerany/pat"
)

func routes() http.Handler {
	mux := pat.New()
	mux.Get("/", handlers.Home)

}
