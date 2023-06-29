package store

import (
	"bufio"
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/domgoodwin/bookscan/book"
	"github.com/domgoodwin/bookscan/notion"
	"github.com/sirupsen/logrus"
)

const csvFilePath = "./books.csv"

type Store struct {
	books      map[string]*book.Book
	databaseID string
}

var store *Store

func SetupStore() *Store {
	if store == nil {
		store = &Store{
			books: make(map[string]*book.Book),
		}
	}
	store.LoadBooksFromCSV()
	store.LoadBooksFromNotion(context.Background(), "")
	return store
}

func (s *Store) StoreBook(b *book.Book) bool {
	_, found := s.books[b.ISBN]
	if !found {
		s.books[b.ISBN] = b
	}
	return found
}

func (s *Store) CheckIfBookInCache(isbn string) (*book.Book, bool) {
	book, found := s.books[isbn]
	return book, found
}

func (s *Store) ClearCache() int {
	logrus.Info("clearing cache")
	oldLength := len(s.books)
	s.books = make(map[string]*book.Book)
	return oldLength
}

func (s *Store) Length() int {
	return len(s.books)
}

func (s *Store) DatabaseID() string {
	return s.databaseID
}

func (s *Store) LoadBooksFromCSV() {
	if os.Getenv("CSV_CACHE") == "false" {
		logrus.Info("CSV cache disabled")
		return
	}
	count := 0
	file, err := os.Open(csvFilePath)
	if err != nil {
		logrus.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count++
		logrus.Debug(scanner.Text())
		book := CSVBookToBook(scanner.Text())
		s.StoreBook(book)
	}

	if err := scanner.Err(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Loaded %v books from CSV", count)
}

func (s *Store) LoadBooksFromNotion(ctx context.Context, databaseID string) error {
	if os.Getenv("NOTION_CACHE") == "false" {
		logrus.Info("Notion cache disabled")
		return nil
	}
	books, databaseID, err := notion.GetPagesFromDatabase(ctx, databaseID)
	if err != nil {
		return err
	}
	for _, b := range books {
		s.StoreBook(b)
	}
	logrus.Infof("Loaded %v books from Notion", len(books))
	s.databaseID = databaseID
	return nil
}

func CSVBookToBook(line string) *book.Book {
	parts := strings.Split(line, ",")
	if len(parts) != 6 {
		logrus.Warn("Line doesn't have 6 parts", map[string]string{"line": line})
		return nil
	}

	pages, err := strconv.Atoi(removeQuotes(parts[3]))
	if err != nil {
		logrus.Error(err)
	}
	return &book.Book{
		Title:    removeQuotes(parts[0]),
		Authors:  strings.Split(removeQuotes(parts[1]), ";"),
		ISBN:     removeQuotes(parts[2]),
		Pages:    pages,
		Link:     removeQuotes(parts[4]),
		CoverURL: removeQuotes(parts[5]),
	}
}

func removeQuotes(in string) string {
	return in[1 : len(in)-1]
}
