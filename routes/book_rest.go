package routes

import (
	"book-store-be/database"
	"book-store-be/models"
	"book-store-be/responses"
	"context"
	"database/sql"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type DatabaseSql struct {
	Db *sql.DB
}

var validate = validator.New()
var meter = otel.Meter("book-counter")

// PostBook godoc
//
//	@Summary	Adds a book
//	@Schemes
//	@Description	Adds new book
//	@Tags			Book
//	@Accept			json
//	@Produce		json
//	@Param			models.Book	body		models.Book	true	"Adds new book"
//	@Success		200			{object}	models.Book
//	@Failure		400			{object}	responses.ResponseErrorJSON
//	@Failure		500			{object}	responses.ResponseErrorJSON
//	@Router			/book/ [post]
//
// PostBook function add new book
func (ds *DatabaseSql) PostBook(c *gin.Context) {
	_, span := otel.Tracer("").Start(context.Background(), "/api/v1/book/")
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

	query := `INSERT INTO books (title, writer, price, summary, cover, genre, quantity, category, IdCover) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := ds.Db.Exec(query, book.Titolo, book.Autore, book.Prezzo, book.Summary, book.Copertina, book.Genere, book.Quantita, book.Categoria, book.IdCopertina)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
		return
	}

	// init a metric
	meterCounter, err := meter.Int64Counter("post-book-counter")
	if err != nil {
		panic(err.Error())
	}

	meterCounter.Add(c.Request.Context(), 1)
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
	// create a span
	_, span := otel.Tracer("").Start(c.Request.Context(), "/api/v1/book/id")
	defer span.End()

	book := new(models.Book)
	bookTitle := c.Param("title")

	// create a query
	query := `SELECT * FROM books WHERE id = $1`

	// init a meter counter
	meterCounter, err := meter.Int64Counter("get-book-counter")
	if err != nil {
		panic(err.Error())
	}

	res, err := ds.Db.Query(query, bookTitle)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
		return
	}
	defer res.Close()

	for res.Next() {
		err = res.Scan(&book.Id, &book.Titolo, &book.Autore, &book.Prezzo, &book.Summary, &book.Copertina, &book.Genere, &book.Quantita, &book.Categoria, &book.IdCopertina)
		if err != nil {
			responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
			return
		}
	}

	// check if is an empty struct
	if *book == (models.Book{}) {
		responses.ResponseMessage(c, http.StatusNotFound, "book not found")

		meterCounter.Add(c.Request.Context(), 1, metric.WithAttributes(
			attribute.String("status", strconv.Itoa(http.StatusNotFound)),
		))
	} else {
		c.JSON(http.StatusOK, &book)

		meterCounter.Add(c.Request.Context(), 1, metric.WithAttributes(
			attribute.String("status", strconv.Itoa(http.StatusOK)),
		))
	}
}

// GetBooks godoc
//
//	@Summary	Get books
//	@Schemes
//	@Description	Get details for every book
//	@Tags			Book
//	@Accept			json
//	@Produce		json
//	@Param			page	query		string	true	"Number of the pagination"
//	@Success		200		{object}	responses.ResponseDatabase  // TODO: fix type
//	@Failure		404		{object}	responses.ResponseErrorJSON
//	@Failure		500		{object}	responses.ResponseErrorJSON
//	@Router			/book [get]
//
// GetBooks function that return a JSON with all books
func (ds *DatabaseSql) GetBooks(c *gin.Context) {
	var query string
	var bookList []models.Book

	book := new(models.Book)
	counter := 0

	// create a span
	_, span := otel.Tracer("").Start(c.Request.Context(), "/api/v1/book/")
	defer span.End()

	// query param
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

	// init a meter counter
	meterCounter, err := meter.Int64Counter("get-books-counter")
	if err != nil {
		panic(err.Error())
	}

	// creates a query
	query = `SELECT * FROM books ORDER BY id DESC LIMIT 10 OFFSET $1;`

	res, err := ds.Db.Query(query, 10*page)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
		return
	}

	// execute all books
	for res.Next() {
		err = res.Scan(&book.Id, &book.Titolo, &book.Autore, &book.Prezzo, &book.Summary, &book.Copertina, &book.Genere, &book.Quantita, &book.Categoria, &book.IdCopertina)
		if err != nil {
			responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
			return
		}

		bookList = append(bookList, *book)
	}

	// query for show a count of all elements into database
	query = database.GET_ALL
	res, err = ds.Db.Query(query)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
		return
	}
	defer res.Close()

	for res.Next() {
		err = res.Scan(&counter)
		if err != nil {
			responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
			return
		}
	}

	// check length of the list
	if len(bookList) > 0 {
		c.JSON(http.StatusOK, responses.ResponseDatabase{
			Data: bookList,
			Pagination: responses.PaginationDatabase{
				TotalRecord: counter,
				Page:        page,
				TotalPages:  int(math.Ceil(float64(counter)/10.0)) - 1,
			}})

		meterCounter.Add(c.Request.Context(), 1, metric.WithAttributes(
			attribute.String("status", strconv.Itoa(http.StatusOK)),
		))
	} else {
		responses.ResponseMessage(c, http.StatusNotFound, "There aren't books")

		meterCounter.Add(c.Request.Context(), 1, metric.WithAttributes(
			attribute.String("status", strconv.Itoa(http.StatusNotFound)),
		))
	}
}
