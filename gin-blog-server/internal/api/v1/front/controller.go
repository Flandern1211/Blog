package front

import (
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/service"
	"gin-blog/pkg/errors"
	"gin-blog/pkg/response"
	"html/template"
	"strconv"

	"gin-blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

type FrontController struct {
	svc         service.FrontService
	articleSvc  service.ArticleService
	interactSvc service.InteractionService
	blogInfoSvc service.BlogInfoService
	systemSvc   service.SystemService
}

func NewFrontController(svc service.FrontService, articleSvc service.ArticleService, interactSvc service.InteractionService, blogInfoSvc service.BlogInfoService, systemSvc service.SystemService) *FrontController {
	return &FrontController{
		svc:         svc,
		articleSvc:  articleSvc,
		interactSvc: interactSvc,
		blogInfoSvc: blogInfoSvc,
		systemSvc:   systemSvc,
	}
}

// BlogInfo
func (ctrl *FrontController) GetHomeInfo(c *gin.Context) {
	data, err := ctrl.svc.GetHomeInfo(c)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, data)
}

// Article
func (ctrl *FrontController) GetArticleList(c *gin.Context) {
	var query request.FArticleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	list, _, err := ctrl.articleSvc.GetBlogArticleList(c, query)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, list)
}

func (ctrl *FrontController) GetArticleInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}

	article, err := ctrl.articleSvc.GetBlogArticle(c, id)
	if err != nil {
		if err == errors.ErrNotFound {
			response.Error(c, errors.CodeNotFound, errors.GetMessage(errors.CodeNotFound))
			return
		}
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, article)
}

func (ctrl *FrontController) GetArchiveList(c *gin.Context) {
	var query request.FArticleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	list, total, err := ctrl.articleSvc.GetBlogArticleList(c, query)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	// Return list of archives (id, title, created_at)
	type ArchiveVO struct {
		ID        int    `json:"id"`
		Title     string `json:"title"`
		CreatedAt string `json:"created_at"`
	}
	var archives []ArchiveVO
	for _, a := range list {
		archives = append(archives, ArchiveVO{
			ID:        a.ID,
			Title:     a.Title,
			CreatedAt: a.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	response.PageSuccess(c, archives, total, query.Page, query.Size)
}

func (ctrl *FrontController) SearchArticle(c *gin.Context) {
	keyword := c.Query("keyword")
	list, err := ctrl.svc.SearchArticle(c, keyword)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, list)
}

func (ctrl *FrontController) LikeArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("article_id"))
	if err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	user, err := middleware.CurrentUserAuth(c)
	if err != nil {
		response.Error(c, errors.CodeNoLogin, errors.GetMessage(errors.CodeNoLogin))
		return
	}
	authId := user.ID
	if err := ctrl.svc.LikeArticle(c, id, authId); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}

// Category
func (ctrl *FrontController) GetCategoryList(c *gin.Context) {
	list, _, err := ctrl.articleSvc.GetCategoryList(c, request.CategoryQuery{PageQuery: request.PageQuery{Page: 1, Size: 1000}})
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, list)
}

// Tag
func (ctrl *FrontController) GetTagList(c *gin.Context) {
	list, _, err := ctrl.articleSvc.GetTagList(c, request.TagQuery{PageQuery: request.PageQuery{Page: 1, Size: 1000}})
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, list)
}

// Message
func (ctrl *FrontController) GetMessageList(c *gin.Context) {
	isReview := true
	list, _, err := ctrl.interactSvc.GetMessageList(c, request.MessageQuery{PageQuery: request.PageQuery{Page: 1, Size: 1000}, IsReview: &isReview})
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, list)
}

func (ctrl *FrontController) SaveMessage(c *gin.Context) {
	var req request.FAddMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	req.Content = template.HTMLEscapeString(req.Content)

	user, err := middleware.CurrentUserAuth(c)
	if err != nil {
		response.Error(c, errors.CodeNoLogin, errors.GetMessage(errors.CodeNoLogin))
		return
	}

	err = ctrl.interactSvc.AddMessage(c, user.ID, req)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}

// Comment
func (ctrl *FrontController) GetCommentList(c *gin.Context) {
	var query request.FCommentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	list, total, err := ctrl.interactSvc.GetFrontCommentList(c, query)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.PageSuccess(c, list, total, query.Page, query.Size)
}

func (ctrl *FrontController) GetCommentReplyList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	var query request.PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}

	list, err := ctrl.interactSvc.GetCommentReplyList(c, id, query.Page, query.Size)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, list)
}

func (ctrl *FrontController) AddComment(c *gin.Context) {
	var req request.FAddCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}

	req.Content = template.HTMLEscapeString(req.Content)

	user, err := middleware.CurrentUserAuth(c)
	if err != nil {
		response.Error(c, errors.CodeNoLogin, errors.GetMessage(errors.CodeNoLogin))
		return
	}

	err = ctrl.interactSvc.AddComment(c, user.ID, req)
	if err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}

func (ctrl *FrontController) LikeComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		response.Error(c, errors.CodeRequestError, errors.GetMessage(errors.CodeRequestError))
		return
	}
	user, err := middleware.CurrentUserAuth(c)
	if err != nil {
		response.Error(c, errors.CodeNoLogin, errors.GetMessage(errors.CodeNoLogin))
		return
	}
	authId := user.ID
	if err := ctrl.svc.LikeComment(c, id, authId); err != nil {
		response.Error(c, errors.CodeDbOpError, errors.GetMessage(errors.CodeDbOpError))
		return
	}
	response.Success(c, nil)
}
