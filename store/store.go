package store

import (
	"context"

	"github.com/domgoodwin/bookscan/notion"
	"github.com/sirupsen/logrus"
)

var BookStore *BookStorer
var RecordStore *RecordStorer

func SetupStore(ctx context.Context) error {
	books, records, err := notion.GetAllPagesFromAllDatabases(ctx)
	if err != nil {
		logrus.Error(err)
		return err
	}
	BookStore = &BookStorer{}
	BookStore.Setup(books)
	RecordStore = &RecordStorer{}
	RecordStore.Setup(records)
	return nil
}
