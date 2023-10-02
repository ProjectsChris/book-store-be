package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ConnectDatabase function create a new connection to the database
func ConnectDatabase(cn string) *mongo.Client {
	// create a deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// connect to database
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cn))
	if err != nil {
		panic(err.Error())
	}

	// does a ping
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err.Error())
	}

	return client
}

func GetCollection(client *mongo.Client, cn string) *mongo.Collection {
	return client.Database("BOOK-STORE").Collection(cn)
}

// InitDatabase function create a new connection to the Postgres database
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
