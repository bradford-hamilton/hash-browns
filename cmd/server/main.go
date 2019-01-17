package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bradford-hamilton/hash-browns/postgres"
	"github.com/bradford-hamilton/hash-browns/server"
)

func main() {
	db, err := postgres.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s := server.New(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/hash", s.ReqTimer(s.Hash()))
	mux.HandleFunc("/stats", s.Stats())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("HB_PORT")), mux))
}
