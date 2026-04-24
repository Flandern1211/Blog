package article

import (
	"github.com/gin-gonic/gin"
)

func RegisterArticleRouter(r *gin.RouterGroup, artCtrl *ArticleController, catCtrl *CategoryController, tagCtrl *TagController) {
	// Article
	art := r.Group("/article")
	{
		art.GET("/list", artCtrl.GetList)
		art.GET("/:id", artCtrl.GetById)
		art.POST("", artCtrl.SaveOrUpdate)
		art.PUT("/top", artCtrl.UpdateTop)
		art.PUT("/soft-delete", artCtrl.SoftDelete)
		art.DELETE("", artCtrl.Delete)
	}

	// Category
	cat := r.Group("/category")
	{
		cat.GET("/list", catCtrl.GetList)
		cat.POST("", catCtrl.SaveOrUpdate)
		cat.DELETE("", catCtrl.Delete)
		cat.GET("/option", catCtrl.GetOption)
	}

	// Tag
	tag := r.Group("/tag")
	{
		tag.GET("/list", tagCtrl.GetList)
		tag.POST("", tagCtrl.SaveOrUpdate)
		tag.DELETE("", tagCtrl.Delete)
		tag.GET("/option", tagCtrl.GetOption)
	}
}
