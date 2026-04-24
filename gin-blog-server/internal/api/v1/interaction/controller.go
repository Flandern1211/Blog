package interaction

import (
	global "gin-blog/internal/global"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/service"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	svc service.InteractionService
}

func NewMessageController(svc service.InteractionService) *MessageController {
	return &MessageController{svc: svc}
}

func (ctrl *MessageController) GetList(c *gin.Context) {
	var query request.MessageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	list, total, err := ctrl.svc.GetMessageList(c, query)
	if err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnPageSuccess(c, list, total, query.Page, query.Size)
}

func (ctrl *MessageController) UpdateReview(c *gin.Context) {
	var req request.UpdateReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	if err := ctrl.svc.UpdateMessagesReview(c, req); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}

func (ctrl *MessageController) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	if err := ctrl.svc.DeleteMessages(c, ids); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}

type CommentController struct {
	svc service.InteractionService
}

func NewCommentController(svc service.InteractionService) *CommentController {
	return &CommentController{svc: svc}
}

func (ctrl *CommentController) GetList(c *gin.Context) {
	var query request.CommentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	list, total, err := ctrl.svc.GetCommentList(c, query)
	if err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnPageSuccess(c, list, total, query.Page, query.Size)
}

func (ctrl *CommentController) UpdateReview(c *gin.Context) {
	var req request.UpdateReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	if err := ctrl.svc.UpdateCommentsReview(c, req); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}

func (ctrl *CommentController) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		global.ReturnError(c, global.ErrRequest, err)
		return
	}
	if err := ctrl.svc.DeleteComments(c, ids); err != nil {
		global.ReturnError(c, global.ErrDbOp, err)
		return
	}
	global.ReturnSuccess(c, nil)
}
