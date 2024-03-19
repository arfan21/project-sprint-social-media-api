package entity

import (
	"time"

	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

type User struct {
	ID        uuid.UUID   `json:"id"`
	Email     null.String `json:"email"`
	Phone     null.String `json:"phone"`
	Name      string      `json:"name"`
	Password  string      `json:"password"`
	ImageUrl  null.String `json:"imageUrl"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}

type Friend struct {
	UserIDAdder uuid.UUID `json:"userIdAdder"`
	UserIDAdded uuid.UUID `json:"userIdAdded"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdateAt    time.Time `json:"updatedAt"`
}

func (Friend) TableName() string {
	return "friends"
}
