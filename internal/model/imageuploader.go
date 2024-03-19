package model

import "mime/multipart"

type FileUploaderImageRequest struct {
	File *multipart.FileHeader `json:"file" form:"file"`
}

type FileUploaderImageResponse struct {
	ImageURL string `json:"imageUrl"`
}
