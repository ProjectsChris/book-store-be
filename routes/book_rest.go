package routes

import (
	"book-store-be/database"
	"book-store-be/models"
	"book-store-be/responses"
	"context"
	"database/sql"
	"go.opentelemetry.io/otel/trace"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
)

type DatabaseSql struct {
	Db *sql.DB
}

var validate = validator.New()
var meter = otel.Meter("book-counter")

var tracer = otel.Tracer("Book Store Be")
var idBook int

// PostBook godoc
//
//	@Summary	Adds a book
//	@Schemes
//	@Description	Adds new book
//	@Tags			Book
//	@Accept			json
//	@Produce		json
//	@Param			models.Book	body		models.Book	true	"Adds new book"
//	@Success		200			{object}	responses.Response
//	@Failure		400			{object}	responses.Response
//	@Failure		500			{object}	responses.Response
//	@Router			/book/ [post]
func (ds *DatabaseSql) PostBook(c *gin.Context) {
	book := new(models.Book)

	// create a span child of "bookSpan"
	spanCtx, postBooksSpan := tracer.Start(c.Request.Context(), "/api/v1/book/")
	defer postBooksSpan.End()

	// take values from body
	if err := c.BindJSON(book); err != nil {
		responses.ErrorServerResponseJson(c, err.Error())
		return
	}

	// check validator
	validationError := validate.Struct(book)
	if validationError != nil {
		responses.ErrorServerResponseJson(c, validationError.Error())
		return
	}

	// executes query for add new book
	_, closeSpan := tracer.Start(spanCtx, "query")
	_, err := ds.Db.Exec(
		database.ADD_NEW_BOOK,
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
	closeSpan.End()
	if err != nil {
		responses.ErrorServerResponseJson(c, err.Error())
		return
	}

	// init a new metric
	meterCounter, err := meter.Int64Counter("post-book-counter")
	if err != nil {
		panic(err.Error())
	}

	// increment counter
	meterCounter.Add(c.Request.Context(), 1)
	c.JSON(http.StatusOK, responses.Response{
		Message: "added new book",
	})
}

// GetBook godoc
//
//	@Summary		Get a book
//	@Description	Get details of a book
//	@Tags			Book
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Title of the book"
//	@Success		200	{object}	models.Book
//	@Failure		404	{object}	responses.Response
//	@Failure		500	{object}	responses.Response
//	@Router			/book/{id} [get]
func (ds *DatabaseSql) GetBook(c *gin.Context) {
	book := new(models.Book)
	bookTitle := c.Param("title")

	// create a span child of "bookSpan"
	_, getBookSpan := tracer.Start(c.Request.Context(), "/api/v1/book/:id")
	defer getBookSpan.End()

	// init a meter counter
	meterCounter, err := meter.Int64Counter("get-book-counter")
	if err != nil {
		panic(err.Error())
	}

	// executes query
	res, err := ds.Db.Query(database.GET_DETAIL_BOOK, bookTitle)
	if err != nil {
		responses.ErrorServerResponseJson(c, err.Error())
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
			responses.ErrorServerResponseJson(c, err.Error())
			return
		}
	}

	// check if is an empty struct
	if *book == (models.Book{}) {
		c.JSON(http.StatusNotFound, responses.Response{Message: "book not found"})

		meterCounter.Add(c.Request.Context(), 1, metric.WithAttributes(
			attribute.String("status", strconv.Itoa(http.StatusNotFound)),
		))
	} else {
		c.JSON(http.StatusOK, gin.H{"data": book})

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
//	@Produce		json
//	@Param			page	query		int	false	"Number of the pagination"
//	@Success		200		{object}	[]models.Book
//	@Failure		404		{object}	responses.Response
//	@Failure		500		{object}	responses.Response
//	@Router			/book [get]
func (ds *DatabaseSql) GetBooks(c *gin.Context) {
	var bookList []models.Book
	book := new(models.Book)
	counter := 0

	// init context
	ctxSpan, cancelFunc := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancelFunc()

	// take values from query params if is not null
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))

	// create a parent span
	spanCtx, spanGetBooks := tracer.Start(ctxSpan, "/api/v1/book")
	defer spanGetBooks.End()

	// init a meter counter
	meterCounter, err := meter.Int64Counter("get-books-counter")
	if err != nil {
		panic(err.Error())
	}

	// create a child span
	_, spanClose := tracer.Start(spanCtx, "query OFFSET_BOOK_PAGINATION")
	res, err := ds.Db.Query(database.OFFSET_BOOK_PAGINATION, 10*page)
	if err != nil {
		responses.ErrorServerResponseJson(c, err.Error())
		spanClose.RecordError(err, trace.WithStackTrace(true))
		spanClose.End()
		return
	}

	// set status child span and close
	spanClose.SetStatus(codes.Ok, "Query OFFSET_BOOK_PAGINATION is ok")
	spanClose.End()

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
			responses.ErrorServerResponseJson(c, err.Error())
			return
		}

		bookList = append(bookList, *book)
	}

	// create a child span
	_, spanGetCountBooks := tracer.Start(spanCtx, "query GET_COUNT_BOOKS")

	// query return a number of all records
	res, err = ds.Db.Query(database.GET_COUNT_BOOKS)
	if err != nil {
		responses.ErrorServerResponseJson(c, err.Error())
		spanGetCountBooks.RecordError(err, trace.WithStackTrace(true))
		spanGetCountBooks.End()
		return
	}
	defer res.Close()

	// set status child span and close
	spanGetCountBooks.SetStatus(codes.Ok, "query GET_COUNT_BOOKS is ok")
	spanGetCountBooks.End()

	// execute query for extrapolate number of counter
	for res.Next() {
		err = res.Scan(&counter)
		if err != nil {
			responses.ErrorServerResponseJson(c, err.Error())
			return
		}
	}

	// checks length of the list
	if len(bookList) > 0 {
		meterCounter.Add(c.Request.Context(), 1, metric.WithAttributes(
			attribute.String("status", strconv.Itoa(http.StatusOK)),
		))

		c.JSON(http.StatusOK, responses.ResponsePagination{
			Data: bookList,
			PaginationDatabase: responses.PaginationDatabase{
				TotalRecord: counter,
				Page:        page,
				TotalPages:  int(math.Ceil(float64(counter)/10.0)) - 1,
			}},
		)
	} else {
		meterCounter.Add(c.Request.Context(), 1, metric.WithAttributes(
			attribute.String("status", strconv.Itoa(http.StatusNotFound)),
		))

		c.JSON(http.StatusNotFound, responses.Response{Message: "there aren't books."})
	}
}

// UpdateBook godoc
//
//	@Summary		Updates the book
//	@Description	Update all details of a book
//	@Description	The title cannot be more than 255 characters.
//	@Tags			Book
//	@Produce		json
//	@Param			id		path		string		true	"Id"
//	@Param			title	body		models.Book	true	"Book"
//	@Success		200		{object}	responses.Response
//	@Failure		404		{object}	responses.Response
//	@Failure		500		{object}	responses.Response
//	@Router			/book/{id} [put]
//
// UpdateBook method for update the book
func (ds *DatabaseSql) UpdateBook(c *gin.Context) {
	idBook, _ = strconv.Atoi(c.Param("id"))
	var book = new(models.Book)

	// init context
	ctx, cancelFunc := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancelFunc()

	// create a parent span
	ctxSpan, spanUpdate := tracer.Start(ctx, "/api/v1/book/{id}")
	defer spanUpdate.End()

	// take values from JSON
	if err := c.BindJSON(book); err != nil {
		responses.ErrorServerResponseJson(c, err.Error())
		return
	}

	// validation of fields
	if book.Titolo != "" {
		if errValidator := validate.StructPartial(book, "Titolo"); errValidator != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Message: errValidator.Error()})
			return
		}
	}

	if book.Autore != "" {
		if errValidator := validate.StructPartial(book, "Autore"); errValidator != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Message: errValidator.Error()})
			return
		}
	}

	if book.Prezzo != 0 {
		if errValidator := validate.StructPartial(book, "Prezzo"); errValidator != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Message: errValidator.Error()})
			return
		}
	}

	if book.Summary != "" {
		if errValidator := validate.StructPartial(book, "Summary"); errValidator != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Message: errValidator.Error()})
			return
		}
	}

	if book.Copertina != "" {
		if errValidator := validate.StructPartial(book, "Copertina"); errValidator != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Message: errValidator.Error()})
			return
		}
	}

	if book.Genere != "" {
		if errValidator := validate.StructPartial(book, "Genere"); errValidator != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Message: errValidator.Error()})
			return
		}
	}

	if book.Quantita != 0 {
		if errValidator := validate.StructPartial(book, "Quantita"); errValidator != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Message: errValidator.Error()})
			return
		}
	}

	if book.Categoria != "" {
		if errValidator := validate.StructPartial(book, "Categoria"); errValidator != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Message: errValidator.Error()})
			return
		}
	}

	if book.IdCopertina != 0 {
		if errValidator := validate.StructPartial(book, "IdCopertina"); errValidator != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Message: errValidator.Error()})
			return
		}
	}

	// create a child span
	_, spanUpdateQuery := tracer.Start(ctxSpan, "query UPDATE_BOOK")

	// exec query
	_, err := ds.Db.Exec(
		database.UPDATE_BOOK,
		book.Titolo,
		book.Autore,
		book.Prezzo,
		book.Summary,
		book.Copertina,
		book.Genere,
		book.Quantita,
		book.Categoria,
		book.IdCopertina,
		idBook,
	)
	if err != nil {
		responses.ErrorServerResponseJson(c, err.Error())
		spanUpdateQuery.RecordError(err, trace.WithStackTrace(true))
		spanUpdateQuery.SetStatus(codes.Error, "")
		spanUpdateQuery.End()
		return
	}

	// response of success
	spanUpdateQuery.SetStatus(codes.Ok, "success")
	spanUpdateQuery.End()

	c.JSON(http.StatusOK, responses.Response{Message: "book updated"})
}

// DeleteBook godoc
//
//	@Summary		Delete a book
//	@Description	Delete a book with same ID
//	@Tags			Book
//	@Produce		json
//	@Param			id	path		int	true	"Id"
//	@Success		200	{object}	models.Book
//	@Failure		404	{object}	responses.Response
//	@Failure		500	{object}	responses.Response
//	@Router			/book/{id} [delete]
//
// DeleteBook method for delete a book
func (ds *DatabaseSql) DeleteBook(c *gin.Context) {
	idBook, _ = strconv.Atoi(c.Param("id"))

	// exec query
	res, err := ds.Db.Exec(database.DELETE_BOOK, idBook)
	if err != nil {
		responses.ErrorServerResponseJson(c, err.Error())
		return
	}

	status, err := res.RowsAffected()
	if err != nil {
		responses.ErrorServerResponseJson(c, err.Error())
		return
	}

	if status > 0 {
		c.JSON(http.StatusOK, responses.Response{Message: "book deleted"})
	} else {
		c.JSON(http.StatusNotFound, responses.Response{Message: "book not found or already deleted."})
	}
}
