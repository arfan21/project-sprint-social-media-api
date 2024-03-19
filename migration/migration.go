package migration

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/arfan21/project-sprint-social-media-api/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var embedMigrations embed.FS

type Migration struct {
	db *sql.DB
	*migrate.Migrate
}

func New(db *sql.DB) (*Migration, error) {
	if db == nil {
		return &Migration{}, errors.New("db is nil")
	}

	source, err := iofs.New(embedMigrations, ".")
	if err != nil {
		return &Migration{}, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return &Migration{}, err
	}

	migrator, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		return &Migration{}, err
	}

	migrator.Log = &loggerMigration{}

	return &Migration{db: db, Migrate: migrator}, nil
}

type loggerMigration struct{}

func (l *loggerMigration) Printf(format string, v ...interface{}) {
	logger.Log(context.Background()).Info().Msgf(format, v...)
}

func (l *loggerMigration) Verbose() bool {
	return true
}
