package lookup

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/domgoodwin/bookscan/book"
)

const (
	openLibraryURL = "https://openlibrary.org"
)

func OpenLibraryLookup(isbn string) (*book.Book, error) {
	url := fmt.Sprintf("%v/isbn/%v.json", openLibraryURL, isbn)
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("code %v", rsp.StatusCode))
	}
	defer rsp.Body.Close()
	edition := &openLibraryEdition{}
	err = json.NewDecoder(rsp.Body).Decode(edition)
	if err != nil {
		return nil, err
	}
	edition.ISBN = isbn

	return edition.Book()
}

func (o openLibraryEdition) Book() (*book.Book, error) {
	authors, err := o.Authors()
	if err != nil {
		return nil, err
	}
	return &book.Book{
		Title:   o.Title,
		Authors: authors,
		ISBN:    o.ISBN,
		Pages:   o.NumberOfPages,
		Link:    fmt.Sprintf("%v%v", openLibraryURL, o.Key),
	}, nil
}

func (o openLibraryEdition) Authors() ([]string, error) {
	var authors []string
	// If author is empty we need to lookup the work from our edition
	authorKeys := o.AuthorKeys
	if len(o.AuthorKeys) == 0 && len(o.Works) > 0 {
		rsp, err := http.Get(fmt.Sprintf("%v%v.json", openLibraryURL, o.Works[0].Key))
		if err != nil {
			return nil, err
		}
		work := &openLibraryWork{}
		err = json.NewDecoder(rsp.Body).Decode(work)
		if err != nil {
			return nil, err
		}
		for _, author := range work.Authors {
			authorKeys = append(authorKeys, author.Author)
		}
	}

	// Lookup author name from keys
	for _, authorKey := range authorKeys {
		rsp, err := http.Get(fmt.Sprintf("%v%v.json", openLibraryURL, authorKey.Key))
		if err != nil {
			return nil, err
		}
		author := &openLibraryAuthor{}
		err = json.NewDecoder(rsp.Body).Decode(author)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author.Name)
	}
	return authors, nil

}

type openLibraryEdition struct {
	Publishers    []string               `json:"publishers"`
	NumberOfPages int                    `json:"number_of_pages"`
	Key           string                 `json:"key"`
	AuthorKeys    []openLibraryKey       `json:"authors"`
	Title         string                 `json:"title"`
	Identifiers   openLibraryIdentifiers `json:"identifiers"`
	PublishDate   string                 `json:"publish_date"`
	Works         []openLibraryKey       `json:"works"`
	ISBN          string
}

type openLibraryWork struct {
	Title   string                 `json:"title"`
	Authors []openLibrarySubAuthor `json:"authors"`
	Key     string                 `json:"key"`
}

type openLibraryAuthor struct {
	Name string `json:"name"`
}

type openLibrarySubAuthor struct {
	Author openLibraryKey `json:"author"`
}

type openLibraryKey struct {
	Key string `json:"key"`
}

type openLibraryIdentifiers struct {
	Goodreads []string `json:"goodreads"`
	Amazon    []string `json:"amazon"`
}
