package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type ApiAuth struct {
	bun.BaseModel `bun:"table:api_auth,alias:u"`

	Token  string `bun:"token,notnull"`
	UserID string `bun:"user_id,notnull"`
	User   *User  `bun:"rel:belongs-to,join:user_id=id"`
}

func CreateAPIAuth(ctx context.Context, userID string) (string, error) {
	apiAuth := new(ApiAuth)
	exists, err := db.NewSelect().Model(apiAuth).Where("user_id = ?", userID).Exists(ctx)
	if err != nil {
		logrus.Error("error checking if token exists")
		return "", err
	}

	token := generateSecureToken(64)
	apiAuth = &ApiAuth{
		UserID: userID,
		Token:  token,
	}
	if !exists {
		_, err := db.NewInsert().Model(apiAuth).Exec(ctx)
		return token, err
	}
	_, err = db.NewUpdate().Model(apiAuth).Where("user_id = ?", userID).Exec(ctx)
	return token, err

}

func CheckAPIToken(ctx context.Context, userID, token string) (bool, error) {
	apiAuth := new(ApiAuth)
	err := db.NewSelect().Model(apiAuth).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		logrus.Error("error reading API token")
		return false, err
	}
	return apiAuth.Token == token, nil
}

func ListNotionTokens(ctx context.Context) ([]NotionToken, error) {
	notionTokens := new([]NotionToken)
	err := db.NewSelect().Model(notionTokens).Scan(ctx)
	if err != nil {
		logrus.Error("error listing API tokens")
		return nil, err
	}
	return *notionTokens, nil
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
