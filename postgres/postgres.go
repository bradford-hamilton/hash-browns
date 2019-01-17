package postgres

import (
	"database/sql"
	"fmt"

	// postgres driver
	_ "github.com/lib/pq"
)

// Database is our database struct used for interacting with the database
type Database struct{ DB }

// DB defines the interface that a database must implement
type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Close() error
}

// New makes a new database using the connection string and
// returns it, otherwise returns the error
func New(connString string) (*Database, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	// Check that our connection is good
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

// ConnString returns a connection string based on the parameters it's given
// This would normally also contain the password, however we're not using one
func ConnString(host string, port int, user string, dbName string) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbName,
	)
}

// InsertReqTime inserts a request time into the database
func (db *Database) InsertReqTime(timeInMicroseconds int64) {
	query := `INSERT INTO req_times (time) VALUES ($1);`

	_, err := db.Exec(query, timeInMicroseconds)
	if err != nil {
		// TODO: ERror handling
		fmt.Println("Error inserting into database", err)
	}
}
