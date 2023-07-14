package database

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var sqldb *sql.DB
var db *bun.DB

var databaseModels = map[string]interface{}{
	"notion_tokens": &NotionToken{},
	"users":         &User{},
}

func Setup() error {
	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		return errors.New("POSTGRES_URL not set")
	}
	sqldb = sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxOpenConns)
	db = bun.NewDB(sqldb, pgdialect.New())

	// A query hook runs before and after executing
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	ctx := context.Background()
	for k, model := range databaseModels {
		logrus.Infof("setting up model %v", k)
		err := setupDatabase(ctx, model)
		if err != nil {
			logrus.Errorf("failed to setup model: %v %v", k, err)
			return err
		}
	}
	logrus.Info("Setup connection to postgres")

	return nil
}

func setupDatabase(ctx context.Context, model interface{}) error {
	_, err := db.NewCreateTable().IfNotExists().Model(model).Exec(ctx)
	return err
}
