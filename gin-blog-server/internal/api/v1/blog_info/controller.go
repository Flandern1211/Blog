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

func (ctrl *BlogInfoController) GetConfigMap(c *gin.Context) {
	data, err := ctrl.svc.GetConfigMap(c)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, data)
}

func (ctrl *BlogInfoController) UpdateConfigMap(c *gin.Context) {
	var m map[string]string
	if err := c.ShouldBindJSON(&m); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	if err := ctrl.svc.UpdateConfigMap(c, m); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}

type PageController struct {
	svc service.BlogInfoService
}

func NewPageController(svc service.BlogInfoService) *PageController {
	return &PageController{svc: svc}
}

func (ctrl *PageController) GetList(c *gin.Context) {
	data, _, err := ctrl.svc.GetPageList(c)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, data)
}

func (ctrl *PageController) SaveOrUpdate(c *gin.Context) {
	var req request.AddOrEditPageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	page, err := ctrl.svc.SaveOrUpdatePage(c, req)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, page)
}

func (ctrl *PageController) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	if err := ctrl.svc.DeletePages(c, ids); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}
