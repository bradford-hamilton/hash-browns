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
	// Set up error logging
	f, err := os.OpenFile("errors.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Create a new instance of pg database
	db, err := postgres.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create our server
	s := server.New(db)

	// Set up graceful shutdown and start server
	sigint := make(chan os.Signal, 1)
	idleConnsClosed := make(chan struct{})
	s.SigChan = sigint

	go func(sigint chan os.Signal) {
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		fmt.Println("Gracefully shutting down server...")

		// Received an interrupt signal, shut down
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
