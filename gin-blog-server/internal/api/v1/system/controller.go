package system

import (
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/service"
	"gin-blog/pkg/errors"
	"gin-blog/pkg/response"

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
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	list, total, err := ctrl.svc.GetLinkList(c, query)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.PageSuccess(c, list, total, query.Page, query.Size)
}

func (ctrl *LinkController) SaveOrUpdate(c *gin.Context) {
	var req request.AddOrEditLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	link, err := ctrl.svc.SaveOrUpdateLink(c, req)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, link)
}

func (ctrl *LinkController) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	if err := ctrl.svc.DeleteLinks(c, ids); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}
