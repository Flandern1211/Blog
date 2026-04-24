package blog_info

import (
	global "gin-blog/internal/global"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/service"

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
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, data)
}

func (ctrl *BlogInfoController) GetAbout(c *gin.Context) {
	data, err := ctrl.svc.GetAbout(c)
	if err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, data)
}

func (ctrl *BlogInfoController) UpdateAbout(c *gin.Context) {
	var req request.AboutReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	if err := ctrl.svc.UpdateAbout(c, req); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, req.Content)
}

func (ctrl *BlogInfoController) Report(c *gin.Context) {
	if err := ctrl.svc.Report(c); err != nil {
		global.ReturnError(c, global.ErrRedisOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}

func (ctrl *BlogInfoController) GetConfigMap(c *gin.Context) {
	data, err := ctrl.svc.GetConfigMap(c)
	if err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, data)
}

func (ctrl *BlogInfoController) UpdateConfigMap(c *gin.Context) {
	var m map[string]string
	if err := c.ShouldBindJSON(&m); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	if err := ctrl.svc.UpdateConfigMap(c, m); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
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
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, data)
}

func (ctrl *PageController) SaveOrUpdate(c *gin.Context) {
	var req request.AddOrEditPageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	page, err := ctrl.svc.SaveOrUpdatePage(c, req)
	if err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, page)
}

func (ctrl *PageController) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	if err := ctrl.svc.DeletePages(c, ids); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}
