package postsvc

import (
	"context"
	"fmt"

	"github.com/arfan21/project-sprint-social-media-api/internal/entity"
	"github.com/arfan21/project-sprint-social-media-api/internal/model"
	"github.com/arfan21/project-sprint-social-media-api/internal/post"
	"github.com/arfan21/project-sprint-social-media-api/pkg/validation"
	"github.com/google/uuid"
)

type Service struct {
	repo post.Repository
}

func New(repo post.Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) Create(ctx context.Context, req model.PostRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("post.service.Create: failed to validate request: %w", err)
		return
	}

	userIdUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		err = fmt.Errorf("post.service.Create: failed to parse user id: %w", err)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		err = fmt.Errorf("user.service.Create: failed to generate product id: %w", err)
		return
	}

	data := entity.Post{
		ID:     id,
		UserID: userIdUUID,
		Body:   req.PostInHtml,
		Tags:   req.Tags,
	}

	err = s.repo.Create(ctx, data)
	if err != nil {
		err = fmt.Errorf("post.service.Create: failed to create post: %w", err)
		return
	}

	return
}
