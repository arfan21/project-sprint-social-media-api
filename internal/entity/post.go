package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type Post struct {
	ID        uuid.UUID             `json:"id"`
	UserID    uuid.UUID             `json:"userId"`
	Body      string                `json:"body"`
	Tags      []string              `json:"tags"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	Comments  []PostCommentNullable `json:"comments"`
	Total     int                   `json:"total"`
}

func (Post) TableName() string {
	return "posts"
}

type PostComment struct {
	ID        uuid.UUID `json:"id"`
	PostID    uuid.UUID `json:"postId"`
	UserID    uuid.UUID `json:"userId"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PostComment) TableName() string {
	return "post_comments"
}

type PostCommentNullable struct {
	ID        uuid.NullUUID `json:"id"`
	PostID    uuid.NullUUID `json:"postId"`
	UserID    uuid.NullUUID `json:"userId"`
	Comment   null.String   `json:"comment"`
	CreatedAt null.Time     `json:"created_at"`
	UpdatedAt null.Time     `json:"updated_at"`
}

type PostCounter struct {
	Count int `json:"count"`
}

func (PostCounter) TableName() string {
	return "post_counter"
}
