package blog_info

import (
	"github.com/gin-gonic/gin"
)

func RegisterBlogInfoRouter(r *gin.RouterGroup, ctrl *BlogInfoController) {
	r.GET("/home", ctrl.GetHomeInfo)
	r.POST("/report", ctrl.Report)
}

func RegisterSettingRouter(r *gin.RouterGroup, ctrl *BlogInfoController) {
	setting := r.Group("/setting")
	{
		setting.GET("/about", ctrl.GetAbout)
		setting.PUT("/about", ctrl.UpdateAbout)
	}
}
