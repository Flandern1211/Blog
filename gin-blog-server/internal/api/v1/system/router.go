package system

import (
	"github.com/gin-gonic/gin"
)

func RegisterLinkRouter(r *gin.RouterGroup, ctrl *LinkController) {
	link := r.Group("/link")
	{
		link.GET("/list", ctrl.GetList)
		link.POST("", ctrl.SaveOrUpdate)
		link.DELETE("", ctrl.Delete)
	}
}
