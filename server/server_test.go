package server

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/bradford-hamilton/hash-browns/postgres"
)

type mockSQLResult struct {
	Query  string
	Timing interface{}
}

func (msr mockSQLResult) LastInsertId() (int64, error) { return 1, nil }
func (msr mockSQLResult) RowsAffected() (int64, error) { return 1, nil }

type mockDB struct{}

func (m mockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	msql := mockSQLResult{}
	return msql, nil
}

func (m mockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	return &sql.Row{}
}

func (m mockDB) Close() error { return nil }

func TestServer_Hash(t *testing.T) {
	m := &postgres.Database{mockDB{}}
	s := &Server{db: m}

	/// Create a request with expected form values to pass to our handler
	req, err := http.NewRequest("POST", "/hash", nil)
	if err != nil {
		t.Fatal(err)
	}
	form := url.Values{}
	form.Add("password", "angryMonkey")
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.Hash)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestServer_ReqTimer(t *testing.T) {
	m := &postgres.Database{mockDB{}}
	s := &Server{db: m}

	req, err := http.NewRequest("POST", "/hash", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.ReqTimer(func(w http.ResponseWriter, r *http.Request) {}))
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestServer_Shutdown(t *testing.T) {
	sigint := make(chan os.Signal, 1)
	defer close(sigint)

	m := &postgres.Database{mockDB{}}
	s := &Server{db: m, SigChan: sigint}

	req, err := http.NewRequest("GET", "/shutdown", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.Shutdown)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Wait for SIGTERM on s.SigChan and if we dont see it within 5 seconds, exit and fail
	select {
	case sig := <-s.SigChan:
		if sig != syscall.SIGTERM {
			t.Fatalf("signal was %v, want %v", sig, syscall.SIGTERM)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("timeout waiting for %v", syscall.SIGTERM)
	}
}
