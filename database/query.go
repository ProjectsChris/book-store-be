package database

const (
	ADD_NEW_BOOK           = `INSERT INTO books (title, writer, price, summary, cover, genre, quantity, category, IdCover) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
	GET_DETAIL_BOOK        = `SELECT * FROM books WHERE id = $1;`
	OFFSET_BOOK_PAGINATION = `SELECT * FROM books ORDER BY id DESC LIMIT 10 OFFSET $1;`
	GET_COUNT_BOOKS        = `SELECT COUNT(*) FROM books;`
)
