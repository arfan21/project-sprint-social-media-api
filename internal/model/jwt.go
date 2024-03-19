package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	// UserID for middleware purpose
	UserID string `json:"-"`

	Name string `json:"name"`
	jwt.RegisteredClaims
}
