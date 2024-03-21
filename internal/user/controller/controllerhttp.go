package userctrl

import (
	"github.com/arfan21/project-sprint-social-media-api/internal/model"
	"github.com/arfan21/project-sprint-social-media-api/internal/user"
	"github.com/arfan21/project-sprint-social-media-api/pkg/constant"
	"github.com/arfan21/project-sprint-social-media-api/pkg/exception"
	"github.com/arfan21/project-sprint-social-media-api/pkg/logger"
	"github.com/arfan21/project-sprint-social-media-api/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
)

type ControllerHTTP struct {
	svc user.Service
}

func New(svc user.Service) *ControllerHTTP {
	return &ControllerHTTP{svc: svc}
}

// @Summary Register user
// @Description Register user
// @Tags user
// @Accept json
// @Produce json
// @Param body body model.UserRegisterRequest true "Payload user Register Request"
// @Success 201 {object} pkgutil.HTTPResponse{data=model.UserLoginResponse}
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/user/register [post]
func (ctrl ControllerHTTP) Register(c *fiber.Ctx) error {
	var req model.UserRegisterRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	res, err := ctrl.svc.Register(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusCreated).JSON(pkgutil.HTTPResponse{
		Message: "User registered successfully",
		Data:    res,
	})
}

// @Summary Login user
// @Description Login user
// @Tags user
// @Accept json
// @Produce json
// @Param body body model.UserLoginRequest true "Payload user Login Request"
// @Success 200 {object} pkgutil.HTTPResponse{data=model.UserLoginResponse}
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/user/login [post]
func (ctrl ControllerHTTP) Login(c *fiber.Ctx) error {
	var req model.UserLoginRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	res, err := ctrl.svc.Login(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Message: "User login successfully",
		Data:    res,
	})
}

// @Summary Add friend
// @Description Add friend
// @Tags friend
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param body body model.FriendRequest true "Payload friend request"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/friend [post]
func (ctrl ControllerHTTP) AddFriend(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	var req model.FriendRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	req.UserIDAdder = claims.UserID

	err = ctrl.svc.AddFriend(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Message: "Friend added successfully",
	})
}

// @Summary Delete friend
// @Description Delete friend
// @Tags friend
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param body body model.FriendRequest true "Payload friend request"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/friend [delete]
func (ctrl ControllerHTTP) DeleteFriend(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	var req model.FriendRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	req.UserIDAdder = claims.UserID

	err = ctrl.svc.DeleteFriend(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Message: "Friend deleted successfully",
	})
}

// @Summary Get list user
// @Description Get list user
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param limit query int false "Limit data"
// @Param offset query int false "Offset data"
// @Param search query string false "Search data"
// @Param sortBy query string false "Sort by data"
// @Param orderBy query string false "Order by data"
// @Param onlyFriend query bool false "Only friend data"
// @Success 200 {object} pkgutil.HTTPResponse{data=[]model.UserResponse}
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/friend [get]
func (ctrl ControllerHTTP) GetList(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	var req model.UserGetListRequest
	err := c.QueryParser(&req)
	exception.PanicIfNeeded(err)

	req.UserID = claims.UserID
	if req.Limit == 0 {
		req.Limit = 5
	}

	res, count, err := ctrl.svc.GetList(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Data: res,
		Meta: pkgutil.MetaResponse{
			Offset: req.Offset,
			Limit:  req.Limit,
			Total:  count,
		},
	})
}

// @Summary Update Phone
// @Description Update Phone
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param body body model.UserPhoneUpdateRequest true "Payload user update phone request"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/user/link/phone [post]
func (ctrl ControllerHTTP) UpdatePhone(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	var req model.UserPhoneUpdateRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	req.UserID = claims.UserID

	err = ctrl.svc.UpdatePhone(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Message: "Phone updated successfully",
	})
}

// @Summary Update Email
// @Description Update Email
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param body body model.UserEmailUpdateRequest true "Payload user update email request"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/user/link/email [post]
func (ctrl ControllerHTTP) UpdateEmail(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	var req model.UserEmailUpdateRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	req.UserID = claims.UserID

	err = ctrl.svc.UpdateEmail(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Message: "Email updated successfully",
	})
}

// @Summary Update Profile
// @Description Update Profile
// @Tags user
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param body body model.UserProfileUpdateRequest true "Payload user update profile request"
// @Success 200 {object} pkgutil.HTTPResponse{data=model.UserResponse}
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/user [patch]
func (ctrl ControllerHTTP) UpdateProfile(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	var req model.UserProfileUpdateRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	req.UserID = claims.UserID

	err = ctrl.svc.UpdateProfile(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Message: "Profile updated successfully",
	})
}
