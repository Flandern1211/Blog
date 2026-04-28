package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterAuthRouter(r *gin.RouterGroup, ctrl *AuthController) {
	r.POST("/login", ctrl.Login)
	r.POST("/register", ctrl.Register)
	r.POST("/code", ctrl.SendCode)
	r.GET("/logout", ctrl.Logout)

	auth := r.Group("/auth")
	{
		auth.GET("/verify", ctrl.VerifyCode)
	}
}

func RegisterAdminAuthRouter(r *gin.RouterGroup, ctrl *AuthController) {
	r.POST("/admin/login", ctrl.AdminLogin)
}
