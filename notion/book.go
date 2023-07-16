package notion

import (
	"context"
	"os"

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
)

func (c *NotionClient) GetBookPagesFromDatabase(ctx context.Context) ([]*items.Book, string, error) {
	var books []*items.Book
	var nextCursor notionapi.Cursor
	for {
		rsp, err := c.Database.Query(ctx, notionapi.DatabaseID(c.BooksDatabaseID), &notionapi.DatabaseQueryRequest{
			PageSize:    100,
			StartCursor: nextCursor,
		})
		if err != nil {
			logrus.Error(err)
			return nil, "", err
		}

		for _, page := range rsp.Results {
			books = append(books, notionPageToBook(page))
		}

		if !rsp.HasMore {
			break
		}
		nextCursor = rsp.NextCursor
	}

	return books, c.BooksDatabaseID, nil
}

func (c *NotionClient) GetBookPageFromISBN(ctx context.Context, ISBN string) (*items.Book, string, error) {
	rsp, err := c.Database.Query(ctx, notionapi.DatabaseID(c.BooksDatabaseID), &notionapi.DatabaseQueryRequest{
		PageSize: 10,
		Filter: notionapi.PropertyFilter{
			Property: columnISBN,
			RichText: &notionapi.TextFilterCondition{
				Equals: ISBN,
			},
		},
	})
	if err != nil {
		return nil, "", err
	}

	if len(rsp.Results) == 0 {
		// Not found don't error
		return nil, "", nil
	}
	return notionPageToBook(rsp.Results[0]), rsp.Results[0].URL, nil
}

func (c *NotionClient) AddBookToDatabase(ctx context.Context, book *items.Book) (string, error) {

	// Check Book doesn't exist
	notionBook, notionURL, err := c.GetBookPageFromISBN(ctx, book.ISBN)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	if notionBook != nil {
		logrus.Infof("Book %v already exists in database: %v", notionBook.ISBN, notionURL)
		return notionURL, nil
	}

	logrus.Infof("saving %v in %v", book.Title, c.BooksDatabaseID)
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
	page, err := c.Page.Create(ctx, &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: notionapi.DatabaseID(c.BooksDatabaseID),
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

func notionPageToBook(p notionapi.Page) *items.Book {
	isbnPropety := p.Properties[columnISBN].(*notionapi.TitleProperty)
	authorsProperty := p.Properties[columnAuthors].(*notionapi.MultiSelectProperty)
	var authors []string
	for _, author := range authorsProperty.MultiSelect {
		authors = append(authors, author.Name)
	}
	titleProperty := p.Properties[columnTitle].(*notionapi.RichTextProperty)
	linkProperty := p.Properties[columnLink].(*notionapi.URLProperty)
	pagesProperty := p.Properties[columnPages].(*notionapi.NumberProperty)

	return &items.Book{
		Title:   titleProperty.RichText[0].Text.Content,
		Authors: authors,
		ISBN:    isbnPropety.Title[0].Text.Content,
		Link:    linkProperty.URL,
		Pages:   int(pagesProperty.Number),
	}

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
