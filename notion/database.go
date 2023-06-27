package notion

import (
	"context"

	"github.com/domgoodwin/bookscan/book"
	"github.com/jomei/notionapi"
)

func GetPagesFromDatabase(ctx context.Context) ([]*book.Book, error) {
	var books []*book.Book
	var nextCursor notionapi.Cursor
	for {

		rsp, err := client.Database.Query(ctx, booksDatabaseID, &notionapi.DatabaseQueryRequest{
			PageSize:    100,
			StartCursor: nextCursor,
		})
		if err != nil {
			return nil, err
		}

		for _, page := range rsp.Results {
			books = append(books, notionPageToBook(page))
		}

		if !rsp.HasMore {
			break
		}
		nextCursor = rsp.NextCursor
	}

	return books, nil
}

func notionPageToBook(p notionapi.Page) *book.Book {
	isbnPropety := p.Properties[columnISBN].(*notionapi.TitleProperty)
	authorsProperty := p.Properties[columnAuthors].(*notionapi.MultiSelectProperty)
	var authors []string
	for _, author := range authorsProperty.MultiSelect {
		authors = append(authors, author.Name)
	}
	titleProperty := p.Properties[columnTitle].(*notionapi.RichTextProperty)
	linkProperty := p.Properties[columnLink].(*notionapi.URLProperty)
	pagesProperty := p.Properties[columnPages].(*notionapi.NumberProperty)

	return &book.Book{
		Title:   titleProperty.RichText[0].Text.Content,
		Authors: authors,
		ISBN:    isbnPropety.Title[0].Text.Content,
		Link:    linkProperty.URL,
		Pages:   int(pagesProperty.Number),
	}

}
