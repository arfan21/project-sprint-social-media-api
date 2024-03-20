package post

import (
	"context"

	"github.com/arfan21/project-sprint-social-media-api/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, data entity.Post) (err error)
}
