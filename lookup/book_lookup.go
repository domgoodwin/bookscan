package lookup

import (
	"errors"

	"github.com/domgoodwin/bookscan/items"
	"github.com/sirupsen/logrus"
)

var bookDatasources = map[BookDataSource]bool{
	openLibraryDataStore{}:  true,
	googlebooksDataSource{}: true,
}

type BookDataSource interface {
	Name() string
	LookupISBN(isbn string) (*items.Book, error)
}

func LookupISBN(isbn string) (*items.Book, error) {
	logrus.Info("looking up isbn ", map[string]string{
		"isbn": isbn,
	})
	for d, enabled := range bookDatasources {
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
