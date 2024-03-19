package userrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/arfan21/project-sprint-social-media-api/internal/entity"
	"github.com/arfan21/project-sprint-social-media-api/pkg/constant"
	dbpostgres "github.com/arfan21/project-sprint-social-media-api/pkg/db/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r Repository) Create(ctx context.Context, data entity.User) (err error) {
	query := `
		INSERT INTO users (id, email, phone, name, password)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = r.db.Exec(ctx, query, data.ID, data.Email, data.Phone, data.Name, data.Password)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLUniqueViolation {
				err = constant.ErrEmailOrPhoneAlreadyRegistered
			}
		}

		err = fmt.Errorf("user.repository.Create: failed to create user: %w", err)
		return
	}

	return
}

func (r Repository) GetByCredential(ctx context.Context, credentialType, credentialValue string) (data entity.User, err error) {
	credType := "email"
	if credentialType == "phone" {
		credType = "phone"
	}
	query := `
		SELECT id, name, password, email, phone
		FROM users
		WHERE ` + credType + ` = $1
	`

	err = r.db.QueryRow(ctx, query, credentialValue).Scan(
		&data.ID,
		&data.Name,
		&data.Password,
		&data.Email,
		&data.Phone,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrUserNotFound
		}

		err = fmt.Errorf("user.repository.GetByCredential: failed to get user by email: %w", err)

		return
	}

	return
}

func (r Repository) GetByID(ctx context.Context, id string) (data entity.User, err error) {
	query := `
		SELECT id, name, email, phone
		FROM users
		WHERE id = $1
	`

	err = r.db.QueryRow(ctx, query, id).Scan(
		&data.ID,
		&data.Name,
		&data.Email,
		&data.Phone,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrUserNotFound
		}

		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLInvalidUUID {
				err = constant.ErrUserNotFound
			}
		}

		err = fmt.Errorf("user.repository.GetByID: failed to get user by id: %w", err)

		return
	}

	return
}

func (r Repository) AddFriend(ctx context.Context, userIdAdder, userIdAdded string) (err error) {
	query := `
		INSERT INTO friends (userIdAdder, userIdAdded)
		VALUES ($1, $2)
	`

	_, err = r.db.Exec(ctx, query, userIdAdder, userIdAdded)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLUniqueViolation {
				err = constant.ErrFriendUseralreadyAdded
			}
		}

		err = fmt.Errorf("user.repository.AddFriend: failed to add friend: %w", err)
		return
	}

	return
}

func (r Repository) DeleteFriend(ctx context.Context, userIdAdder, userIdAdded string) (err error) {
	query := `
		DELETE FROM friends
		WHERE( userIdAdder = $1 AND userIdAdded = $2 ) OR ( userIdAdder = $2 AND userIdAdded = $1)
	`

	cmd, err := r.db.Exec(ctx, query, userIdAdder, userIdAdded)
	if err != nil {
		err = fmt.Errorf("user.repository.DeleteFriend: failed to delete friend: %w", err)
		return
	}

	if cmd.RowsAffected() == 0 {
		err = fmt.Errorf("user.repository.DeleteFriend: failed to delete friend: %w", constant.ErrFriendUserNotAdded)
		return
	}

	return
}
