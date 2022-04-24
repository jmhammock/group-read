package models

type Book struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Location string `json:"location"`
}

type Books []Book

type BookModel struct {
	db DB
}

func NewBookModel(db DB) *BookModel {
	return &BookModel{
		db: db,
	}
}

func (b BookModel) Get(limit, offset uint16) (*Books, error) {
	q := `SELECT *
		FROM books
		LIMIT $1
		OFFSET $2;`

	rows, err := b.db.Query(q, limit, offset)
	if err != nil {
		return nil, err
	}

	var books Books
	for rows.Next() {
		var book Book
		err := rows.Scan(&book)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return &books, nil
}

func (b BookModel) GetById(id string) (*Book, error) {
	q := `SELECT *
		FROM books
		WHERE id = $1;`

	var book *Book
	err := b.db.QueryRow(q, id).Scan(&book)

	return book, err
}

func (b BookModel) Insert(book *Book) (string, error) {
	q := `INSERT INTO books (title, author, location)
		VALUES($1, $2, $3)
		RETURNING id;`

	var id string
	args := []interface{}{
		book.Title,
		book.Author,
		book.Location,
	}
	err := b.db.QueryRow(q, args...).Scan(&id)

	return id, err
}

func (b BookModel) Delete(id string) error {
	q := `DELTE FROM books
		WHERE id = $1;`

	_, err := b.db.Execute(q, id)
	return err
}
