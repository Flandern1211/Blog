package user

import (
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/service"
	"gin-blog/pkg/errors"
	"gin-blog/pkg/response"

	g "gin-blog/internal/global"
	"gi
	"github.com/gin-contrib/sessions"
	g "gin-blog/internal/global"
)

type UserController struct {
	svc service.UserService
}

func NewUserController(svc service.UserService) *UserController {
	return &UserController{svc: svc}
}

func (ctrl *UserController) GetInfo(c *gin.Context) {
	val := sessions.Default(c).Get(g.CTX_USER_AUTH)
	if val == nil {
		response.Error(c, errors.CodeNoLogin, errors.GetMessage(errors.CodeNoLogin))
		return
	}
	authId := val.(int)
	vo, err := ctrl.svc.GetInfo(c, authId)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, vo)
}

func (ctrl *UserController) UpdateCurrent(c *gin.Context) {
	var req request.UpdateCurrentUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}

	authIdVal := sessions.Default(c).Get(g.CTX_USER_AUTH)
	if authIdVal == nil {
		response.Error(c, errors.CodeNoLogin, errors.GetMessage(errors.CodeNoLogin))
		return
	}
	authId := authIdVal.(int)
	if err := ctrl.svc.UpdateCurrent(c, authId, req); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}

func (ctrl *UserController) Update(c *gin.Context) {
	var req request.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}

	if err := ctrl.svc.Update(c, req); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}

func (ctrl *UserController) UpdateDisable(c *gin.Context) {
	var req request.UpdateUserDisableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}

	if err := ctrl.svc.UpdateDisable(c, req); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}

func (ctrl *UserController) GetList(c *gin.Context) {
	var query request.UserQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}

	list, total, err := ctrl.svc.GetList(c, query)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.PageSuccess(c, list, total, query.Page, query.Size)
}
