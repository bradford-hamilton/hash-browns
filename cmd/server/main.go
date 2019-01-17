package main

import (
	"log"
	"net/http"

	"github.com/bradford-hamilton/hash-browns/server"
)

func main() {
	s := server.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/hash", s.Hash())

	// Serve
	log.Fatal(http.ListenAndServe(":4000", mux))
}
