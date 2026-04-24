package front

import (
	"github.com/gin-gonic/gin"
)

func RegisterFrontRouter(r *gin.RouterGroup, ctrl *FrontController) {
	r.GET("/home", ctrl.GetHomeInfo)

	article := r.Group("/article")
	{
		article.GET("/list", ctrl.GetArticleList)
		article.GET("/:id", ctrl.GetArticleInfo)
		article.GET("/archive", ctrl.GetArchiveList)
		article.GET("/search", ctrl.SearchArticle)
	}

	category := r.Group("/category")
	{
		category.GET("/list", ctrl.GetCategoryList)
	}

	tag := r.Group("/tag")
	{
		tag.GET("/list", ctrl.GetTagList)
	}

	message := r.Group("/message")
	{
		message.GET("/list", ctrl.GetMessageList)
	}

	comment := r.Group("/comment")
	{
		comment.GET("/list", ctrl.GetCommentList)
		comment.GET("/replies/:comment_id", ctrl.GetCommentReplyList)
	}

	// 需要登录的接口 (由外部控制中间件，这里只负责路由定义)
	// 或者在这里也定义 Group？
}

// 需要登录的路由
func RegisterFrontAuthRouter(r *gin.RouterGroup, ctrl *FrontController) {
	r.POST("/message", ctrl.SaveMessage)
	r.POST("/comment", ctrl.AddComment)
	r.GET("/comment/like/:comment_id", ctrl.LikeComment)
	r.GET("/article/like/:article_id", ctrl.LikeArticle)
}
