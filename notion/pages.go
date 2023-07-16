package notion

import (
	"context"

	"github.com/domgoodwin/bookscan/database"
	"github.com/domgoodwin/bookscan/items"
	"github.com/jomei/notionapi"
	"github.com/sirupsen/logrus"
)

func (c *NotionClient) GetDatabaseIDs(ctx context.Context, pageID string) (string, string, error) {
	var bookDBID, recordDBID string
	blocksRsp, err := c.Block.GetChildren(ctx, notionapi.BlockID(pageID), &notionapi.Pagination{PageSize: 100})
	if err != nil {
		logrus.Errorf("error getting block children: %v", pageID)
		return "", "", err
	}
	for _, listBlock := range blocksRsp.Results {
		if listBlock.GetType() == notionapi.BlockTypeChildDatabase {
			block := listBlock.(*notionapi.ChildDatabaseBlock)
			switch block.ChildDatabase.Title {
			case "Books":
				bookDBID = string(block.GetID())
			case "Records":
				recordDBID = string(block.GetID())
			}
		}
	}
	return bookDBID, recordDBID, nil
}

func GetAllPagesFromAllDatabases(ctx context.Context) (map[string]map[string]*items.Book, map[string]map[string]*items.Record, error) {
	booksToDatabaseID := make(map[string]map[string]*items.Book)
	recordsToDatabaseID := make(map[string]map[string]*items.Record)
	tokens, err := database.ListNotionTokens(ctx)
	if err != nil {
		logrus.Error(err)
		return nil, nil, err
	}
	for _, token := range tokens {
		booksToDatabaseID[token.BookDatabaseID] = make(map[string]*items.Book)
		recordsToDatabaseID[token.RecordDatabaseID] = make(map[string]*items.Record)
		notionClient := GetClient(token.AccessToken)
		c := &NotionClient{notionClient, token.AccessToken, token.UserID, token.DuplicatedTemplateID, token.BookDatabaseID, token.RecordDatabaseID}
		books, _, err := c.GetBookPagesFromDatabase(ctx)
		if err != nil {
			logrus.Error(err)
			return nil, nil, err
		}
		for _, book := range books {
			book := book
			booksToDatabaseID[token.BookDatabaseID][book.ISBN] = book
		}
		records, _, err := c.GetRecordPagesFromDatabase(ctx)
		if err != nil {
			logrus.Error(err)
			return nil, nil, err
		}
		for _, record := range records {
			record := record
			recordsToDatabaseID[token.RecordDatabaseID][record.Barcode] = record
		}
		logrus.Debugf("loading %v books and %v records from notion for %v", len(booksToDatabaseID[token.BookDatabaseID]), len(recordsToDatabaseID[token.RecordDatabaseID]), token.UserID)
	}
	return booksToDatabaseID, recordsToDatabaseID, nil
}
