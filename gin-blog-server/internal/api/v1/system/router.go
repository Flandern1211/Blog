package system

import (
	"github.com/gin-gonic/gin"
)

func RegisterSystemRouter(r *gin.RouterGroup, linkCtrl *LinkController, optCtrl *OperationLogController) {
	// FriendLink
	link := r.Group("/link")
	{
		link.GET("/list", linkCtrl.GetList)
		link.POST("", linkCtrl.SaveOrUpdate)
		link.DELETE("", linkCtrl.Delete)
	}

	// OperationLog
	opt := r.Group("/operation/log")
	{
		opt.GET("/list", optCtrl.GetList)
		opt.DELETE("", optCtrl.Delete)
	}
}
