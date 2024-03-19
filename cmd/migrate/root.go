package migration

import (
	"github.com/arfan21/project-sprint-social-media-api/migration"
	dbpostgres "github.com/arfan21/project-sprint-social-media-api/pkg/db/postgres"
	"github.com/jackc/pgx/v5/stdlib"
)

func initMigration() (*migration.Migration, error) {
	db, err := dbpostgres.NewPgx()
	if err != nil {
		return nil, err
	}

	sqlDB := stdlib.OpenDBFromPool(db)

	migration, err := migration.New(sqlDB)
	if err != nil {
		return nil, err
	}

	return migration, nil

}
