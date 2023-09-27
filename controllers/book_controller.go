package controllers

import (
	"book-store-be/database"
	"book-store-be/models"
	"book-store-be/responses"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var bookCollection *mongo.Collection = database.ConnectDatabase().Database("BOOK-STORE").Collection("Books")
var counterCollection *mongo.Collection = database.ConnectDatabase().Database("BOOK-STORE").Collection("Counters")

// PostBook function add new book
func PostBook(c *gin.Context) {
	// create a deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// take values from body
	var book models.Book
	if err := c.BindJSON(&book); err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error", err.Error())
		return
	}

	// auto increment id
	var auto_increment map[string]any = bson.M{}
	err := counterCollection.FindOneAndUpdate(
		ctx,
		bson.D{{
			Key:   "_id",
			Value: "bookId",
		}},
		bson.D{{
			Key: "$inc",
			Value: bson.D{{
				Key:   "seq",
				Value: 1,
			}},
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true),
	).Decode(&auto_increment)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error", err.Error())
		return
	}

	// change value of the id
	book.Id = auto_increment["seq"].(int32)

	// adds new book
	_, err = bookCollection.InsertOne(ctx, book)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error", err.Error())
		return
	}

	responses.ResponseMessage(c, http.StatusOK, "success", "added new book")
}

// GetBooks return a JSON with a list of books
func GetBooks(c *gin.Context) {
	// create a deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// find all books
	cur, err := bookCollection.Find(ctx, bson.D{})
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error", err.Error())
		return
	}
	defer cur.Close(ctx)

	// create an array
	var books []models.Book

	// execute all content of cursor
	for cur.Next(ctx) {
		var b models.Book

		err := cur.Decode(&b)
		if err != nil {
			responses.ResponseMessage(c, http.StatusInternalServerError, "error", err.Error())
			return
		}

		// add new book inside an array
		books = append(books, b)
	}
	if err = cur.Err(); err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error", err.Error())
		return
	}

	responses.ResponseMessage(c, http.StatusOK, "success", books)
}
