package user

import (
	global "gin-blog/internal/global"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	svc service.UserService
}

func NewUserController(svc service.UserService) *UserController {
	return &UserController{svc: svc}
}

func (ctrl *UserController) GetInfo(c *gin.Context) {
	authId := sessions.Default(c).Get(global.CTX_USER_AUTH).(int)
	vo, err := ctrl.svc.GetInfo(c, authId)
	if err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, vo)
}

func (ctrl *UserController) UpdateCurrent(c *gin.Context) {
	var req request.UpdateCurrentUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}

	authId := sessions.Default(c).Get(global.CTX_USER_AUTH).(int)
	if err := ctrl.svc.UpdateCurrent(c, authId, req); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}

func (ctrl *UserController) Update(c *gin.Context) {
	var req request.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}

	if err := ctrl.svc.Update(c, req); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}

func (ctrl *UserController) UpdateDisable(c *gin.Context) {
	var req request.UpdateUserDisableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}

	if err := ctrl.svc.UpdateDisable(c, req); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}

func (ctrl *UserController) GetList(c *gin.Context) {
	var query request.UserQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}

	list, total, err := ctrl.svc.GetList(c, query)
	if err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnPageSuccess(c, list, total, query.Page, query.Size)
}
