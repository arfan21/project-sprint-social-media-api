package model

type PostRequest struct {
	PostInHtml string   `json:"postInHtml" validate:"required,min=3"`
	Tags       []string `json:"tags" validate:"required,dive,required"`
	UserID     string   `json:"-" validate:"required"`
}
