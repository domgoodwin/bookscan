package store

import (
	"github.com/domgoodwin/bookscan/items"
	"github.com/sirupsen/logrus"
)

const csvFilePath = "./books.csv"

type BookStorer struct {
	// map of database ids to map of book isbns to books
	books      map[string]map[string]*items.Book
	databaseID string
}

func (s *BookStorer) Setup(in map[string]map[string]*items.Book) {
	s.books = in
}

func (s *BookStorer) StoreItem(databaseID string, b *items.Book) bool {
	_, databaseFound := s.books[databaseID]
	if !databaseFound {
		s.books[databaseID] = make(map[string]*items.Book)
	}
	_, bookFound := s.books[databaseID][b.ISBN]
	if !bookFound {
		s.books[databaseID][b.ISBN] = b
	}
	return bookFound
}

func (s *BookStorer) CheckIfItemInCache(databaseID, isbn string) (*items.Book, bool) {
	_, found := s.books[databaseID]
	if !found {
		s.books[databaseID] = make(map[string]*items.Book)
	}
	book, found := s.books[databaseID][isbn]
	return book, found
}

func (s *BookStorer) ClearCache(databaseID string) int {
	logrus.Info("clearing cache")
	oldLength := len(s.books[databaseID])
	s.books[databaseID] = make(map[string]*items.Book)
	return oldLength
}

func (s *BookStorer) Length() int {
	return len(s.books)
}

func (s *BookStorer) DatabaseID() string {
	return s.databaseID
}

func removeQuotes(in string) string {
	return in[1 : len(in)-1]
}
