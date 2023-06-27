package book

import (
	"fmt"
	"os"
	"strings"
)

const csvFilePath = "./books.csv"

type Book struct {
	Title    string
	Authors  []string
	ISBN     string
	Pages    int
	Link     string
	CoverURL string
}

func (b Book) Author() string {
	if len(b.Authors) == 0 {
		return ""
	}
	return b.Authors[0]
}

func (b Book) Info() string {
	return fmt.Sprintf("%v by %v", b.Title, b.Author())
}

func (b Book) FullInfo() string {
	return fmt.Sprintf(`
	Title: %v
	Authors: %v
	ISBN: %v
	Pages: %v
	Link: %v
	Cover URL: %v
	`, b.Title, b.Authors, b.ISBN, b.Pages, b.Link, b.CoverURL)
}

func (b Book) FullInfoFields() map[string]string {
	return map[string]string{
		"title":     b.Title,
		"authors":   strings.Join(b.Authors, ";"),
		"isbn":      b.ISBN,
		"pages":     fmt.Sprint(b.Pages),
		"link":      b.Link,
		"cover_url": b.CoverURL,
	}
}

func (b Book) csvLine() string {
	return fmt.Sprintf("\"%v\",\"%v\",\"%v\",\"%v\",\"%v\",\"%v\"", b.Title, b.Author(), b.ISBN, b.Pages, b.Link, b.CoverURL)
}

func (b Book) StoreInCSV() error {
	if os.Getenv("CSV_SAVE") == "false" {
		return nil
	}
	f, err := os.OpenFile(csvFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%v\n", b.csvLine())); err != nil {
		return err
	}
	return nil
}
