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

// PostBook method add new book
func (ds *DatabaseSql) PostBook(c *gin.Context) {
	book := new(models.Book)

	// creates a new span
	_, span := otel.Tracer("").Start(context.Background(), "/api/v1/book/")
	defer span.End()

	// take values from body
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

	// executes query for add new book
	_, err := ds.Db.Exec(database.ADD_NEW_BOOK,
		book.Titolo,
		book.Autore,
		book.Prezzo,
		book.Summary,
		book.Copertina,
		book.Genere,
		book.Quantita,
		book.Categoria,
		book.IdCopertina,
	)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
		return
	}

	// init a new metric
	meterCounter, err := meter.Int64Counter("post-book-counter")
	if err != nil {
		panic(err.Error())
	}

	// increment counter
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
//	@Param			id	path		int	true	"Title of the book"
//	@Success		200	{object}	models.Book
//	@Failure		404	{object}	responses.ResponseErrorJSON
//	@Failure		500	{object}	responses.ResponseErrorJSON
//	@Router			/book/{id} [get]

// GetBook method return detail of a specific book with same id
func (ds *DatabaseSql) GetBook(c *gin.Context) {
	book := new(models.Book)
	bookTitle := c.Param("title")

	// create a span
	_, span := otel.Tracer("").Start(c.Request.Context(), "/api/v1/book/id")
	defer span.End()

	// init a meter counter
	meterCounter, err := meter.Int64Counter("get-book-counter")
	if err != nil {
		panic(err.Error())
	}

	// executes query
	res, err := ds.Db.Query(database.GET_DETAIL_BOOK, bookTitle)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
		return
	}
	defer res.Close()

	// search result of the query
	for res.Next() {
		err = res.Scan(
			&book.Id,
			&book.Titolo,
			&book.Autore,
			&book.Prezzo,
			&book.Summary,
			&book.Copertina,
			&book.Genere,
			&book.Quantita,
			&book.Categoria,
			&book.IdCopertina,
		)
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
//	@Summary		Get books
//	@Schemes		http https
//	@Description	Get details for every book
//	@Tags			Book
//	@Produce		json
//	@Param			page	query		int	false	"Number of the pagination"
//	@Success		200		{object}	responses.ResponseDatabase
//	@Failure		404		{object}	responses.ResponseErrorJSON
//	@Failure		500		{object}	responses.ResponseErrorJSON
//	@Router			/book [get]

// GetBooks method that return a JSON with all books
func (ds *DatabaseSql) GetBooks(c *gin.Context) {
	var bookList *[]models.Book

	book := new(models.Book)
	counter := 0

	// take values from query params if is not null
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

	// create a span
	_, span := otel.Tracer("").Start(c.Request.Context(), "/api/v1/book")
	defer span.End()

	// init a meter counter
	meterCounter, err := meter.Int64Counter("get-books-counter")
	if err != nil {
		panic(err.Error())
	}

	// execute query
	res, err := ds.Db.Query(database.OFFSET_BOOK_PAGINATION, 10*page)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
		return
	}

	// execute result for adds books inside array
	for res.Next() {
		err = res.Scan(
			&book.Id,
			&book.Titolo,
			&book.Autore,
			&book.Prezzo,
			&book.Summary,
			&book.Copertina,
			&book.Genere,
			&book.Quantita,
			&book.Categoria,
			&book.IdCopertina,
		)
		if err != nil {
			responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
			return
		}

		*bookList = append(*bookList, *book)
	}

	// query return a number of all records
	res, err = ds.Db.Query(database.GET_COUNT_BOOKS)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
		return
	}
	defer res.Close()

	// execute query for extrapolate number of counter
	for res.Next() {
		err = res.Scan(&counter)
		if err != nil {
			responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
			return
		}
	}

	// checks length of the list
	if len(*bookList) > 0 {
		meterCounter.Add(c.Request.Context(), 1, metric.WithAttributes(
			attribute.String("status", strconv.Itoa(http.StatusOK)),
		))

		c.JSON(http.StatusOK, responses.ResponseDatabase{
			Data: *bookList,
			PaginationDatabase: responses.PaginationDatabase{
				TotalRecord: counter,
				Page:        page,
				TotalPages:  int(math.Ceil(float64(counter)/10.0)) - 1,
			}})
	} else {
		meterCounter.Add(c.Request.Context(), 1, metric.WithAttributes(
			attribute.String("status", strconv.Itoa(http.StatusNotFound)),
		))

		responses.ResponseMessage(c, http.StatusNotFound, "There aren't books")
	}
}
