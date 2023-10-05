package routes

import (
	"book-store-be/models"
	"book-store-be/observability"
	"book-store-be/responses"
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel"
)

type DatabaseSql struct {
	Db *sql.DB
}

var validate = validator.New()
var meter = otel.Meter("book-counter")

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
//	@Failure		400			{object}	responses.ResponseErrorJSON
//	@Failure		500			{object}	responses.ResponseErrorJSON
//	@Router			/book/ [post]
//
// PostBook function add new book
func (ds *DatabaseSql) PostBook(c *gin.Context) {
	t := observability.Tracer
	_, span := t.Tracer("").Start(context.Background(), "post-book")
	defer span.End()

	// take values from body
	book := new(models.Book)
	if err := c.BindJSON(book); err != nil {
		responses.ResponseMessage(c, http.StatusBadRequest, "error: "+err.Error())
		return
	}

	// check validator
	validationError := validate.Struct(book)
	if validationError != nil {
		responses.ResponseMessage(c, http.StatusBadRequest, "error: "+validationError.Error())
		return
	}

	query := `INSERT INTO books (title, writer, price, summary, cover, genre, quantity, IdCover) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := ds.Db.Exec(query, book.Titolo, book.Autore, book.Prezzo, book.Summary, book.Copertina, book.Genere, book.Quantita, book.IdCopertina)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
		return
	}

	meterCounter, err := meter.Int64Counter("post-book-counter")
	if err != nil {
		panic(err.Error())
	}

	meterCounter.Add(context.Background(), 1)

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
func (ds *DatabaseSql) GetBook(c *gin.Context) {
	t := observability.Tracer
	_, span := t.Tracer("").Start(context.Background(), "get-book")
	defer span.End()

	book := new(models.Book)
	bookTitle := c.Param("title")

	// create a query
	query := `SELECT * FROM books WHERE id = $1`

	meterCounter, err := meter.Int64Counter("get-book-counter")
	if err != nil {
		panic(err.Error())
	}

	meterCounter.Add(context.Background(), 1)

	res, err := ds.Db.Query(query, bookTitle)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
		return
	}
	defer res.Close()

	for res.Next() {
		err = res.Scan(&book.Id, &book.Titolo, &book.Autore, &book.Prezzo, &book.Summary, &book.Copertina, &book.Genere, &book.Quantita, &book.IdCopertina)
		if err != nil {
			responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
			return
		}
	}

	// check if is an empty struct
	if *book == (models.Book{}) {
		responses.ResponseMessage(c, http.StatusNotFound, "book not found")
	} else {
		c.JSON(http.StatusOK, &book)
	}
}
