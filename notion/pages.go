package notion

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/domgoodwin/bookscan/items"
	"github.com/jomei/notionapi"
	"github.com/sirupsen/logrus"
)

const (
	// Both
	columnTitle = "Title"
	columnLink  = "Link"
	// Book
	columnAuthors = "Authors"
	columnISBN    = "ISBN"
	columnPages   = "Pages"
	// Record
	colummArtists = "Artists"
	colummYear    = "Year"
	columnBarcode = "Barcode"
)

func AddBookToDatabase(ctx context.Context, book *items.Book, databaseID string) (string, error) {
	if databaseID == "" {
		return "", errors.New("book database ID must be supplied")
	}
	logrus.Infof("saving %v in %v", book.Title, databaseID)
	if os.Getenv("NOTION_SAVE") == "false" {
		return "", nil
	}

	var bookCover *notionapi.Image
	if book.CoverURL != "" {
		bookCover = &notionapi.Image{
			Type: notionapi.FileTypeExternal,
			External: &notionapi.FileObject{
				URL: book.CoverURL,
			},
		}
	}
	page, err := client.Page.Create(ctx, &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: notionapi.DatabaseID(databaseID),
		},
		Properties: bookToDatabaseProperties(book),
		Cover:      bookCover,
	})
	if err != nil {
		logrus.Errorf("error: %v", err)
		return "", err
	}
	logrus.Infof("saving %v, url: %v", book.Title, page.URL)
	return page.URL, nil
}

func AddRecordToDatabase(ctx context.Context, record *items.Record, databaseID string) (string, error) {
	if record == nil {
		return "", errors.New("nil record")
	}
	logrus.Infof("Adding record to database: %v", record)
	if databaseID == "" {
		return "", errors.New("book database ID must be supplied")
	}
	logrus.Infof("saving %v in %v", record.Title, databaseID)
	if os.Getenv("NOTION_SAVE") == "false" {
		return "", nil
	}

	var pageCover *notionapi.Image
	if record.CoverURL != "" {
		pageCover = &notionapi.Image{
			Type: notionapi.FileTypeExternal,
			External: &notionapi.FileObject{
				URL: record.CoverURL,
			},
		}
	}
	page, err := client.Page.Create(ctx, &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: notionapi.DatabaseID(databaseID),
		},
		Properties: recordToDatabaseProperties(record),
		Cover:      pageCover,
	})
	if err != nil {
		logrus.Errorf("error: %v", err)
		return "", err
	}
	logrus.Infof("saving %v, url: %v", record.Title, page.URL)
	return page.URL, nil
}

func bookToDatabaseProperties(b *items.Book) notionapi.Properties {

	authors := notionapi.MultiSelectProperty{
		Type:        notionapi.PropertyTypeMultiSelect,
		MultiSelect: []notionapi.Option{},
	}
	for _, author := range b.Authors {
		authors.MultiSelect = append(
			authors.MultiSelect,
			notionapi.Option{
				Name: sanitseSelectValue(author),
			},
		)
	}
	return notionapi.Properties{
		columnISBN: notionapi.TitleProperty{
			Type: notionapi.PropertyTypeTitle,
			Title: []notionapi.RichText{
				{
					Text: &notionapi.Text{
						Content: b.ISBN,
					},
				},
			},
		},
		columnAuthors: authors,
		columnTitle: notionapi.RichTextProperty{
			Type: notionapi.PropertyTypeRichText,
			RichText: []notionapi.RichText{
				{
					Text: &notionapi.Text{
						Content: b.Title,
					},
				},
			},
		},
		columnPages: notionapi.NumberProperty{
			Type:   notionapi.PropertyTypeNumber,
			Number: float64(b.Pages),
		},
		columnLink: notionapi.URLProperty{
			Type: notionapi.PropertyTypeURL,
			URL:  b.Link,
		},
	}
}

func recordToDatabaseProperties(b *items.Record) notionapi.Properties {
	artists := notionapi.MultiSelectProperty{
		Type:        notionapi.PropertyTypeMultiSelect,
		MultiSelect: []notionapi.Option{},
	}
	for _, artist := range b.Artists {
		artists.MultiSelect = append(
			artists.MultiSelect,
			notionapi.Option{
				Name: sanitseSelectValue(artist),
			},
		)
	}
	return notionapi.Properties{
		columnTitle: notionapi.TitleProperty{
			Type: notionapi.PropertyTypeTitle,
			Title: []notionapi.RichText{
				{
					Text: &notionapi.Text{
						Content: b.Title,
					},
				},
			},
		},
		colummArtists: artists,
		columnBarcode: notionapi.RichTextProperty{
			Type: notionapi.PropertyTypeRichText,
			RichText: []notionapi.RichText{
				{
					Text: &notionapi.Text{
						Content: b.Barcode,
					},
				},
			},
		},
		colummYear: notionapi.NumberProperty{
			Type:   notionapi.PropertyTypeNumber,
			Number: float64(b.Year),
		},
		columnLink: notionapi.URLProperty{
			Type: notionapi.PropertyTypeURL,
			URL:  b.Link,
		},
	}
}

func sanitseSelectValue(in string) string {
	return strings.Replace(in, ",", " ", -1)
}
