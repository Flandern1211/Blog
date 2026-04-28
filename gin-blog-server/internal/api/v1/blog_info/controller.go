package blog_info

import (
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/service"
	"gin-blog/pkg/errors"
	"gin-blog/pkg/response"

	"github.com/gin-gonic/gin"
)

type BlogInfoController struct {
	svc service.BlogInfoService
}

func NewBlogInfoController(svc service.BlogInfoService) *BlogInfoController {
	return &BlogInfoController{svc: svc}
}

func (ctrl *BlogInfoController) GetHomeInfo(c *gin.Context) {
	data, err := ctrl.svc.GetHomeInfo(c)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, data)
}

func (ctrl *BlogInfoController) GetAbout(c *gin.Context) {
	data, err := ctrl.svc.GetAbout(c)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, data)
}

func (ctrl *BlogInfoController) UpdateAbout(c *gin.Context) {
	var req request.AboutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	if err := ctrl.svc.UpdateAbout(c, req); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, req.Content)
}

func (ctrl *BlogInfoController) Report(c *gin.Context) {
	if err := ctrl.svc.Report(c); err != nil {
		response.Error(c, errors.CodeRedisOpError, errors.GetMessage(errors.CodeRedisOpError))
		return
	}
	response.Success(c, nil)
}
