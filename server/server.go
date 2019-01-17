package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bradford-hamilton/hash-browns/pkg/hashbrown"
	"github.com/bradford-hamilton/hash-browns/postgres"
)

// Server will hold connection to the db as well as handlers
type Server struct {
	db *postgres.Database
}

// New returns a new Server with db dependency
func New(db *postgres.Database) *Server {
	return &Server{db}
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

// Hash handles http calls to /hash
func (s *Server) Hash() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

// Stats handles http calls to /stats
func (s *Server) Stats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		s := s.db.GetStats()

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s)
	}
}
