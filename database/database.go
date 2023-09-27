package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// TODO: change values with Mongo Atlas
const uri string = "mongodb://localhost:27017"

// ConnectDatabase function create a new connection to the database
func ConnectDatabase() *mongo.Client {
	// create a deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// connect to database
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err.Error())
	}

	// does a ping
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err.Error())
	}

	return client
}
