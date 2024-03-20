package server

import (
	fileuploaderctrl "github.com/arfan21/project-sprint-social-media-api/internal/fileuploader/controller"
	fileuploadersvc "github.com/arfan21/project-sprint-social-media-api/internal/fileuploader/service"
	userctrl "github.com/arfan21/project-sprint-social-media-api/internal/user/controller"
	userrepo "github.com/arfan21/project-sprint-social-media-api/internal/user/repository"
	usersvc "github.com/arfan21/project-sprint-social-media-api/internal/user/service"
	"github.com/arfan21/project-sprint-social-media-api/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Routes() {

	api := s.app.Group("")
	api.Get("/health-check", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	userRepo := userrepo.New(s.db)
	userSvc := usersvc.New(userRepo)
	userCtrl := userctrl.New(userSvc)

	fileUploaderSvc := fileuploadersvc.New()
	fileUploaderCtrl := fileuploaderctrl.New(fileUploaderSvc)

	s.RoutesCustomer(api, userCtrl)
	s.RoutesFileUploader(api, fileUploaderCtrl)
}

func (s Server) RoutesCustomer(route fiber.Router, ctrl *userctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	usersV1 := v1.Group("/user")
	usersV1.Post("/register", ctrl.Register)
	usersV1.Post("/login", ctrl.Login)

	friend := v1.Group("/friend", middleware.JWTAuth)
	friend.Post("", ctrl.AddFriend)
	friend.Delete("", ctrl.DeleteFriend)
	friend.Get("", ctrl.GetList)
}
func (s Server) RoutesFileUploader(route fiber.Router, ctrl *fileuploaderctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	fileUploaderV1 := v1.Group("/image", middleware.JWTAuth)
	fileUploaderV1.Post("", ctrl.UploadImage)
}
