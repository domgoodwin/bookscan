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
	// Record
	colummArtists = "Artists"
	colummYear    = "Year"
	columnBarcode = "Barcode"
)

func (c *NotionClient) GetRecordPagesFromDatabase(ctx context.Context) ([]*items.Record, string, error) {
	var records []*items.Record
	var nextCursor notionapi.Cursor
	for {

		rsp, err := c.Database.Query(ctx, notionapi.DatabaseID(c.RecordsDatabaseID), &notionapi.DatabaseQueryRequest{
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

	return records, c.RecordsDatabaseID, nil
}

func (c *NotionClient) AddRecordToDatabase(ctx context.Context, record *items.Record) (string, error) {
	if record == nil {
		return "", errors.New("nil record")
	}
	logrus.Infof("Adding record to database: %v", record)
	logrus.Infof("saving %v in %v", record.Title, c.RecordsDatabaseID)
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
	page, err := c.Page.Create(ctx, &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			Type:       notionapi.ParentTypeDatabaseID,
			DatabaseID: notionapi.DatabaseID(c.RecordsDatabaseID),
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
