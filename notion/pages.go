package notion

import (
	"context"

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
