package migrate

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	psql "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	sslMode          = "sslmode=disable"
	dataBase         = "postgres"
	pathToMigrations = "file://internal/data/"
)

func Up(dataBaseURL string) error {

	db, err := sql.Open(dataBase, dataBaseURL+sslMode)
	if err != nil {
		return err
	}

	driver, err := psql.WithInstance(db, &psql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		pathToMigrations,
		dataBase,
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != migrate.ErrNoChange {
		return err
	}

	return nil
}
