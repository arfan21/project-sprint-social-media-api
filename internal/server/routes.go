package server

import (
	fileuploaderctrl "github.com/arfan21/project-sprint-social-media-api/internal/fileuploader/controller"
	fileuploadersvc "github.com/arfan21/project-sprint-social-media-api/internal/fileuploader/service"
	postctrl "github.com/arfan21/project-sprint-social-media-api/internal/post/controller"
	postrepo "github.com/arfan21/project-sprint-social-media-api/internal/post/repository"
	postsvc "github.com/arfan21/project-sprint-social-media-api/internal/post/service"
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

	postRepo := postrepo.New(s.db)
	postSvc := postsvc.New(postRepo, userSvc)
	postCtrl := postctrl.New(postSvc)

	s.RoutesCustomer(api, userCtrl)
	s.RoutesFileUploader(api, fileUploaderCtrl)
	s.RoutesPost(api, postCtrl)
}

func (s Server) RoutesCustomer(route fiber.Router, ctrl *userctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	usersV1 := v1.Group("/user")
	usersV1.Post("/register", ctrl.Register)
	usersV1.Post("/login", ctrl.Login)
	usersV1.Patch("", middleware.JWTAuth, ctrl.UpdateProfile)

	friend := v1.Group("/friend", middleware.JWTAuth)
	friend.Post("", ctrl.AddFriend)
	friend.Delete("", ctrl.DeleteFriend)
	friend.Get("", ctrl.GetList)

	linkV1 := usersV1.Group("/link", middleware.JWTAuth)
	linkV1.Post("/phone", ctrl.UpdatePhone)
	linkV1.Post("/", ctrl.UpdateEmail)
}
func (s Server) RoutesFileUploader(route fiber.Router, ctrl *fileuploaderctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	fileUploaderV1 := v1.Group("/image", middleware.JWTAuth)
	fileUploaderV1.Post("", ctrl.UploadImage)
}

func (s Server) RoutesPost(route fiber.Router, ctrl *postctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	postV1 := v1.Group("/post", middleware.JWTAuth)
	postV1.Post("", ctrl.Create)
	postV1.Post("/comment", ctrl.CreateComment)
	postV1.Get("", ctrl.GetList)
}
