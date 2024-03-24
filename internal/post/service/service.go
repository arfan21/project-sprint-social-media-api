package postsvc

import (
	"context"
	"fmt"

	"github.com/arfan21/project-sprint-social-media-api/internal/entity"
	"github.com/arfan21/project-sprint-social-media-api/internal/model"
	"github.com/arfan21/project-sprint-social-media-api/internal/post"
	"github.com/arfan21/project-sprint-social-media-api/internal/user"
	"github.com/arfan21/project-sprint-social-media-api/pkg/constant"
	"github.com/arfan21/project-sprint-social-media-api/pkg/logger"
	"github.com/arfan21/project-sprint-social-media-api/pkg/validation"
	"github.com/google/uuid"
)

type Service struct {
	repo    post.Repository
	userSvc user.Service
}

func New(repo post.Repository, userSvc user.Service) *Service {
	return &Service{repo: repo, userSvc: userSvc}
}

func (s Service) Create(ctx context.Context, req model.PostRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("post.service.Create: failed to validate request: %w", err)
		return
	}

	userIdUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		err = fmt.Errorf("post.service.Create: failed to parse user id: %w", err)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		err = fmt.Errorf("user.service.Create: failed to generate product id: %w", err)
		return
	}

	data := entity.Post{
		ID:     id,
		UserID: userIdUUID,
		Body:   req.PostInHtml,
		Tags:   req.Tags,
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("post.service.Create: failed to begin transaction: %w", err)
		return
	}

	defer func() {
		logger.Log(ctx).Err(err).Msg("post.service.Create: defer func rollback")
		if err != nil {
			errRb := tx.Rollback(ctx)
			if errRb != nil {
				err = fmt.Errorf("post.service.Create: failed  to rollback: %w", errRb)
				return
			}
			return
		}

		err = tx.Commit(ctx)
		if err != nil {
			err = fmt.Errorf("post.service.Create: failed  to commit: %w", err)
			return
		}
	}()

	err = s.repo.WithTx(tx).Create(ctx, data)
	if err != nil {
		err = fmt.Errorf("post.service.Create: failed to create post: %w", err)
		return
	}

	err = s.repo.WithTx(tx).IncrementCount(ctx)
	if err != nil {
		err = fmt.Errorf("post.service.Create: failed to increment count: %w", err)
		return
	}

	return
}

func (s Service) CreateComment(ctx context.Context, req model.PostCommentRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("post.service.CreateComment: failed to validate request: %w", err)
		return
	}

	postData, err := s.repo.GetByID(ctx, req.PostID)
	if err != nil {
		err = fmt.Errorf("post.service.CreateComment: failed to get post: %w", err)
		return
	}
	if req.UserID != postData.UserID.String() {
		isFriend, err := s.userSvc.IsFriend(ctx, req.UserID, postData.UserID.String())
		if err != nil {
			err = fmt.Errorf("post.service.CreateComment: failed to check is friend: %w", err)
			return err
		}

		if !isFriend {
			err = fmt.Errorf("post.service.CreateComment: user is not friend with post owner, %w", constant.ErrUserNotFriend)
			return err
		}
	}

	postIdUUID, err := uuid.Parse(req.PostID)
	if err != nil {
		err = fmt.Errorf("post.service.CreateComment: failed to parse post id: %w", err)
		return
	}

	userIdUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		err = fmt.Errorf("post.service.CreateComment: failed to parse user id: %w", err)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		err = fmt.Errorf("user.service.CreateComment: failed to generate product id: %w", err)
		return
	}

	data := entity.PostComment{
		ID:      id,
		PostID:  postIdUUID,
		Comment: req.Comment,
		UserID:  userIdUUID,
	}

	err = s.repo.CreateComment(ctx, data)
	if err != nil {
		err = fmt.Errorf("post.service.CreateComment: failed to create comment: %w", err)
		return
	}

	return
}

func (s Service) GetList(ctx context.Context, req model.PostGetListRequest) (res []model.PostListResponse, count int, err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("post.service.CreateComment: failed to validate request: %w", err)
		return
	}

	data, postIDs, userIDsUnique, err := s.repo.GetList(ctx, req)
	if err != nil {
		err = fmt.Errorf("post.service.GetList: failed to get list of post: %w", err)
		return
	}
	commentsMap, err := s.repo.GetCommentsByPostIDsMap(ctx, postIDs, userIDsUnique)
	if err != nil {
		err = fmt.Errorf("post.service.GetList: failed to get comments by post ids: %w", err)
		return
	}
	userIDs := make([]string, len(userIDsUnique))
	i := 0
	for k := range userIDsUnique {
		userIDs[i] = k
		i++
	}

	userMap, err := s.userSvc.GetListMap(ctx, model.UserGetListRequest{
		UserIDs:       userIDs,
		DisableOffset: true,
		DisableOrder:  true,
	})

	if err != nil {
		err = fmt.Errorf("post.service.GetList: failed to get list of user: %w", err)
		return
	}

	count, err = s.repo.GetCountList(ctx, req)
	if err != nil {
		err = fmt.Errorf("post.service.GetList: failed to get count list of post: %w", err)
		return
	}

	res = make([]model.PostListResponse, len(data))

	for i, v := range data {

		res[i] = model.PostListResponse{
			PostID: v.ID.String(),
			Post: model.PostResponse{
				PostInHtml: v.Body,
				Tags:       v.Tags,
				CreatedAt:  v.CreatedAt.Format(constant.TimeISO8601Format),
			},
			Creator: userMap[v.UserID.String()],
		}

		comments := commentsMap[v.ID.String()]
		res[i].Comments = make([]model.PostCommentResponse, len(comments))
		for j, comment := range comments {
			res[i].Comments[j] = model.PostCommentResponse{
				Comment:   comment.Comment,
				CreatedAt: comment.CreatedAt.Format(constant.TimeISO8601Format),
				Creator:   userMap[comment.UserID.String()],
			}
		}
	}

	return
}
