package fileuploadersvc

import (
	"context"
	"fmt"

	"github.com/arfan21/project-sprint-social-media-api/config"
	"github.com/arfan21/project-sprint-social-media-api/internal/model"
	"github.com/arfan21/project-sprint-social-media-api/pkg/s3"
	"github.com/arfan21/project-sprint-social-media-api/pkg/validation"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) UploadImage(ctx context.Context, req model.FileUploaderImageRequest) (res model.FileUploaderImageResponse, err error) {
	fieldName := "file"

	err = validation.ValidateContentType(fieldName, req.File.Header.Get("Content-Type"), validation.WithValidateContentTypeImage())
	if err != nil {
		err = fmt.Errorf("fileuploader.service.Upload: failed to validate content type: %w", err)
		return
	}
	// 10KB
	minSize := 10 * 1024
	err = validation.ValidateFileSize(fieldName, req.File.Size, validation.WithValidateFileSizeMinSize(int64(minSize)))
	if err != nil {
		err = fmt.Errorf("imageuploader.service.Upload: failed to validate file size: %w", err)
		return
	}

	client, err := s3.New()
	if err != nil {
		err = fmt.Errorf("imageuploader.service.Upload: failed to create s3 client: %w", err)
		return
	}

	folder := "images"
	bucket := config.Get().S3.Bucket
	resStr, err := client.Upload(ctx, bucket, folder, req.File)
	if err != nil {
		err = fmt.Errorf("imageuploader.service.Upload: failed to upload file: %w", err)
		return
	}
	url := client.GetURL(bucket, resStr)

	res.ImageURL = url

	return res, nil
}
