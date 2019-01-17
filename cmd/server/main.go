package main

import (
	"log"
	"net/http"

	"github.com/bradford-hamilton/hash-browns/postgres"
	"github.com/bradford-hamilton/hash-browns/server"
)

func main() {
	// Create a new connection to our pg database
	db, err := postgres.New(
		postgres.ConnString("localhost", 5432, "bradfordlamson-scribner", "hash_browns_db"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s := server.New(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/hash", s.ReqTimer(s.Hash()))

	log.Fatal(http.ListenAndServe(":4000", mux))
}
