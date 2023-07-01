package lookup

import (
	"errors"

	"github.com/domgoodwin/bookscan/book"
	"github.com/sirupsen/logrus"
)

var datasources = map[BookDataSource]bool{
	openLibraryDataStore{}:  false,
	googlebooksDataSource{}: true,
}

type BookDataSource interface {
	Name() string
	LookupISBN(isbn string) (*book.Book, error)
}

func LookupISBN(isbn string) (*book.Book, error) {
	logrus.Info("looking up isbn ", map[string]string{
		"isbn": isbn,
	})
	for d, enabled := range datasources {
		if !enabled {
			logrus.Info("skipping data source: ", d.Name())
			continue
		}
		book, err := d.LookupISBN(isbn)
		if err != nil {
			return nil, err
		}
		if book != nil {
			return book, nil
		}
	}
	return nil, errors.New("not found")
}
