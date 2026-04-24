package upload

import (
	global "gin-blog/internal/global"
	"gin-blog/internal/service"

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
		global.ReturnError(c, global.ErrRequest, err)
		return
	}

	url, err := ctrl.svc.UploadFile(c, file)
	if err != nil {
		global.ReturnError(c, global.ErrFileUpload, err)
		return
	}

	global.ReturnSuccess(c, url)
}
