package notion

import (
	"github.com/jomei/notionapi"
	"github.com/sirupsen/logrus"
)

func GetClient(token string) *notionapi.Client {
	logrus.Info("Setting up notion client")
	return notionapi.NewClient(notionapi.Token(token))
}

type NotionClient struct {
	*notionapi.Client
	Token             string
	UserID            string
	PageID            string
	BooksDatabaseID   string
	RecordsDatabaseID string
}
