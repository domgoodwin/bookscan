package database

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID          string `bun:"id,pk"`
	Name        string `bun:"name"`
	AvatarURL   string `bun:"avatar_url"`
	Email       string `bun:"email"`
	LatestBotID string `bun:"latest_bot_id"`
}

func UserExists(ctx context.Context, id string) (bool, error) {
	user := new(User)
	return db.NewSelect().Model(user).Where("id = ?", id).Exists(ctx)
}

func GetUserByID(ctx context.Context, id string) (*User, error) {
	user := new(User)
	err := db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	return user, err
}

func CreateUser(ctx context.Context, user *User) error {
	logrus.Infof("creating user in database, user:%v", user.ID)
	_, err := db.NewUpdate().Model(user).Exec(ctx)
	return err
}
