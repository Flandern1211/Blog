package permission

import (
	"gin-blog/internal/middleware"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/service"
	"gin-blog/pkg/errors"
	"gin-blog/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	svc service.PermissionService
}

func NewRoleController(svc service.PermissionService) *RoleController {
	return &RoleController{svc: svc}
}

func (ctrl *RoleController) GetOption(c *gin.Context) {
	list, err := ctrl.svc.GetRoleOption(c)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, list)
}

func (ctrl *RoleController) GetTreeList(c *gin.Context) {
	var query request.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	keyword := c.Query("keyword")

	list, total, err := ctrl.svc.GetRoleList(c, query, keyword)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.PageSuccess(c, list, total, query.Page, query.Size)
}

func (ctrl *RoleController) SaveOrUpdate(c *gin.Context) {
	var req request.AddOrEditRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	if err := ctrl.svc.SaveOrUpdateRole(c, req); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}

func (ctrl *RoleController) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	if err := ctrl.svc.DeleteRoles(c, ids); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}

type MenuController struct {
	svc service.PermissionService
}

func NewMenuController(svc service.PermissionService) *MenuController {
	return &MenuController{svc: svc}
}

func (ctrl *MenuController) GetUserMenu(c *gin.Context) {
	authId := middleware.GetUserID(c)
	if authId == 0 {
		response.Error(c, errors.CodeNoLogin, errors.GetMessage(errors.CodeNoLogin))
		return
	}

	isSuper := middleware.IsSuper(c)

	list, err := ctrl.svc.GetUserMenu(c, authId, isSuper)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, list)
}

func (ctrl *MenuController) GetTreeList(c *gin.Context) {
	keyword := c.Query("keyword")
	list, err := ctrl.svc.GetMenuTreeList(c, keyword)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, list)
}

func (ctrl *MenuController) GetOption(c *gin.Context) {
	list, err := ctrl.svc.GetMenuOption(c)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, list)
}

func (ctrl *MenuController) SaveOrUpdate(c *gin.Context) {
	var req request.AddOrEditMenuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	if err := ctrl.svc.SaveOrUpdateMenu(c, req); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}

func (ctrl *MenuController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctrl.svc.DeleteMenu(c, id); err != nil {
		response.BizError(c, err)
		return
	}
	response.Success(c, nil)
}

type ResourceController struct {
	svc service.PermissionService
}

func NewResourceController(svc service.PermissionService) *ResourceController {
	return &ResourceController{svc: svc}
}

func (ctrl *ResourceController) GetTreeList(c *gin.Context) {
	keyword := c.Query("keyword")
	list, err := ctrl.svc.GetResourceTreeList(c, keyword)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, list)
}

func (ctrl *ResourceController) GetOption(c *gin.Context) {
	list, err := ctrl.svc.GetResourceOption(c)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, list)
}

func (ctrl *ResourceController) SaveOrUpdate(c *gin.Context) {
	var req request.AddOrEditResourceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	if err := ctrl.svc.SaveOrUpdateResource(c, req); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}

func (ctrl *ResourceController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctrl.svc.DeleteResource(c, id); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}

func (ctrl *ResourceController) UpdateAnonymous(c *gin.Context) {
	var req request.EditAnonymousReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	if err := ctrl.svc.UpdateResourceAnonymous(c, req); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}
