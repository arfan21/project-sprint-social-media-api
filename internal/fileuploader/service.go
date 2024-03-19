package fileuploader

import (
	"context"

	"github.com/arfan21/project-sprint-social-media-api/internal/model"
)

type Service interface {
	UploadImage(ctx context.Context, req model.FileUploaderImageRequest) (res model.FileUploaderImageResponse, err error)
}
