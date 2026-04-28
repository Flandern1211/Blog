package article

import (
	"gin-blog/internal/middleware"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/service"
	"gin-blog/pkg/errors"
	"gin-blog/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	svc service.ArticleService
}

func NewArticleController(svc service.ArticleService) *ArticleController {
	return &ArticleController{svc: svc}
}

// 获取文章列表
func (ctrl *ArticleController) GetList(c *gin.Context) {
	var query request.ArticleQuery
	//参数绑定
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

func (ctrl *ArticleController) GetById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	vo, err := ctrl.svc.GetById(c, id)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, vo)
}

func (ctrl *ArticleController) SaveOrUpdate(c *gin.Context) {
	var req request.AddOrEditArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	authId := middleware.GetUserID(c)
	if authId == 0 {
		response.Error(c, errors.CodeNoLogin, errors.GetMessage(errors.CodeNoLogin))
		return
	}
	if err := ctrl.svc.SaveOrUpdate(c, authId, req); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}

func (ctrl *ArticleController) UpdateTop(c *gin.Context) {
	var req request.UpdateArticleTopReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	if err := ctrl.svc.UpdateTop(c, req); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}

func (ctrl *ArticleController) SoftDelete(c *gin.Context) {
	var req request.SoftDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	if err := ctrl.svc.SoftDelete(c, req); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}

func (ctrl *ArticleController) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	if err := ctrl.svc.Delete(c, ids); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}
