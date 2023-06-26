package book

import (
	"fmt"
	"os"
)

type Book struct {
	Title   string
	Authors []string
	ISBN    string
	Pages   int
	Link    string
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
	`, b.Title, b.Authors, b.ISBN, b.Pages, b.Link)
}

func (b Book) FullInfoFields() map[string]string {
	return map[string]string{
		"title":   b.Title,
		"authors": fmt.Sprint(b.Authors),
		"isbn":    b.ISBN,
		"pages":   fmt.Sprint(b.Pages),
		"link":    b.Link,
	}
}

func (b Book) csvLine() string {
	return fmt.Sprintf("\"%v\",\"%v\",\"%v\",\"%v\",\"%v\"", b.Title, b.Author(), b.ISBN, b.Pages, b.Link)
}

func (b Book) StoreInCSV(filePath string) error {
	f, err := os.OpenFile(filePath,
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
