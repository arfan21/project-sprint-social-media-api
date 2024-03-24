package post

import (
	"context"

	"github.com/arfan21/project-sprint-social-media-api/internal/entity"
	"github.com/arfan21/project-sprint-social-media-api/internal/model"
	postrepo "github.com/arfan21/project-sprint-social-media-api/internal/post/repository"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Begin(ctx context.Context) (tx pgx.Tx, err error)
	WithTx(tx pgx.Tx) *postrepo.Repository

	Create(ctx context.Context, data entity.Post) (err error)
	GetByID(ctx context.Context, id string) (data entity.Post, err error)
	CreateComment(ctx context.Context, data entity.PostComment) (err error)
	GetList(ctx context.Context, filter model.PostGetListRequest) (
		res []entity.Post,
		postIDs []string,
		userIdUnique map[string]struct{},
		err error,
	)
	GetCountList(ctx context.Context, filter model.PostGetListRequest) (count int, err error)
	GetCommentsByPostIDsMap(ctx context.Context, postIDs []string, userIDsUnique map[string]struct{}) (res map[string][]entity.PostComment, err error)
}
