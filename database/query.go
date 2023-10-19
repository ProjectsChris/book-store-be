package database

const (
	ADD_NEW_BOOK           = `INSERT INTO books (title, writer, price, summary, cover, genre, quantity, category, IdCover) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
	GET_DETAIL_BOOK        = `SELECT * FROM books WHERE id = $1;`
	OFFSET_BOOK_PAGINATION = `SELECT * FROM books ORDER BY id DESC LIMIT 10 OFFSET $1;`
	GET_COUNT_BOOKS        = `SELECT COUNT(*) FROM books;`
	UPDATE_TITLE_BOOK      = `UPDATE books SET title = $2 WHERE id = $1;`
	UPDATE_WRITER_BOOK     = `UPDATE books SET writer = $2 WHERE id = $1;`
	UPDATE_PRICE_BOOK      = `UPDATE books SET price = $2 WHERE id = $1;`
	UPDATE_SUMMARY_BOOK    = `UPDATE books SET summary = $2 WHERE id = $1;`
	UPDATE_COVER_BOOK      = `UPDATE books SET cover = $2 WHERE id = $1;`
	UPDATE_GENRE_BOOK      = `UPDATE books SET genre = $2 WHERE id = $1;`
	UPDATE_QUANTITY_BOOK   = `UPDATE books SET quantity = $2 WHERE id = $1;`
	UPDATE_CATEGORY_BOOK   = `UPDATE books SET category = $2 WHERE id = $1;`
	UPDATE_ID_COVER_BOOK   = `UPDATE books SET idcover = $2 WHERE id = $1;`
)
