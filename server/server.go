package server

import (
	"net/http"
	"time"

	"github.com/bradford-hamilton/hash-browns/pkg/hashbrown"
)

// Server will hold connection to the db as well as handlers
type Server struct{}

// New returns a new Server with db dependency
func New() *Server {
	return &Server{}
}

// Hash handles http calls to /hash
func (s *Server) Hash() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "405 Method Not Allowed", 405)
			return
		}

		p := r.PostFormValue("password")
		h := hashbrown.Create(p)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(h))

		time.Sleep(5 * time.Second)
	}
}
