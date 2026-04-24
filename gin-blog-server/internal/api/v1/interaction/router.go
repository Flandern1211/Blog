package interaction

import (
	"github.com/gin-gonic/gin"
)

func RegisterInteractionRouter(r *gin.RouterGroup, msgCtrl *MessageController, cmtCtrl *CommentController) {
	// Message
	msg := r.Group("/message")
	{
		msg.GET("/list", msgCtrl.GetList)
		msg.PUT("/review", msgCtrl.UpdateReview)
		msg.DELETE("", msgCtrl.Delete)
	}

	// Comment
	cmt := r.Group("/comment")
	{
		cmt.GET("/list", cmtCtrl.GetList)
		cmt.PUT("/review", cmtCtrl.UpdateReview)
		cmt.DELETE("", cmtCtrl.Delete)
	}
}
