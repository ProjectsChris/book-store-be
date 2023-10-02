package controllers

import (
	db "book-store-be/database"
	"book-store-be/models"
	"book-store-be/responses"
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoClient struct {
	Db *mongo.Client
}

type DatabaseSql struct {
	Db *sql.DB
}

// var bookCollection *mongo.Collection = database.ConnectDatabase().Database("BOOK-STORE").Collection("Books")
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
func (dp *DatabaseSql) PostBook(c *gin.Context) {
	// take values from body
	book := new(models.Book)
	if err := c.BindJSON(book); err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
		return
	}

	// check validator
	validationError := validate.Struct(book)
	if validationError != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+validationError.Error())
		return
	}

	query := `INSERT INTO books (title, writer, price, summary, cover, genre, quantity, IdCover) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := dp.Db.Exec(query, book.Titolo, book.Autore, book.Prezzo, book.Summary, book.Copertina, book.Genere, book.Quantita, book.Id)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
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
func (mc *MongoClient) GetBook(c *gin.Context) {
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
	// Get collection
	bookCollection := db.GetCollection(mc.Db, "Books")

	// finds a book with same "Titolo"
	err := bookCollection.FindOne(ctx, filter).Decode(&book)
	if err == mongo.ErrNoDocuments {
		responses.ResponseMessage(c, http.StatusNotFound, "error: book not found")
		return
	} else if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, &book)
}
