package controllers

import (
	"book-store-be/database"
	"book-store-be/models"
	"book-store-be/responses"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var bookCollection *mongo.Collection = database.ConnectDatabase().Database("BOOK-STORE").Collection("Books")
var validate *validator.Validate = validator.New()

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

	// check validator
	validationError := validate.Struct(&book)
	if validationError != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error", validationError.Error())
		return
	}

	// adds new book
	_, err := bookCollection.InsertOne(ctx, book)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error", err.Error())
		return
	}

	responses.ResponseMessage(c, http.StatusOK, "success", "added new book")
}

// GetBook function that return a JSON with detail book
func GetBook(c *gin.Context) {
	// create a deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// set filter
	var book models.Book
	filter := bson.D{{
		Key:   "titolo",
		Value: c.Request.Header.Get("Title"),
	}}

	// finds a book with same "Titolo"
	err := bookCollection.FindOne(ctx, filter).Decode(&book)
	if err == mongo.ErrNoDocuments {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error", "book not found")
		return
	} else if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error", err.Error())
		return
	}

	responses.ResponseMessage(c, http.StatusOK, "success", book)
}
