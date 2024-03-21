package user

import (
	"context"

	"github.com/arfan21/project-sprint-social-media-api/internal/entity"
	"github.com/arfan21/project-sprint-social-media-api/internal/model"
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
	GetList(ctx context.Context, filter model.UserGetListRequest) (data []entity.User, err error)
	GetCountList(ctx context.Context, filter model.UserGetListRequest) (count int, err error)
	IsFriend(ctx context.Context, userIdAdder, userIdAdded string) (isFriend bool, err error)
	GetListMap(ctx context.Context, filter model.UserGetListRequest) (data map[string]entity.User, err error)
	UpdatePhone(ctx context.Context, userId, phone string) (err error)
	UpdateEmail(ctx context.Context, userId, email string) (err error)
	UpdateProfile(ctx context.Context, data entity.User) (err error)
}
