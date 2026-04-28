package user

import (
	"strconv"

	"gin-blog/internal/middleware"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/service"
	"gin-blog/pkg/errors"
	"gin-blog/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	svc service.UserService
}

func NewUserController(svc service.UserService) *UserController {
	return &UserController{svc: svc}
}

func (ctrl *UserController) GetInfo(c *gin.Context) {
	authId := middleware.GetUserID(c)
	if authId == 0 {
		response.Error(c, errors.CodeNoLogin, errors.GetMessage(errors.CodeNoLogin))
		return
	}
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

	authId := middleware.GetUserID(c)
	if authId == 0 {
		response.Error(c, errors.CodeNoLogin, errors.GetMessage(errors.CodeNoLogin))
		return
	}
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

// 前台用户通过验证码修改密码
func (ctrl *UserController) UpdatePasswordByCode(c *gin.Context) {
	var req request.UpdatePasswordByCodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}

	authId := middleware.GetUserID(c)
	if authId == 0 {
		response.Error(c, errors.CodeNoLogin, errors.GetMessage(errors.CodeNoLogin))
		return
	}

	if err := ctrl.svc.UpdatePasswordByCode(c, authId, req); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

func (ctrl *UserController) GetOnlineList(c *gin.Context) {
	keyword := c.Query("keyword")
	list, err := ctrl.svc.GetOnlineList(c, keyword)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, list)
}

func (ctrl *UserController) ForceOffline(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}

	authId := middleware.GetUserID(c)
	if authId == 0 {
		response.Error(c, errors.CodeNoLogin, errors.GetMessage(errors.CodeNoLogin))
		return
	}

	if authId == id {
		response.Error(c, errors.CodeForceOfflineSelf, errors.GetMessage(errors.CodeForceOfflineSelf))
		return
	}

	if err := ctrl.svc.ForceOffline(c, authId, id); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}
