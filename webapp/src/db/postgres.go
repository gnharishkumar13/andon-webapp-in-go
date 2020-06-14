package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	//Load the driver for side-effect
	_ "github.com/lib/pq"
)

const (
	// TimeLayoutTimestamp is the layout that should be used when parsing Timestamps that are returned from Postgres queries
	TimeLayoutTimestamp = "2006-01-02T15:04:05Z"
)

var (
	db *sql.DB
	m  sync.Mutex
)

func GetDB() (*sql.DB, error) {
	m.Lock()
	var err error

	if db == nil {
		connStr := "user=postgres password=postgres dbname=andon host=db sslmode=disable"
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Print("Error connecting to db", err)
			return nil, fmt.Errorf("Error connecting to db %v ", err)
		}
	}
	m.Unlock()
	return db, nil
}

// FromTimestamp takes a timestamp from the database and attempts to convert it to a time.Time.
func FromTimestamp(timestamp string) (time.Time, error) {
	return time.Parse(TimeLayoutTimestamp, timestamp)
}

// ToTimestamp takes a time and converts it to a format that the database will understand.
func ToTimestamp(t time.Time) string {
	return t.Format(TimeLayoutTimestamp)
}
