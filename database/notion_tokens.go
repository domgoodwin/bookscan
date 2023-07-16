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

	BookDatabaseID   string `bun:"book_database_id"`
	RecordDatabaseID string `bun:"record_database_id"`
}

func SaveToken(ctx context.Context, user *User, token *NotionToken) error {
	logrus.Infof("saving token in database, user:%v token:%v", user.ID, token.AccessToken)
	// Check user exists, if not create
	exists, err := UserExists(ctx, user.ID)
	if err != nil {
		logrus.Error("error checking if user exists")
		return err
	}
	if !exists {
		err := CreateUser(ctx, user)
		if err != nil {
			logrus.Error("error creating user")
			return err
		}
	}

	exists, err = NotionTokenExists(ctx, token.BotID)
	if err != nil {
		logrus.Error("error checking if token exists")
		return err
	}
	if !exists {
		_, err = db.NewInsert().Model(token).Exec(ctx)
		if err != nil {
			logrus.Error("error creating notion token")
			return err
		}
		return nil
	}

	_, err = db.NewUpdate().Model(token).Where("id = ?", token.BotID).Exec(ctx)
	return err
}

func GetNotionTokenByUserID(ctx context.Context, userID string) (*NotionToken, error) {
	notionToken := new(NotionToken)
	err := db.NewSelect().Model(notionToken).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		logrus.Error("error selecting token")
		return nil, err
	}
	return notionToken, nil
}
func NotionTokenExists(ctx context.Context, id string) (bool, error) {
	notionToken := new(NotionToken)
	return db.NewSelect().Model(notionToken).Where("id = ?", id).Exists(ctx)
}
