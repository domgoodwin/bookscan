package items

import (
	"fmt"
	"os"
	"strings"
)

const recordCsvFilePath = "./records.csv"

type Record struct {
	Title    string
	Artists  []string
	Barcode  string
	Year     int
	Link     string
	CoverURL string
}

func (r Record) Type() string {
	return ItemTypeRecord
}

func (r Record) Artist() string {
	if len(r.Artists) == 0 {
		return ""
	}
	return r.Artists[0]
}

func (r Record) Info() string {
	return fmt.Sprintf("%v by %v", r.Title, r.Artist())
}

func (r Record) FullInfo() string {
	return fmt.Sprintf(`
	Title: %v
	Artists: %v
	Barcode: %v
	Year: %v
	Link: %v
	Cover URL: %v
	`, r.Title, r.Artists, r.Barcode, r.Year, r.Link, r.CoverURL)
}

func (r Record) FullInfoFields() map[string]string {
	return map[string]string{
		"title":     r.Title,
		"authors":   strings.Join(r.Artists, ";"),
		"barcode":   r.Barcode,
		"year":      fmt.Sprint(r.Year),
		"link":      r.Link,
		"cover_url": r.CoverURL,
	}
}

func (r Record) csvLine() string {
	return fmt.Sprintf("\"%v\",\"%v\",\"%v\",\"%v\",\"%v\",\"%v\"", r.Title, r.Artist(), r.Barcode, r.Year, r.Link, r.CoverURL)
}

func (r Record) StoreInCSV() error {
	if os.Getenv("CSV_SAVE") == "false" {
		return nil
	}
	f, err := os.OpenFile(recordCsvFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%v\n", r.csvLine())); err != nil {
		return err
	}
	return nil
}
