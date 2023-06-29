package notion

import (
	"context"
	"os"

	"github.com/domgoodwin/bookscan/book"
	"github.com/jomei/notionapi"
	"github.com/sirupsen/logrus"
)

const (
	booksDatabaseID = "4f311bbe86ce4dd4bdae93fa1206328f"
	columnTitle     = "Title"
	columnAuthors   = "Authors"
	columnISBN      = "ISBN"
	columnLink      = "Link"
	columnPages     = "Pages"
)

func AddBookToDatabase(ctx context.Context, book *book.Book, databaseID string) (string, error) {
	if databaseID == "" {
		databaseID = booksDatabaseID
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

func bookToDatabaseProperties(b *book.Book) notionapi.Properties {

	authors := notionapi.MultiSelectProperty{
		Type:        notionapi.PropertyTypeMultiSelect,
		MultiSelect: []notionapi.Option{},
	}
	for _, author := range b.Authors {
		authors.MultiSelect = append(
			authors.MultiSelect,
			notionapi.Option{
				Name: author,
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
