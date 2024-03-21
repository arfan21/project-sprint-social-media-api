package postctrl

import (
	"github.com/arfan21/project-sprint-social-media-api/internal/model"
	"github.com/arfan21/project-sprint-social-media-api/internal/post"
	"github.com/arfan21/project-sprint-social-media-api/pkg/constant"
	"github.com/arfan21/project-sprint-social-media-api/pkg/exception"
	"github.com/arfan21/project-sprint-social-media-api/pkg/logger"
	"github.com/arfan21/project-sprint-social-media-api/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
)

type ControllerHTTP struct {
	svc post.Service
}

func New(svc post.Service) *ControllerHTTP {
	return &ControllerHTTP{svc: svc}
}

// @Summary Create post
// @Description Create post
// @Tags post
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param body body model.PostRequest true "Payload post request"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/post [post]
func (ctrl ControllerHTTP) Create(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	var req model.PostRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	req.UserID = claims.UserID

	err = ctrl.svc.Create(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.JSON(pkgutil.HTTPResponse{
		Message: "Post created successfully",
	})
}

// @Summary Create comment
// @Description Create comment
// @Tags post
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param body body model.PostCommentRequest true "Payload post comment request"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/post/comment [post]
func (ctrl ControllerHTTP) CreateComment(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	var req model.PostCommentRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	req.UserID = claims.UserID

	err = ctrl.svc.CreateComment(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.JSON(pkgutil.HTTPResponse{
		Message: "Comment created successfully",
	})
}

// @Summary Get list post
// @Description Get list post
// @Tags post
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param limit query int false "Limit data"
// @Param offset query int false "Offset data"
// @Param search query string false "Search data"
// @Param searchTag query string false "Search tag data"
// @Success 200 {object} pkgutil.HTTPResponse{data=[]model.PostListResponse}
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/post [get]
func (ctrl ControllerHTTP) GetList(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	var req model.PostGetListRequest
	err := c.QueryParser(&req)
	exception.PanicIfNeeded(err)

	req.UserID = claims.UserID
	if req.Limit == 0 {
		req.Limit = 5
	}

	data, count, err := ctrl.svc.GetList(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.JSON(pkgutil.HTTPResponse{
		Data: data,
		Meta: pkgutil.MetaResponse{
			Offset: req.Offset,
			Limit:  req.Limit,
			Total:  count,
		},
	})
}
