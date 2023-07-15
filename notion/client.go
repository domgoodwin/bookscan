package notion

import (
	"os"

	"github.com/jomei/notionapi"
)

var client *notionapi.Client

func SetupClient() {
	if client == nil {
		// apiKey := os.Getenv("NOTION_API_KEY")
		// if apiKey == "" {
		// 	panic("NOTION_API_KEY is not set")
		// }
		clientID := os.Getenv("NOTION_OAUTH_CLIENT_ID")
		clientSecret := os.Getenv("NOTION_OAUTH_CLIENT_SECRET")
		if clientID == "" || clientSecret == "" {
			panic("client id and secret must be set")
		}
		client = notionapi.NewClient(notionapi.Token("placeholder"),
			notionapi.WithOAuthAppCredentials(clientID, clientSecret),
		)
	}
}
