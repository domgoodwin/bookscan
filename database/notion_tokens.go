package database

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type NotionToken struct {
	bun.BaseModel `bun:"table:notion_tokens,alias:nt"`

	BotID  string `bun:"id,pk,default:gen_random_uuid()"`
	UserID string `bun:"user_id,notnull"`
	User   *User  `bun:"rel:belongs-to,join:user_id=id"`

	AccessToken string `bun:"access_token,notnull"`
	// Other notion info: https://developers.notion.com/docs/authorization#step-4-notion-responds-with-an-access_token-and-additional-information
	DuplicatedTemplateID string `bun:"duplicated_template_id"`
	WorkspaceIcon        string `bun:"workspace_icon"`
	WorkspaceID          string `bun:"workspace_id"`
	WorkspaceName        string `bun:"workspace_name"`
}

func SaveToken(ctx context.Context, user *User, token *NotionToken) error {
	logrus.Infof("saving token in database, user:%v token:%v", user.ID, token.AccessToken)
	// Check user exists, if not create
	if exists, _ := UserExists(ctx, user.ID); !exists {
		err := CreateUser(ctx, user)
		if err != nil {
			return err
		}
	}

	_, err := db.NewInsert().Model(token).Exec(ctx)
	return err
}
