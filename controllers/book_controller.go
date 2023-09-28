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

// PostBook godoc
//
//	@Summary	Get a book
//	@Schemes
//	@Description	Get details of a book
//	@Tags			Book
//	@Accept			json
//	@Produce		json
//	@Param			models.Book	body		models.Book	true "Add new book"
//	@Success		200			{object}	models.Book
//	@Failure		500			{object}	responses.ResponseErrorJSON
//	@Router			/book/ [post]
//
// PostBook function add new book
func PostBook(c *gin.Context) {
	// create a deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// take values from body
	var book models.Book
	if err := c.BindJSON(&book); err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: " +  err.Error())
		return
	}

	// check validator
	validationError := validate.Struct(&book)
	if validationError != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: " + validationError.Error())
		return
	}

	// adds new book
	_, err := bookCollection.InsertOne(ctx, book)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+ err.Error())
		return
	}

	responses.ResponseMessage(c, http.StatusOK, "success: added new book")
}

// GetBook godoc
//
//	@Summary	Get a book
//	@Schemes
//	@Description	Get details of a book
//	@Tags			Book
//	@Accept			json
//	@Produce		json
//	@Param			title	path		string	true	"Title of the book"
//	@Success		200		{object}	models.Book
//	@Failure		404		{object}	responses.ResponseErrorJSON
//	@Failure		500		{object}	responses.ResponseErrorJSON
//	@Router			/book/{title} [get]
//
// GetBook function that return a JSON with detail book
func GetBook(c *gin.Context) {
	// create a deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//utilizziamo gli endpoint non gli header
	bookTitle := c.Param("title")

	// set filter
	book := new(models.Book)
	filter := bson.D{{
		Key:   "titolo",
		Value: bookTitle,
	}}

	// finds a book with same "Titolo"
	err := bookCollection.FindOne(ctx, filter).Decode(&book)
	if err == mongo.ErrNoDocuments {
		responses.ResponseMessage(c, http.StatusNotFound, "error: book not found")
		return
	} else if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, &book)
}
