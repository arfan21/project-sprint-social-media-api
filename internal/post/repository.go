package post

import (
	"context"

	"github.com/arfan21/project-sprint-social-media-api/internal/entity"
	"github.com/arfan21/project-sprint-social-media-api/internal/model"
)

type Repository interface {
	Create(ctx context.Context, data entity.Post) (err error)
	GetByID(ctx context.Context, id string) (data entity.Post, err error)
	CreateComment(ctx context.Context, data entity.PostComment) (err error)
	GetList(ctx context.Context, filter model.PostGetListRequest) (res []entity.Post, userIDs []string, err error)
	GetCountList(ctx context.Context, filter model.PostGetListRequest) (count int, err error)
}
