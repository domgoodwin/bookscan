package store

import (
	"bufio"
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/domgoodwin/bookscan/items"
	"github.com/domgoodwin/bookscan/notion"
	"github.com/sirupsen/logrus"
)

const csvFilePath = "./books.csv"

type BookStorer struct {
	books      map[string]*items.Book
	databaseID string
}

func (s *BookStorer) Setup() {
	s.books = make(map[string]*items.Book)
	s.LoadBooksFromCSV()
	s.LoadBooksFromNotion(context.Background(), "")
}

func (s *BookStorer) StoreItem(b *items.Book) bool {
	_, found := s.books[b.ISBN]
	if !found {
		s.books[b.ISBN] = b
	}
	return found
}

func (s *BookStorer) CheckIfItemInCache(isbn string) (*items.Book, bool) {
	book, found := s.books[isbn]
	return book, found
}

func (s *BookStorer) ClearCache() int {
	logrus.Info("clearing cache")
	oldLength := len(s.books)
	s.books = make(map[string]*items.Book)
	return oldLength
}

func (s *BookStorer) Length() int {
	return len(s.books)
}

func (s *BookStorer) DatabaseID() string {
	return s.databaseID
}

func (s *BookStorer) LoadBooksFromCSV() {
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
		book := s.CSVBookToBook(scanner.Text())
		s.StoreItem(book)
	}

	if err := scanner.Err(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Loaded %v books from CSV", count)
}

func (s *BookStorer) LoadBooksFromNotion(ctx context.Context, databaseID string) error {
	if os.Getenv("NOTION_CACHE") == "false" {
		logrus.Info("Notion cache disabled")
		return nil
	}
	books, databaseID, err := notion.GetBookPagesFromDatabase(ctx, databaseID)
	if err != nil {
		return err
	}
	for _, b := range books {
		s.StoreItem(b)
	}
	logrus.Infof("Loaded %v books from Notion", len(books))
	s.databaseID = databaseID
	return nil
}

func (s *BookStorer) CSVBookToBook(line string) *items.Book {
	parts := strings.Split(line, ",")
	if len(parts) != 6 {
		logrus.Warn("Line doesn't have 6 parts", map[string]string{"line": line})
		return nil
	}

	pages, err := strconv.Atoi(removeQuotes(parts[3]))
	if err != nil {
		logrus.Error(err)
	}
	return &items.Book{
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
