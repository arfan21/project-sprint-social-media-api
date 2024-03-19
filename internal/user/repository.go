package user

import (
	"context"

	"github.com/arfan21/project-sprint-social-media-api/internal/entity"
	userrepo "github.com/arfan21/project-sprint-social-media-api/internal/user/repository"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Begin(ctx context.Context) (tx pgx.Tx, err error)
	WithTx(tx pgx.Tx) *userrepo.Repository

	Create(ctx context.Context, data entity.User) (err error)
	GetByCredential(ctx context.Context, credentialType, credentialValue string) (data entity.User, err error)
	GetByID(ctx context.Context, id string) (data entity.User, err error)
	AddFriend(ctx context.Context, userIdAdder, userIdAdded string) (err error)
	DeleteFriend(ctx context.Context, userIdAdder, userIdAdded string) (err error)
}
