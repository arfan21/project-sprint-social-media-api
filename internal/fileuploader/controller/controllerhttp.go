package fileuploaderctrl

import (
	"errors"

	"github.com/arfan21/project-sprint-social-media-api/internal/fileuploader"
	"github.com/arfan21/project-sprint-social-media-api/internal/model"
	"github.com/arfan21/project-sprint-social-media-api/pkg/exception"
	"github.com/arfan21/project-sprint-social-media-api/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type ControllerHTTP struct {
	service fileuploader.Service
}

func New(service fileuploader.Service) *ControllerHTTP {
	return &ControllerHTTP{service: service}
}

// @Summary Upload Image
// @Description Upload image to s3
// @Tags Image Uploader
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/image [post]
func (ctrl ControllerHTTP) UploadImage(c *fiber.Ctx) error {
	var req model.FileUploaderImageRequest
	file, err := c.FormFile("file")
	if err != nil {
		if errors.Is(err, fasthttp.ErrMissingFile) {
			return c.Status(fiber.StatusBadRequest).JSON(pkgutil.HTTPResponse{
				Code:    fiber.StatusBadRequest,
				Message: "file required",
			})
		}

		exception.PanicIfNeeded(err)
	}

	req.File = file

	res, err := ctrl.service.UploadImage(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(res)
}
