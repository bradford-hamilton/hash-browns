package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	s.KillItWithFire = sigint

	go func(sigint chan os.Signal) {
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := s.Srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}(sigint)

	if err := s.Srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Printf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
