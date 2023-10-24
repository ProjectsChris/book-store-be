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
//	@Success		200			{object}	models.Book
//	@Failure		400			{object}	responses.ResponseErrorJSON
//	@Failure		500			{object}	responses.ResponseErrorJSON
//	@Router			/book/ [post]
func (ds *DatabaseSql) PostBook(c *gin.Context) {
	book := new(models.Book)

	// create a span child of "bookSpan"
	spanCtx, postBooksSpan := tracer.Start(c.Request.Context(), "/api/v1/book/")
	defer postBooksSpan.End()

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
	_, closeSpan := tracer.Start(spanCtx, "query")
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
	closeSpan.End()
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
//	@Summary		Get a book
//	@Description	Get details of a book
//	@Tags			Book
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Title of the book"
//	@Success		200	{object}	models.Book
//	@Failure		404	{object}	responses.ResponseErrorJSON
//	@Failure		500	{object}	responses.ResponseErrorJSON
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
//	@Summary	Get books
//	@Schemes
//	@Description	Get details for every book
//	@Tags			Book
//	@Produce		json
//	@Param			page	query		int	false	"Number of the pagination"
//	@Success		200		{object}	responses.ResponseDatabase
//	@Failure		404		{object}	responses.ResponseErrorJSON
//	@Failure		500		{object}	responses.ResponseErrorJSON
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
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
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
			responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
			return
		}

		bookList = append(bookList, *book)
	}

	// create a child span
	_, spanGetCountBooks := tracer.Start(spanCtx, "query GET_COUNT_BOOKS")

	// query return a number of all records
	res, err = ds.Db.Query(database.GET_COUNT_BOOKS)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
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
			responses.ResponseMessage(c, http.StatusInternalServerError, "error: "+err.Error())
			return
		}
	}

	// checks length of the list
	if len(bookList) > 0 {
		meterCounter.Add(c.Request.Context(), 1, metric.WithAttributes(
			attribute.String("status", strconv.Itoa(http.StatusOK)),
		))

		c.JSON(http.StatusOK, responses.ResponseDatabase{
			Data: bookList,
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

// UpdateTitleBook godoc
//
//	@Summary		Update title
//	@Description	Update title of a book.
//	@Description	The title cannot be more than 255 characters.
//	@Tags			Book
//	@Produce		json
//	@Param			id		path		string	true	"Id"
//	@Param			title	body		object	true	"Title"
//	@Success		200		{object}	responses.ResponseErrorJSON
//	@Failure		404		{object}	responses.ResponseErrorJSON
//	@Failure		500		{object}	responses.ResponseErrorJSON
//	@Router			/book/titolo/{id} [put]
//
// UpdateTitleBook method for update "titolo"
func (ds *DatabaseSql) UpdateTitleBook(c *gin.Context) {
	idBook, _ = strconv.Atoi(c.Param("id"))
	var book = new(models.Book)

	// init context
	ctx, cancelFunc := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancelFunc()

	// create a parent span
	ctxSpan, spanUpdate := tracer.Start(ctx, "/api/v1/titolo/{id}")
	defer spanUpdate.End()

	// take values from JSON
	if err := c.BindJSON(book); err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	// validation of field
	validationError := validate.StructPartial(book, "Titolo")
	if validationError != nil {
		responses.ResponseMessage(c, http.StatusBadRequest, "error: "+validationError.Error())
		return
	}

	// create a child span
	_, spanUpdateQuery := tracer.Start(ctxSpan, "query UPDATE_TITLE_BOOK")

	// exec query
	_, err := ds.Db.Exec(database.UPDATE_TITLE_BOOK, idBook, book.Titolo)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		spanUpdateQuery.RecordError(err, trace.WithStackTrace(true))
		spanUpdateQuery.SetStatus(codes.Error, "")
		spanUpdateQuery.End()
		return
	}

	// response of success
	responses.ResponseMessage(c, http.StatusOK, "value of 'titolo' is updated")
	spanUpdateQuery.SetStatus(codes.Ok, "success")
	spanUpdateQuery.End()
}

// UpdateWriterBook godoc
//
//	@Summary		Update writer
//	@Description	Update name of the writer.
//	@Description	The name of the writer can't be more than 64 characters.
//	@Tags			Book
//	@Produce		json
//	@Param			id		path		string	true	"Id"
//	@Param			writer	body		object	true	"Writer"
//	@Success		200		{object}	responses.ResponseErrorJSON
//	@Failure		404		{object}	responses.ResponseErrorJSON
//	@Failure		500		{object}	responses.ResponseErrorJSON
//	@Router			/book/autore/{id} [put]
//
// UpdateWriterBook method for update "autore"
func (ds *DatabaseSql) UpdateWriterBook(c *gin.Context) {
	idBook, _ = strconv.Atoi(c.Param("id"))
	var book = new(models.Book)

	// init context
	ctx, cancelFunc := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancelFunc()

	// create a parent span
	ctxSpan, spanUpdate := tracer.Start(ctx, "/api/v1/autore/{id}")
	defer spanUpdate.End()

	// take values from JSON
	if err := c.BindJSON(book); err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	// validation of field
	validationError := validate.StructPartial(book, "Autore")
	if validationError != nil {
		responses.ResponseMessage(c, http.StatusBadRequest, "error: "+validationError.Error())
		return
	}

	// create a child span
	_, spanUpdateQuery := tracer.Start(ctxSpan, "query UPDATE_WRITER_BOOK")

	// exec query
	_, err := ds.Db.Exec(database.UPDATE_WRITER_BOOK, idBook, book.Autore)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		spanUpdateQuery.RecordError(err, trace.WithStackTrace(true))
		spanUpdateQuery.SetStatus(codes.Error, "")
		spanUpdateQuery.End()
		return
	}

	// response of success
	c.JSON(http.StatusOK, gin.H{
		"success": "value of 'autore' is updated",
	})
	spanUpdateQuery.SetStatus(codes.Ok, "success")
	spanUpdateQuery.End()
}

// UpdatePriceBook godoc
//
//	@Summary		Update price
//	@Description	Update the price of a book.
//	@Description	The type of the price is a float (e.g. 15.90).
//	@Tags			Book
//	@Produce		json
//	@Param			id		path		string	true	"Id"
//	@Param			price	body		object	true	"Price"
//	@Success		200		{object}	responses.ResponseErrorJSON
//	@Failure		404		{object}	responses.ResponseErrorJSON
//	@Failure		500		{object}	responses.ResponseErrorJSON
//	@Router			/book/prezzo/{id} [put]
//
// UpdatePriceBook method for update price
func (ds *DatabaseSql) UpdatePriceBook(c *gin.Context) {
	idBook, _ = strconv.Atoi(c.Param("id"))
	var book = new(models.Book)

	// init context
	ctx, cancelFunc := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancelFunc()

	// create a parent span
	ctxSpan, spanUpdate := tracer.Start(ctx, "/api/v1/prezzo/{id}")
	defer spanUpdate.End()

	// take values from JSON
	if err := c.BindJSON(book); err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	// validation of field
	validationError := validate.StructPartial(book, "Prezzo")
	if validationError != nil {
		responses.ResponseMessage(c, http.StatusBadRequest, "error: "+validationError.Error())
		return
	}

	// create a child span
	_, spanUpdateQuery := tracer.Start(ctxSpan, "query UPDATE_PRICE_BOOK")

	// exec query
	_, err := ds.Db.Exec(database.UPDATE_PRICE_BOOK, idBook, book.Prezzo)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		spanUpdateQuery.RecordError(err, trace.WithStackTrace(true))
		spanUpdateQuery.SetStatus(codes.Error, "")
		spanUpdateQuery.End()
		return
	}

	// response of success
	c.JSON(http.StatusOK, gin.H{
		"success": "value of 'prezzo' is updated",
	})
	spanUpdateQuery.SetStatus(codes.Ok, "success")
	spanUpdateQuery.End()
}

// UpdateSummaryBook godoc
//
//	@Summary		Update summary
//	@Description	Update the summary of a book.
//	@Description	The summary of the book cannot be less than 512 characters.
//	@Tags			Book
//	@Produce		json
//	@Param			id		path		string	true	"Id"
//	@Param			summary	body		object	true	"Summary"
//	@Success		200		{object}	responses.ResponseErrorJSON
//	@Failure		404		{object}	responses.ResponseErrorJSON
//	@Failure		500		{object}	responses.ResponseErrorJSON
//	@Router			/book/summary/{id} [put]
//
// UpdateSummaryBook method for update summary
func (ds *DatabaseSql) UpdateSummaryBook(c *gin.Context) {
	idBook, _ = strconv.Atoi(c.Param("id"))
	var book = new(models.Book)

	// init context
	ctx, cancelFunc := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancelFunc()

	// create a parent span
	ctxSpan, spanUpdate := tracer.Start(ctx, "/api/v1/summary/{id}")
	defer spanUpdate.End()

	// take values from JSON
	if err := c.BindJSON(book); err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	// validation of field
	validationError := validate.StructPartial(book, "Summary")
	if validationError != nil {
		responses.ResponseMessage(c, http.StatusBadRequest, "error: "+validationError.Error())
		return
	}

	// create a child span
	_, spanUpdateQuery := tracer.Start(ctxSpan, "query UPDATE_SUMMARY_BOOK")

	// exec query
	_, err := ds.Db.Exec(database.UPDATE_SUMMARY_BOOK, idBook, book.Summary)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		spanUpdateQuery.RecordError(err, trace.WithStackTrace(true))
		spanUpdateQuery.SetStatus(codes.Error, "")
		spanUpdateQuery.End()
		return
	}

	// response of success
	c.JSON(http.StatusOK, gin.H{
		"success": "value of 'summary' is updated",
	})
	spanUpdateQuery.SetStatus(codes.Ok, "success")
	spanUpdateQuery.End()
}

// UpdateCoverBook godoc
//
//	@Summary		Update cover
//	@Description	Update the type cover of a book.
//	@Description	There are two types of cover: Hard Cover or Flexible Cover.
//	@Tags			Book
//	@Produce		json
//	@Param			id		path		string	true	"Id"
//	@Param			cover	body		object	true	"Cover"
//	@Success		200		{object}	responses.ResponseErrorJSON
//	@Failure		404		{object}	responses.ResponseErrorJSON
//	@Failure		500		{object}	responses.ResponseErrorJSON
//	@Router			/book/copertina/{id} [put]
//
// UpdateCoverBook method for update cover
func (ds *DatabaseSql) UpdateCoverBook(c *gin.Context) {
	idBook, _ = strconv.Atoi(c.Param("id"))
	var book = new(models.Book)

	// init context
	ctx, cancelFunc := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancelFunc()

	// create a parent span
	ctxSpan, spanUpdate := tracer.Start(ctx, "/api/v1/copertina/{id}")
	defer spanUpdate.End()

	// take values from JSON
	if err := c.BindJSON(book); err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	// validation of field
	validationError := validate.StructPartial(book, "Copertina")
	if validationError != nil {
		responses.ResponseMessage(c, http.StatusBadRequest, "error: "+validationError.Error())
		return
	}

	// create a child span
	_, spanUpdateQuery := tracer.Start(ctxSpan, "query UPDATE_COVER_BOOK")

	// exec query
	_, err := ds.Db.Exec(database.UPDATE_COVER_BOOK, idBook, book.Copertina)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		spanUpdateQuery.RecordError(err, trace.WithStackTrace(true))
		spanUpdateQuery.SetStatus(codes.Error, "")
		spanUpdateQuery.End()
		return
	}

	// response of success
	c.JSON(http.StatusOK, gin.H{
		"success": "value of 'copertina' is updated",
	})
	spanUpdateQuery.SetStatus(codes.Ok, "success")
	spanUpdateQuery.End()
}

// UpdateGenreBook godoc
//
//	@Summary		Update genre
//	@Description	Update the genre of a book.
//	@Description	Here a list of all genre supported: Action, Adventure, Business, Cookbooks, Drama, Detective, Fantasy, Fiction, History, Horror, Romance, Psychology, Science Fiction, Short Stories, Thriller.
//	@Tags			Book
//	@Produce		json
//	@Param			id		path		string	true	"Id"
//	@Param			genre	body		object	true	"Genre"
//	@Success		200		{object}	responses.ResponseErrorJSON
//	@Failure		404		{object}	responses.ResponseErrorJSON
//	@Failure		500		{object}	responses.ResponseErrorJSON
//	@Router			/book/genere/{id} [put]
//
// UpdateGenreBook method for update genre
func (ds *DatabaseSql) UpdateGenreBook(c *gin.Context) {
	idBook, _ = strconv.Atoi(c.Param("id"))
	var book = new(models.Book)

	// init context
	ctx, cancelFunc := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancelFunc()

	// create a parent span
	ctxSpan, spanUpdate := tracer.Start(ctx, "/api/v1/genere/{id}")
	defer spanUpdate.End()

	// take values from JSON
	if err := c.BindJSON(book); err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	// validation of field
	validationError := validate.StructPartial(book, "Genere")
	if validationError != nil {
		responses.ResponseMessage(c, http.StatusBadRequest, "error: "+validationError.Error())
		return
	}

	// create a child span
	_, spanUpdateQuery := tracer.Start(ctxSpan, "query UPDATE_GENRE_BOOK")

	// exec query
	_, err := ds.Db.Exec(database.UPDATE_GENRE_BOOK, idBook, book.Genere)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		spanUpdateQuery.RecordError(err, trace.WithStackTrace(true))
		spanUpdateQuery.SetStatus(codes.Error, "")
		spanUpdateQuery.End()
		return
	}

	// response of success
	c.JSON(http.StatusOK, gin.H{
		"success": "value of 'genere' is updated",
	})
	spanUpdateQuery.SetStatus(codes.Ok, "success")
	spanUpdateQuery.End()
}

// UpdateQuantityBook godoc
//
//	@Summary		Update quantity
//	@Description	Update the quantity of a book
//	@Description	You can choose the quantity from 1 to 5
//	@Tags			Book
//	@Produce		json
//	@Param			id			path		string	true	"Id"
//	@Param			quantity	body		object	true	"Quantity"
//	@Success		200			{object}	responses.ResponseErrorJSON
//	@Failure		404			{object}	responses.ResponseErrorJSON
//	@Failure		500			{object}	responses.ResponseErrorJSON
//	@Router			/book/quantita/{id} [put]
//
// UpdateQuantityBook method for update quantity
func (ds *DatabaseSql) UpdateQuantityBook(c *gin.Context) {
	idBook, _ = strconv.Atoi(c.Param("id"))
	var book = new(models.Book)

	// init context
	ctx, cancelFunc := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancelFunc()

	// create a parent span
	ctxSpan, spanUpdate := tracer.Start(ctx, "/api/v1/quantita/{id}")
	defer spanUpdate.End()

	// take values from JSON
	if err := c.BindJSON(book); err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	// validation of field
	validationError := validate.StructPartial(book, "Quantita")
	if validationError != nil {
		responses.ResponseMessage(c, http.StatusBadRequest, "error: "+validationError.Error())
		return
	}

	// create a child span
	_, spanUpdateQuery := tracer.Start(ctxSpan, "query UPDATE_QUANTITY_BOOK")

	// exec query
	_, err := ds.Db.Exec(database.UPDATE_QUANTITY_BOOK, idBook, book.Quantita)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		spanUpdateQuery.RecordError(err, trace.WithStackTrace(true))
		spanUpdateQuery.SetStatus(codes.Error, "")
		spanUpdateQuery.End()
		return
	}

	// response of success
	c.JSON(http.StatusOK, gin.H{
		"success": "value of 'quantita' is updated",
	})
	spanUpdateQuery.SetStatus(codes.Ok, "success")
	spanUpdateQuery.End()
}

// UpdateCategoryBook godoc
//
//	@Summary		Update category
//	@Description	Update the category of a book
//	@Description	The types of categories are: Best Seller, New Releases, and Best Offers.
//	@Tags			Book
//	@Produce		json
//	@Param			id			path		string	true	"Id"
//	@Param			category	body		object	true	"Category"
//	@Success		200			{object}	responses.ResponseErrorJSON
//	@Failure		404			{object}	responses.ResponseErrorJSON
//	@Failure		500			{object}	responses.ResponseErrorJSON
//	@Router			/book/categoria/{id} [put]
//
// UpdateCategoryBook method for update category
func (ds *DatabaseSql) UpdateCategoryBook(c *gin.Context) {
	idBook, _ = strconv.Atoi(c.Param("id"))
	var book = new(models.Book)

	// init context
	ctx, cancelFunc := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancelFunc()

	// create a parent span
	ctxSpan, spanUpdate := tracer.Start(ctx, "/api/v1/categoria/{id}")
	defer spanUpdate.End()

	// take values from JSON
	if err := c.BindJSON(book); err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	// validation of field
	validationError := validate.StructPartial(book, "Categoria")
	if validationError != nil {
		responses.ResponseMessage(c, http.StatusBadRequest, "error: "+validationError.Error())
		return
	}

	// create a child span
	_, spanUpdateQuery := tracer.Start(ctxSpan, "query UPDATE_CATEGORY_BOOK")

	// exec query
	_, err := ds.Db.Exec(database.UPDATE_CATEGORY_BOOK, idBook, book.Categoria)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		spanUpdateQuery.RecordError(err, trace.WithStackTrace(true))
		spanUpdateQuery.SetStatus(codes.Error, "")
		spanUpdateQuery.End()
		return
	}

	// response of success
	c.JSON(http.StatusOK, gin.H{
		"success": "value of 'categoria' is updated",
	})
	spanUpdateQuery.SetStatus(codes.Ok, "success")
	spanUpdateQuery.End()
}

// UpdateIdCoverBook godoc
//
//	@Summary		Update id cover
//	@Description	Update the id cover of a book
//	@Tags			Book
//	@Produce		json
//	@Param			id		path		string	true	"Id"
//	@Param			idcover	body		object	true	"Id Cover"
//	@Success		200		{object}	responses.ResponseErrorJSON
//	@Failure		404		{object}	responses.ResponseErrorJSON
//	@Failure		500		{object}	responses.ResponseErrorJSON
//	@Router			/book/id-copertina/{id} [put]
//
// UpdateIdCoverBook method for update id cover
func (ds *DatabaseSql) UpdateIdCoverBook(c *gin.Context) {
	idBook, _ = strconv.Atoi(c.Param("id"))
	var book = new(models.Book)

	// init context
	ctx, cancelFunc := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancelFunc()

	// create a parent span
	ctxSpan, spanUpdate := tracer.Start(ctx, "/api/v1/id-cover/{id}")
	defer spanUpdate.End()

	// take values from JSON
	if err := c.BindJSON(book); err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	// validation of field
	validationError := validate.StructPartial(book, "IdCopertina")
	if validationError != nil {
		responses.ResponseMessage(c, http.StatusBadRequest, "error: "+validationError.Error())
		return
	}

	// create a child span
	_, spanUpdateQuery := tracer.Start(ctxSpan, "query UPDATE_ID_COVER_BOOK")

	// exec query
	_, err := ds.Db.Exec(database.UPDATE_ID_COVER_BOOK, idBook, book.IdCopertina)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		spanUpdateQuery.RecordError(err, trace.WithStackTrace(true))
		spanUpdateQuery.SetStatus(codes.Error, "")
		spanUpdateQuery.End()
		return
	}

	// response of success
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "value of 'id_copertina' is updated",
	})
	spanUpdateQuery.SetStatus(codes.Ok, "success")
	spanUpdateQuery.End()
}

// DeleteBook godoc
//
//	@Summary		Delete a book
//	@Description	Delete a book with same ID
//	@Tags			Book
//	@Produce		json
//	@Param			id	path		int	true	"Id"
//	@Success		200	{object}	responses.ResponseErrorJSON
//	@Failure		404	{object}	responses.ResponseErrorJSON
//	@Failure		500	{object}	responses.ResponseErrorJSON
//	@Router			/book/{id} [delete]
//
// DeleteBook method for delete a book
func (ds *DatabaseSql) DeleteBook(c *gin.Context) {
	idBook, _ = strconv.Atoi(c.Param("id"))

	// exec query
	res, err := ds.Db.Exec(database.DELETE_BOOK, idBook)
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	status, err := res.RowsAffected()
	if err != nil {
		responses.ResponseMessage(c, http.StatusInternalServerError, err.Error())
		return
	}

	if status > 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": "book deleted",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error": "book not found or already deleted.",
		})
	}

}
