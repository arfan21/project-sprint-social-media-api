package post

import (
	"context"

	"github.com/arfan21/project-sprint-social-media-api/internal/model"
)

type Service interface {
	Create(ctx context.Context, req model.PostRequest) (err error)
}
