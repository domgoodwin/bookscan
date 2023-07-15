package notion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/domgoodwin/bookscan/database"
	"github.com/jomei/notionapi"
	"github.com/sirupsen/logrus"
)

func GetToken(ctx context.Context, code string) (string, error) {
	logrus.Infof("Creating token %v", code)
	rsp, err := getNotionToken(ctx, code)
	if err != nil {
		logrus.Error("error creating token")
		return "", err
	}
	err = database.SaveToken(ctx,
		&database.User{
			ID:          string(rsp.Owner.ID),
			Name:        string(rsp.Owner.Name),
			AvatarURL:   string(rsp.Owner.AvatarURL),
			Email:       string(rsp.Owner.Person.Email),
			LatestBotID: rsp.BotID,
		},
		&database.NotionToken{
			BotID:                rsp.BotID,
			UserID:               string(rsp.Owner.ID),
			AccessToken:          rsp.AccessToken,
			DuplicatedTemplateID: rsp.DuplicatedTemplateID,
			WorkspaceIcon:        rsp.WorkspaceIcon,
			WorkspaceID:          rsp.WorkspaceID,
			WorkspaceName:        rsp.WorkspaceName,
		},
	)
	if err != nil {
		logrus.Error("error saving token in db")
		logrus.Error(err)
		return "", err
	}
	return rsp.AccessToken, nil
}

func getNotionToken(ctx context.Context, code string) (*NotionOauthTokenResponse, error) {
	clientID := os.Getenv("NOTION_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("NOTION_OAUTH_CLIENT_SECRET")
	url := fmt.Sprintf("https://%v:%vapi.notion.com/v1/oauth/token?grant_type=authorization_code&code=%v,response_type=code,owner=user", clientID, clientSecret, code)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("User-Agent", "Bookscan/0.1 +https://dgood.win")
	logrus.Debugf("notion oauth token req: %v", req)
	rsp, err := client.Do(req)
	if err != nil {
		logrus.Error("error getting notion token")
		logrus.Error(err)
		return nil, err
	}

	defer rsp.Body.Close()
	response := &NotionOauthTokenResponse{}
	err = json.NewDecoder(rsp.Body).Decode(response)
	if err != nil {
		return nil, err
	}

	return response, err
}

type NotionOauthTokenResponse struct {
	AccessToken          string `json:"access_token"`
	BotID                string `json:"bot_id"`
	DuplicatedTemplateID string `json:"duplicated_template_id"`
	Owner                notionapi.User
	WorkspaceIcon        string `json:"workspace_icon"`
	WorkspaceID          string `json:"workspace_id"`
	WorkspaceName        string `json:"workspace_name"`
}
