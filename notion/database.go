package notion

import (
	"context"
	"errors"

	"github.com/domgoodwin/bookscan/items"
	"github.com/jomei/notionapi"
	"github.com/sirupsen/logrus"
)

func GetBookPagesFromDatabase(ctx context.Context, databaseID string) ([]*items.Book, string, error) {
	if databaseID == "" {
		return nil, "", errors.New("book database ID must be supplied")
	}
	var books []*items.Book
	var nextCursor notionapi.Cursor
	for {
		rsp, err := client.Database.Query(ctx, notionapi.DatabaseID(databaseID), &notionapi.DatabaseQueryRequest{
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

	return books, databaseID, nil
}

func GetBookPageFromISBN(ctx context.Context, databaseID, ISBN string) (*items.Book, string, error) {
	if databaseID == "" {
		return nil, "", errors.New("book database ID must be supplied")
	}

	rsp, err := client.Database.Query(ctx, notionapi.DatabaseID(databaseID), &notionapi.DatabaseQueryRequest{
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

func GetRecordPagesFromDatabase(ctx context.Context, databaseID string) ([]*items.Record, string, error) {
	if databaseID == "" {
		return nil, "", errors.New("book database ID must be supplied")
	}
	var records []*items.Record
	var nextCursor notionapi.Cursor
	for {

		rsp, err := client.Database.Query(ctx, notionapi.DatabaseID(databaseID), &notionapi.DatabaseQueryRequest{
			PageSize:    100,
			StartCursor: nextCursor,
		})
		if err != nil {
			return nil, "", err
		}

		for _, page := range rsp.Results {
			records = append(records, notionPageToRecord(page))
		}

		if !rsp.HasMore {
			break
		}
		nextCursor = rsp.NextCursor
	}

	return records, databaseID, nil
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

func notionPageToRecord(p notionapi.Page) *items.Record {
	titleProperty := p.Properties[columnTitle].(*notionapi.TitleProperty)
	artistsProperty := p.Properties[colummArtists].(*notionapi.MultiSelectProperty)
	var artists []string
	for _, author := range artistsProperty.MultiSelect {
		artists = append(artists, author.Name)
	}
	barcodeProperty := p.Properties[columnBarcode].(*notionapi.RichTextProperty)
	linkProperty := p.Properties[columnLink].(*notionapi.URLProperty)
	yearProperty := p.Properties[colummYear].(*notionapi.NumberProperty)

	return &items.Record{
		Barcode: barcodeProperty.RichText[0].Text.Content,
		Title:   titleProperty.Title[0].Text.Content,
		Artists: artists,
		Link:    linkProperty.URL,
		Year:    int(yearProperty.Number),
	}

}
