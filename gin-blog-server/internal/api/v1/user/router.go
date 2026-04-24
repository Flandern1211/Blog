package user

import (
	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.RouterGroup, ctrl *UserController) {
	r.GET("/info", ctrl.GetInfo)
	r.PUT("/current", ctrl.UpdateCurrent)
	r.PUT("/update", ctrl.Update)
	r.PUT("/disable", ctrl.UpdateDisable)
	r.GET("/list", ctrl.GetList)
}
