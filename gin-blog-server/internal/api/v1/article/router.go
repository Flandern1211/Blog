package article

import (
	"github.com/gin-gonic/gin"
)

func RegisterArticleRouter(r *gin.RouterGroup, ctrl *ArticleController) {
	// Article
	art := r.Group("/article")
	{
		art.GET("/list", ctrl.GetList)
		art.GET("/:id", ctrl.GetById)
		art.POST("", ctrl.SaveOrUpdate)
		art.PUT("/top", ctrl.UpdateTop)
		art.PUT("/soft-delete", ctrl.SoftDelete)
		art.DELETE("", ctrl.Delete)
	}
}
