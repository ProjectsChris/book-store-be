package database

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
)

// InitDatabase function create a new connection to the database
func InitDatabase(cn string) *sql.DB {
	// open database
	uri, err := pq.ParseURL(cn)
	//sql.Open("postgres", cn)
	if err != nil {
		log.Fatal("error: ", err.Error())
	}
	db, err := sql.Open("postgres", uri)
	// check database
	if err = db.Ping(); err != nil {
		log.Fatal("error: ", err.Error())
	}

	return db
}
