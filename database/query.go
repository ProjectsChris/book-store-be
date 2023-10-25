package database

const (
	ADD_NEW_BOOK           = `INSERT INTO books (title, writer, price, summary, cover, genre, quantity, category, IdCover) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
	GET_DETAIL_BOOK        = `SELECT * FROM books WHERE id = $1;`
	OFFSET_BOOK_PAGINATION = `SELECT * FROM books ORDER BY id DESC LIMIT 10 OFFSET $1;`
	GET_COUNT_BOOKS        = `SELECT COUNT(*) FROM books;`
	DELETE_BOOK            = `DELETE FROM books WHERE id = $1;`
	UPDATE_BOOK            = `UPDATE books SET title = COALESCE(NULLIF($1, ''), title), writer = COALESCE(NULLIF($2, ''), writer), price = COALESCE(NULLIF($3, 0), price), summary = COALESCE(NULLIF($4, ''), summary), cover = COALESCE(NULLIF($5, ''), cover), genre = COALESCE(NULLIF($6, ''), genre), quantity = COALESCE(NULLIF($7, 0), quantity), category = COALESCE(NULLIF($8, ''), category), idcover = COALESCE(NULLIF($9, 0), idcover) WHERE id = $10;`
)
