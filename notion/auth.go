package notion

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/domgoodwin/bookscan/database"
	"github.com/sirupsen/logrus"
)

func GetToken(ctx context.Context, code, redirectURI string) (string, string, error) {
	logrus.Infof("Creating token %v", code)
	rsp, err := getNotionToken(ctx, code, redirectURI)
	if err != nil {
		logrus.Error("error creating token")
		return "", "", err
	}
	logrus.Debug(rsp)
	logrus.Debug(rsp.Owner)
	userID := rsp.Owner.User.ID
	err = database.SaveToken(ctx,
		&database.User{
			ID:          userID,
			Name:        rsp.Owner.User.Name,
			AvatarURL:   rsp.Owner.User.AvatarURL,
			Email:       rsp.Owner.User.Person.Email,
			LatestBotID: rsp.BotID,
		},
		&database.NotionToken{
			BotID:                rsp.BotID,
			UserID:               userID,
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
		return "", "", err
	}

	token, err := database.CreateAPIAuth(ctx, userID)
	return userID, token, err
}

func getNotionToken(ctx context.Context, code, redirectURI string) (*NotionOauthTokenResponse, error) {
	clientID := os.Getenv("NOTION_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("NOTION_OAUTH_CLIENT_SECRET")
	url := "https://api.notion.com/v1/oauth/token"
	redirect := os.Getenv("NOTION_REDIRECT_URL")
	if redirectURI != "" {
		redirect = redirectURI
	}
	reqBody := &NotionOAuthRequest{
		Code:        code,
		GrantType:   "authorization_code",
		RedirectURI: redirect,
	}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	logrus.Debug(req)
	req.Header.Set("User-Agent", "Bookscan/0.1 +https://dgood.win")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(clientID+":"+clientSecret)))
	req.BasicAuth()
	logrus.Debugf("notion oauth token req: %v", req)
	rsp, err := client.Do(req)
	if err != nil {
		logrus.Error("error getting notion token")
		logrus.Error(err)
		return nil, err
	}
	if rsp.StatusCode != 200 {
		if rsp.Body != nil {
			defer rsp.Body.Close()
			errResponse := &NotionErrorResponse{}
			err = json.NewDecoder(rsp.Body).Decode(errResponse)
			if err != nil {
				return nil, err
			}

			logrus.Error(errResponse)
		}
		logrus.Errorf("non 200 response code: %v %v", rsp.StatusCode, rsp)
		return nil, fmt.Errorf("code %v", rsp.StatusCode)
	}

	defer rsp.Body.Close()
	logrus.Debugf("response: %v", rsp)
	response := &NotionOauthTokenResponse{}
	err = json.NewDecoder(rsp.Body).Decode(response)
	if err != nil {
		return nil, err
	}

	return response, err
}

type NotionOAuthRequest struct {
	Code        string `json:"code"`
	GrantType   string `json:"grant_type"`
	RedirectURI string `json:"redirect_uri"`
}

type NotionOauthTokenResponse struct {
	AccessToken          string `json:"access_token"`
	BotID                string `json:"bot_id"`
	DuplicatedTemplateID string `json:"duplicated_template_id"`
	Owner                struct {
		Type string `json:"type"`
		User struct {
			AvatarURL string `json:"avatar_url"`
			ID        string `json:"id"`
			Name      string `json:"name"`
			Person    struct {
				Email string `json:"email"`
			} `json:"person"`
		} `json:"user"`
	} `json:"owner"`
	WorkspaceIcon string `json:"workspace_icon"`
	WorkspaceID   string `json:"workspace_id"`
	WorkspaceName string `json:"workspace_name"`
}

type NotionErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
