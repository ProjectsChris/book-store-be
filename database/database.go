package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

// InitDatabase function create a new connection to the database
func InitDatabase(cn string) *sql.DB {
	// open database
	db, err := sql.Open("postgres", cn)
	if err != nil {
		log.Fatal("error: ", err.Error())
	}

	// check database
	if err = db.Ping(); err != nil {
		log.Fatal("error: ", err.Error())
	}

	return db
}
