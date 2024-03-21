package post

import (
	"context"

	"github.com/arfan21/project-sprint-social-media-api/internal/model"
)

type Service interface {
	Create(ctx context.Context, req model.PostRequest) (err error)
	CreateComment(ctx context.Context, req model.PostCommentRequest) (err error)
	GetList(ctx context.Context, req model.PostGetListRequest) (res []model.PostListResponse, count int, err error)
}
