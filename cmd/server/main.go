package main

import (
	"context"
	"fmt"
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

	sigint := make(chan os.Signal, 1)
	idleConnsClosed := make(chan struct{})

	s.SigChan = sigint

	go func(sigint chan os.Signal) {
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		// Nice little notification
		fmt.Println("Gracefully shutting down server...")

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
