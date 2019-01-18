package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	// postgres driver
	_ "github.com/lib/pq"
)

// Database is our database struct used for interacting with the database
type Database struct{ DB }

// DB defines the interface that a database must implement
type DB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Close() error
}

// Stats represents the shape needed for statistics endpoint
type Stats struct {
	Total   int     `json:"total"`
	Average float64 `json:"average"`
}

// New makes a new database using the connection string and
// returns it, otherwise returns the error - Normally connStr
// would contain a password as well but no need for this challenge
func New() (*Database, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("HB_DB_HOST"),
		os.Getenv("HB_DB_PORT"),
		os.Getenv("HB_DB_USER"),
		os.Getenv("HB_DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
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

// InsertReqTime inserts a request time into the database
func (db *Database) InsertReqTime(timeInMicroseconds int64) {
	query := `INSERT INTO req_times (time) VALUES ($1);`

	_, err := db.Exec(query, timeInMicroseconds)
	if err != nil {
		log.Printf("Error inserting time of: %d into database: %s", timeInMicroseconds, err)
	}
}

// GetStats finds the total password hash requests as well as the average time
func (db *Database) GetStats() *Stats {
	query := `SELECT COUNT(time), AVG(time) FROM req_times;`
	stats := Stats{}
	var avg string

	row := db.QueryRow(query)

	// No need to handle error below - worst case scenario we
	// get 0 & 0 values for the very first request which would be accurate
	row.Scan(&stats.Total, &avg)

	// AVG returns a string with the pq driver and challenge example is
	// looking for int/float. No need to handle err as this will always be from pg
	f, _ := strconv.ParseFloat(avg, 64)
	stats.Average = f

	return &stats
}
