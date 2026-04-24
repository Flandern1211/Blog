package blog_info

import (
	"github.com/gin-gonic/gin"
)

func RegisterBlogInfoRouter(r *gin.RouterGroup, blogCtrl *BlogInfoController, pageCtrl *PageController) {
	// BlogInfo & Config
	blog := r.Group("/")
	{
		blog.GET("/home", blogCtrl.GetHomeInfo)
		blog.GET("/about", blogCtrl.GetAbout)
		blog.PUT("/about", blogCtrl.UpdateAbout)
		blog.POST("/report", blogCtrl.Report)

		blog.GET("/config", blogCtrl.GetConfigMap)
		blog.PUT("/config", blogCtrl.UpdateConfigMap)
	}

	// Page
	page := r.Group("/page")
	{
		page.GET("/list", pageCtrl.GetList)
		page.POST("", pageCtrl.SaveOrUpdate)
		page.DELETE("", pageCtrl.Delete)
	}
}
