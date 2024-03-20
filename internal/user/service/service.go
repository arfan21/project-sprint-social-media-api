package usersvc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arfan21/project-sprint-social-media-api/config"
	"github.com/arfan21/project-sprint-social-media-api/internal/entity"
	"github.com/arfan21/project-sprint-social-media-api/internal/model"
	"github.com/arfan21/project-sprint-social-media-api/internal/user"
	"github.com/arfan21/project-sprint-social-media-api/pkg/constant"
	"github.com/arfan21/project-sprint-social-media-api/pkg/validation"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

type Service struct {
	repo user.Repository
}

func New(repo user.Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) Register(ctx context.Context, req model.UserRegisterRequest) (res model.UserLoginResponse, err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("user.service.Register: failed to validate request: %w", err)
		return
	}

	cost := bcrypt.DefaultCost
	if config.Get().Bcrypt.Salt > 0 {
		cost = config.Get().Bcrypt.Salt
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), cost)
	if err != nil {
		err = fmt.Errorf("user.service.Register: failed to hash password: %w", err)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		err = fmt.Errorf("user.service.Create: failed to generate product id: %w", err)
		return
	}

	data := entity.User{
		ID:       id,
		Name:     req.Name,
		Password: string(hashedPassword),
	}

	if req.CredentialType == "email" {
		data.Email = null.StringFrom(req.CredentialValue)
	} else {
		data.Phone = null.StringFrom(req.CredentialValue)
	}

	err = s.repo.Create(ctx, data)
	if err != nil {
		err = fmt.Errorf("user.service.Register: failed to register user: %w", err)
		return
	}

	return s.login(data, true)
}

func (s Service) Login(ctx context.Context, req model.UserLoginRequest) (res model.UserLoginResponse, err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("user.service.Login: failed to validate request: %w", err)
		return
	}

	var data entity.User

	data, err = s.repo.GetByCredential(ctx, req.CredentialType, req.CredentialValue)
	if err != nil {
		err = fmt.Errorf("user.service.Login: failed to get user by phone: %w", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			err = constant.ErrUsernameOrPasswordInvalid
		}
		err = fmt.Errorf("user.service.Login: failed to compare password: %w", err)
		return
	}

	return s.login(data, false)
}

func (s Service) login(data entity.User, isRegister bool) (res model.UserLoginResponse, err error) {
	accessTokenExpire := time.Duration(config.Get().JWT.ExpireIn) * time.Second

	accessToken, err := s.CreateJWTWithExpiry(
		data.ID.String(),
		data.Name,
		config.Get().JWT.Secret,
		accessTokenExpire,
	)

	if err != nil {
		err = fmt.Errorf("user.service.login: failed to create access token: %w", err)
		return
	}
	res = model.UserLoginResponse{
		Name:        data.Name,
		AccessToken: accessToken,
	}

	if isRegister {
		res.Email = data.Email.Ptr()
		res.Phone = data.Phone.Ptr()
	} else {
		phone := data.Phone.ValueOrZero()
		email := data.Email.ValueOrZero()
		res.Email = &email
		res.Phone = &phone
	}

	return
}

func (s Service) CreateJWTWithExpiry(id, name, secret string, expiry time.Duration) (token string, err error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, model.JWTClaims{
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Get().Service.Name,
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	token, err = jwtToken.SignedString([]byte(secret))
	if err != nil {
		err = fmt.Errorf("usecase: failed to create jwt token: %w", err)
		return
	}

	return
}

func (s Service) AddFriend(ctx context.Context, req model.FriendRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("user.service.AddFriend: failed to validate request: %w", err)
		return
	}

	if req.UserID == req.UserIDAdder {
		err = constant.ErrFriendSelfAdding
		return
	}

	_, err = s.repo.GetByID(ctx, req.UserID)
	if err != nil {
		err = fmt.Errorf("user.service.AddFriend: failed to get user by id: %w", err)
		return
	}

	err = s.repo.AddFriend(ctx, req.UserIDAdder, req.UserID)
	if err != nil {
		err = fmt.Errorf("user.service.AddFriend: failed to add friend: %w", err)
		return
	}

	return
}

func (s Service) DeleteFriend(ctx context.Context, req model.FriendRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("user.service.DeleteFriend: failed to validate request: %w", err)
		return
	}

	if req.UserID == req.UserIDAdder {
		err = constant.ErrFriendSelfDeleting
		return
	}

	_, err = s.repo.GetByID(ctx, req.UserID)
	if err != nil {
		err = fmt.Errorf("user.service.DeleteFriend: failed to get user by id: %w", err)
		return
	}

	err = s.repo.DeleteFriend(ctx, req.UserIDAdder, req.UserID)
	if err != nil {
		err = fmt.Errorf("user.service.DeleteFriend: failed to delete friend: %w", err)
		return
	}

	return
}

func (s Service) GetList(ctx context.Context, req model.UserGetListRequest) (res []model.UserResponse, count int, err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("user.service.GetList: failed to validate request: %w", err)
		return
	}

	if req.Limit == 0 {
		req.Limit = 5
	}

	resDB, err := s.repo.GetList(ctx, req)
	if err != nil {
		err = fmt.Errorf("user.service.GetList: failed to get list user: %w", err)
		return
	}

	count, err = s.repo.GetCountList(ctx, req)
	if err != nil {
		err = fmt.Errorf("user.service.GetList: failed to get count list user: %w", err)
		return
	}

	res = make([]model.UserResponse, len(resDB))

	for i, v := range resDB {
		res[i] = model.UserResponse{
			UserID:      v.ID.String(),
			Name:        v.Name,
			ImageUrl:    v.ImageUrl.ValueOrZero(),
			FriendCount: v.FriendCount,
			CreatedAt:   v.CreatedAt.Format(constant.TimeISO8601Format),
		}
	}

	return
}
