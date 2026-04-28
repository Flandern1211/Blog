package user

import (
	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.RouterGroup, ctrl *UserController) {
	user := r.Group("/user")
	{
		user.GET("/list", ctrl.GetList)
		user.PUT("", ctrl.Update)
		user.PUT("/disable", ctrl.UpdateDisable)
		user.GET("/info", ctrl.GetInfo)
		user.PUT("/current", ctrl.UpdateCurrent)
		user.GET("/online", ctrl.GetOnlineList)
		user.POST("/offline/:id", ctrl.ForceOffline)
	}
}
