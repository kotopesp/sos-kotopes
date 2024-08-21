package app

import (
	"context"
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	psql "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	_ "github.com/lib/pq"
)

const (
	sslMode          = "?sslmode=disable"
	dataBase         = "postgres"
	pathToMigrations = "file:///app/internal/data/"
)

func migrateUp(ctx context.Context, dataBaseURL string) error {

	db, err := sql.Open(dataBase, dataBaseURL+sslMode)
	if err != nil {
		logger.Log().Fatal(ctx, err.Error())
	}

	driver, err := psql.WithInstance(db, &psql.Config{})
	if err != nil {
		logger.Log().Fatal(ctx, err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		pathToMigrations,
		dataBase, driver)
	if err != nil {
		logger.Log().Fatal(ctx, err.Error())
	}

	if err := m.Up(); err != migrate.ErrNoChange {
		return err
	}

	return nil
}
