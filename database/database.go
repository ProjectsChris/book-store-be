package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

// InitDatabase creates new connection to the database
func InitDatabase(cn string) *sql.DB {
	// parsing string
	// uri, err := pq.ParseURL(cn)
	// open database
	db, err := sql.Open("postgres", cn)
	if err != nil {
		log.Fatal("Impossibile connettersi al database")
	}

	// check database
	if err = db.Ping(); err != nil {
		log.Fatal("error: ", err.Error())
	}

	return db
}
