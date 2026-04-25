package upload

import (
	"gin-blog/internal/service"
	"gin-blog/pkg/errors"
	"gin-blog/pkg/response"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	svc service.UploadService
}

func NewUploadController(svc service.UploadService) *UploadController {
	return &UploadController{svc: svc}
}

func (ctrl *UploadController) UploadFile(c *gin.Context) {
	_, file, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}

	url, err := ctrl.svc.UploadFile(c, file)
	if err != nil {
		response.Error(c, errors.CodeFileUploadErr, errors.GetMessage(errors.CodeFileUploadErr))
		return
	}

	response.Success(c, url)
}
