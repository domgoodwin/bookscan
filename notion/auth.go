package notion

import (
	"context"

	"github.com/domgoodwin/bookscan/database"
	"github.com/jomei/notionapi"
	"github.com/sirupsen/logrus"
)

func GetToken(ctx context.Context, code string) (string, error) {
	if client == nil {
		SetupClient()
	}
	logrus.Infof("Creating token %v", code)
	rsp, err := client.Authentication.CreateToken(ctx, &notionapi.TokenCreateRequest{
		Code: code,
	})
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	userOwner := rsp.Owner.(notionapi.User)
	err = database.SaveToken(ctx,
		&database.User{
			ID:          string(userOwner.ID),
			Name:        string(userOwner.Name),
			AvatarURL:   string(userOwner.AvatarURL),
			Email:       string(userOwner.Person.Email),
			LatestBotID: rsp.BotId,
		},
		&database.NotionToken{
			BotID:                rsp.BotId,
			UserID:               string(userOwner.ID),
			AccessToken:          rsp.AccessToken,
			DuplicatedTemplateID: rsp.DuplicatedTemplateId,
			WorkspaceIcon:        rsp.WorkspaceIcon,
			WorkspaceID:          rsp.WorkspaceId,
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
