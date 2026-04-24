package permission

import (
	"github.com/gin-gonic/gin"
)

func RegisterPermissionRouter(r *gin.RouterGroup, roleCtrl *RoleController, menuCtrl *MenuController, resourceCtrl *ResourceController) {
	// Role
	role := r.Group("/role")
	{
		role.GET("/option", roleCtrl.GetOption)
		role.GET("/list", roleCtrl.GetTreeList)
		role.POST("", roleCtrl.SaveOrUpdate)
		role.DELETE("", roleCtrl.Delete)
	}

	// Menu
	menu := r.Group("/menu")
	{
		menu.GET("/user", menuCtrl.GetUserMenu)
		menu.GET("/list", menuCtrl.GetTreeList)
		menu.POST("", menuCtrl.SaveOrUpdate)
		menu.DELETE("/:id", menuCtrl.Delete)
	}

	// Resource
	resource := r.Group("/resource")
	{
		resource.GET("/list", resourceCtrl.GetTreeList)
		resource.GET("/option", resourceCtrl.GetOption)
		resource.POST("", resourceCtrl.SaveOrUpdate)
		resource.DELETE("/:id", resourceCtrl.Delete)
		resource.PUT("/anonymous", resourceCtrl.UpdateAnonymous)
	}
}
