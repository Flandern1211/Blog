package system

import (
	global "gin-blog/internal/global"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/service"

	"github.com/gin-gonic/gin"
)

type LinkController struct {
	svc service.SystemService
}

func NewLinkController(svc service.SystemService) *LinkController {
	return &LinkController{svc: svc}
}

func (ctrl *LinkController) GetList(c *gin.Context) {
	var query request.FriendLinkQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	list, total, err := ctrl.svc.GetLinkList(c, query)
	if err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnPageSuccess(c, list, total, query.Page, query.Size)
}

func (ctrl *LinkController) SaveOrUpdate(c *gin.Context) {
	var req request.AddOrEditLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	link, err := ctrl.svc.SaveOrUpdateLink(c, req)
	if err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, link)
}

func (ctrl *LinkController) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	if err := ctrl.svc.DeleteLinks(c, ids); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}

type OperationLogController struct {
	svc service.SystemService
}

func NewOperationLogController(svc service.SystemService) *OperationLogController {
	return &OperationLogController{svc: svc}
}

func (ctrl *OperationLogController) GetList(c *gin.Context) {
	var query request.OperationLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	list, total, err := ctrl.svc.GetOperationLogList(c, query)
	if err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnPageSuccess(c, list, total, query.Page, query.Size)
}

func (ctrl *OperationLogController) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	if err := ctrl.svc.DeleteOperationLogs(c, ids); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}
