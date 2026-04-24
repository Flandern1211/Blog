package article

import (
	global "gin-blog/internal/global"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/service"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	svc service.ArticleService
}

func NewArticleController(svc service.ArticleService) *ArticleController {
	return &ArticleController{svc: svc}
}

func (ctrl *ArticleController) GetList(c *gin.Context) {
	var query request.ArticleQuery
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

func (ctrl *ArticleController) GetById(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	vo, err := ctrl.svc.GetById(c, id)
	if err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, vo)
}

func (ctrl *ArticleController) SaveOrUpdate(c *gin.Context) {
	var req request.AddOrEditArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	authId := sessions.Default(c).Get(global.CTX_USER_AUTH).(int)
	if err := ctrl.svc.SaveOrUpdate(c, authId, req); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}

func (ctrl *ArticleController) UpdateTop(c *gin.Context) {
	var req request.UpdateArticleTopReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	if err := ctrl.svc.UpdateTop(c, req); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}

func (ctrl *ArticleController) SoftDelete(c *gin.Context) {
	var req request.SoftDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	if err := ctrl.svc.SoftDelete(c, req); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}

func (ctrl *ArticleController) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	if err := ctrl.svc.Delete(c, ids); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}

type CategoryController struct {
	svc service.ArticleService
}

func NewCategoryController(svc service.ArticleService) *CategoryController {
	return &CategoryController{svc: svc}
}

func (ctrl *CategoryController) GetList(c *gin.Context) {
	var query request.CategoryQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	list, total, err := ctrl.svc.GetCategoryList(c, query)
	if err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnPageSuccess(c, list, total, query.Page, query.Size)
}

func (ctrl *CategoryController) SaveOrUpdate(c *gin.Context) {
	var req request.AddOrEditCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	if err := ctrl.svc.SaveOrUpdateCategory(c, req); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}

func (ctrl *CategoryController) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	if err := ctrl.svc.DeleteCategories(c, ids); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}

func (ctrl *CategoryController) GetOption(c *gin.Context) {
	list, err := ctrl.svc.GetCategoryOption(c)
	if err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, list)
}

type TagController struct {
	svc service.ArticleService
}

func NewTagController(svc service.ArticleService) *TagController {
	return &TagController{svc: svc}
}

func (ctrl *TagController) GetList(c *gin.Context) {
	var query request.TagQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	list, total, err := ctrl.svc.GetTagList(c, query)
	if err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnPageSuccess(c, list, total, query.Page, query.Size)
}

func (ctrl *TagController) SaveOrUpdate(c *gin.Context) {
	var req request.AddOrEditTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	if err := ctrl.svc.SaveOrUpdateTag(c, req); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}

func (ctrl *TagController) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	if err := ctrl.svc.DeleteTags(c, ids); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}

func (ctrl *TagController) GetOption(c *gin.Context) {
	list, err := ctrl.svc.GetTagOption(c)
	if err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, list)
}
