package postgres

import (
	"database/sql"
	"fmt"
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
	Total   int
	Average float64
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

// GetStats finds the total password hash requests as well as the average time
func (db *Database) GetStats() *Stats {
	query := `SELECT COUNT(time), AVG(time) FROM req_times;`
	stats := Stats{}
	var avg string

	row := db.QueryRow(query)
	err := row.Scan(&stats.Total, &avg)
	if err != nil {
		// TODO: ERror handling
		fmt.Println("Error selecting stats from db", err)
	}

	// AVG returns a string with the pq driver and challenge example is looking for int/float
	f, _ := strconv.ParseFloat(avg, 64)
	stats.Average = f

	return &stats
}
