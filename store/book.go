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

func (s *BookStorer) Setup() {
	s.books = make(map[string]map[string]*items.Book)
	// for _, id := range []string{"4f311bbe86ce4dd4bdae93fa1206328f", "7b929608b3dc460c98d804e20de882c2"} {
	// 	s.LoadBooksFromNotion(context.Background(), id)
	// }
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

// func (s *BookStorer) LoadBooksFromNotion(ctx context.Context, databaseID string) error {
// 	if os.Getenv("NOTION_CACHE") == "false" {
// 		logrus.Info("Notion cache disabled")
// 		return nil
// 	}
// 	books, databaseID, err := notion.GetBookPagesFromDatabase(ctx, databaseID)
// 	if err != nil {
// 		return err
// 	}
// 	for _, b := range books {
// 		s.StoreItem(databaseID, b)
// 	}
// 	logrus.Infof("Loaded %v books from Notion:%v", len(books), databaseID)
// 	s.databaseID = databaseID
// 	return nil
// }

func removeQuotes(in string) string {
	return in[1 : len(in)-1]
}
