package entity

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID     `json:"id"`
	UserID    uuid.UUID     `json:"userId"`
	Body      string        `json:"body"`
	Tags      []string      `json:"tags"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Comments  []PostComment `json:"comments"`
}

func (Post) TableName() string {
	return "posts"
}

type PostComment struct {
	ID      uuid.UUID `json:"id"`
	PostID  uuid.UUID `json:"postId"`
	UserID  uuid.UUID `json:"userId"`
	Comment string    `json:"comment"`
}

func (PostComment) TableName() string {
	return "post_comments"
}
