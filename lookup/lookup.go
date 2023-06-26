package lookup

import (
	"errors"

	"github.com/domgoodwin/bookscan/book"
	"github.com/sirupsen/logrus"
)

func LookupISBN(isbn string) (*book.Book, error) {
	logrus.Info("looking up isbn ", map[string]string{
		"isbn": isbn,
	})
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
