package notion

import (
	"os"

	"github.com/jomei/notionapi"
)

var client *notionapi.Client

func SetupClient() {
	if client == nil {
		apiKey := os.Getenv("NOTION_API_KEY")
		if apiKey == "" {
			panic("NOTION_API_KEY is not set")
		}
		client = notionapi.NewClient(notionapi.Token(apiKey))
	}
}
