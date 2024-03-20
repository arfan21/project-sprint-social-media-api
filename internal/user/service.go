package user

import (
	"context"

	"github.com/arfan21/project-sprint-social-media-api/internal/model"
)

type Service interface {
	Register(ctx context.Context, req model.UserRegisterRequest) (res model.UserLoginResponse, err error)
	Login(ctx context.Context, req model.UserLoginRequest) (res model.UserLoginResponse, err error)
	AddFriend(ctx context.Context, req model.FriendRequest) (err error)
	DeleteFriend(ctx context.Context, req model.FriendRequest) (err error)
	GetList(ctx context.Context, req model.UserGetListRequest) (res []model.UserResponse, count int, err error)
}
