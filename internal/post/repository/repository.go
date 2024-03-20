package postrepo

import (
	"context"
	"fmt"

	"github.com/arfan21/project-sprint-social-media-api/internal/entity"
	dbpostgres "github.com/arfan21/project-sprint-social-media-api/pkg/db/postgres"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db dbpostgres.Queryer
}

func New(db dbpostgres.Queryer) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) Begin(ctx context.Context) (tx pgx.Tx, err error) {
	return r.db.Begin(ctx)
}

func (r Repository) WithTx(tx pgx.Tx) *Repository {
	r.db = tx
	return &r
}

func (r Repository) Create(ctx context.Context, data entity.Post) (err error) {
	query := `
		INSERT INTO posts (id, userId, body, tags)
		VALUES ($1, $2, $3, $4)
	`

	_, err = r.db.Exec(ctx, query, data.ID, data.UserID, data.Body, data.Tags)
	if err != nil {
		err = fmt.Errorf("post.repository.Create: failed to create post: %w", err)
		return
	}

	return
}
