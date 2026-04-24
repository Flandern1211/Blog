package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterAuthRouter(r *gin.RouterGroup, ctrl *AuthController) {
	r.POST("/login", ctrl.Login)
	r.POST("/register", ctrl.Register)
	r.GET("/logout", ctrl.Logout)
	r.GET("/verify", ctrl.VerifyCode)
}
