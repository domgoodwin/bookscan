package notion

import (
	"context"

	"github.com/domgoodwin/bookscan/book"
	"github.com/jomei/notionapi"
)

const (
	booksDatabaseID = "4f311bbe86ce4dd4bdae93fa1206328f"
	columnTitle      = "Title"
	columnAuthors   = "Authors"
	columnISBN      = "ISBN"
	columnLink      = "Link"
	columnPages     = "Pages"
)

func AddBookToDatabase(ctx context.Context, book *book.Book) (string, error) {
	page, err := client.Page.Create(ctx, &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: booksDatabaseID,
		},
		Properties: bookToDatabaseProperties(book),
		Cover: &notionapi.Image{
			Type: notionapi.FileTypeExternal,
			External: &notionapi.FileObject{
				URL: book.CoverURL,
			},
		},
	})
	if err != nil {
		return "", err
	}
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
