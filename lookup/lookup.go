package lookup

import (
	"errors"

	"github.com/domgoodwin/bookscan/book"
)

func LookupISBN(isbn string) (*book.Book, error) {
	// Look through different ISBN lookup provider services
	book, err := OpenLibraryLookup(isbn)
	if err != nil {
		return nil, err
	}
	if book != nil {
		return book, nil
	}

	return nil, errors.New("not found")
}
