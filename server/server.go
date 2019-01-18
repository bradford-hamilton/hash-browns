package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/bradford-hamilton/hash-browns/pkg/hashbrown"
	"github.com/bradford-hamilton/hash-browns/postgres"
)

// Server struct will hold connection to the db, the server and handlers,
// as well as the sigChan needed for graceful shutdown
type Server struct {
	db      *postgres.Database
	Srv     *http.Server
	SigChan chan os.Signal
}

// New takes a db (*postgres.Database) and returns a new Server
// with db, server, and handlers
func New(db *postgres.Database) *Server {
	s := &Server{db: db}
	mux := http.NewServeMux()

	// Routing
	mux.HandleFunc("/hash", s.ReqTimer(s.Hash))
	mux.HandleFunc("/stats", s.Stats)
	mux.HandleFunc("/shutdown", s.Shutdown)

	s.Srv = &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("HB_PORT")),
		Handler: mux,
	}

	return s
}

// ReqTimer serves as middleware to time requests
func (s *Server) ReqTimer(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		microSecs := time.Since(start).Nanoseconds() / 1000
		s.db.InsertReqTime(microSecs)
	}
}

// Hash is the handler for http calls to /hash - It takes the form value with
// the key of "password", runs it through hash brown, and returns it in plain text
func (s *Server) Hash(w http.ResponseWriter, r *http.Request) {
	defer time.Sleep(5 * time.Second)

	if r.Method != "POST" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	p := r.PostFormValue("password")
	h := hashbrown.Create(p)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(h))
}

// Stats is the handler for http calls to /stats. It makes a GetStats() db
// call and returns the stats as json
func (s *Server) Stats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := s.db.GetStats()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

// Shutdown handles http calls to /shutdown and gracefully terminates the server
// by passing syscall.SIGTERM onto the servers SigChan
func (s *Server) Shutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	s.SigChan <- syscall.SIGTERM

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully shutdown server"))
}
