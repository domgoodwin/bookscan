package lookup

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/domgoodwin/bookscan/book"
	"github.com/sirupsen/logrus"
)

const (
	googlebooksURL = "https://www.googleapis.com/books/v1"
)

type googlebooksDataSource struct{}

func (g googlebooksDataSource) Name() string {
	return "googlebooks"
}

func (g googlebooksDataSource) LookupISBN(isbn string) (*book.Book, error) {
	logrus.Infof("looking up isbn with google books API: %v", isbn)
	url := fmt.Sprintf("%v/volumes?q=isbn:%v", googlebooksURL, isbn)
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("code %v", rsp.StatusCode)
	}
	defer rsp.Body.Close()
	response := &googlebooksVolumesResponse{}
	err = json.NewDecoder(rsp.Body).Decode(response)
	if err != nil {
		return nil, err
	}
	if len(response.Items) == 0 {
		return nil, fmt.Errorf("not found")
	}
	// not sure if we need this
	if len(response.Items) > 1 {
		logrus.Infof("found too many books: %v", response)
		return nil, fmt.Errorf("found too many")
	}
	book := response.Items[0].Book()
	logrus.Infof("found book with google books API: %v", book)
	return book, nil
}

type googlebooksVolumesResponse struct {
	Kind  string              `json:"kind"`
	Items []googlebooksVolume `json:"items"`
}

type googlebooksVolume struct {
	Kind       string                `json:"kind"`
	ID         string                `json:"id"`
	ETag       string                `json:"etag"`
	SelfLink   string                `json:"selfLink"`
	VolumeInfo googlebooksVolumeInfo `json:"volumeInfo"`
}

func (v googlebooksVolume) Book() *book.Book {
	return &book.Book{
		Title:    v.VolumeInfo.Title,
		Authors:  v.VolumeInfo.Authors,
		ISBN:     v.ISBN(),
		Pages:    v.VolumeInfo.PageCount,
		Link:     v.SelfLink,
		CoverURL: v.CoverURL(),
	}
}

func (v googlebooksVolume) ISBN() string {
	var isbn13 string
	var isbn10 string
	for _, identifier := range v.VolumeInfo.IndustryIdentifiers {
		if identifier.Type == "ISBN_13" {
			isbn13 = identifier.Identifier
		}
		if identifier.Type == "ISBN_10" {
			isbn10 = identifier.Identifier
		}
	}
	if isbn13 == "" {
		return isbn10
	}
	return isbn13
}

func (v googlebooksVolume) CoverURL() string {
	if v.VolumeInfo.ImageLinks.Thumbnail == "" {
		return v.VolumeInfo.ImageLinks.SmallThunbnail
	}
	return v.VolumeInfo.ImageLinks.Thumbnail
}

type googlebooksVolumeInfo struct {
	Title               string                          `json:"title"`
	Subtitle            string                          `json:"subtitle"`
	Authors             []string                        `json:"authors"`
	PageCount           int                             `json:"pageCount"`
	ImageLinks          googlebooksImageLink            `json:"imageLinks"`
	PreviewLink         string                          `json:"previewLink"`
	IndustryIdentifiers []googlebooksIndustryIdentifier `json:"industryIdentifiers"`
}

type googlebooksImageLink struct {
	SmallThunbnail string `json:"smallThumbnail"`
	Thumbnail      string `json:"thumbnail"`
}

type googlebooksIndustryIdentifier struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}
