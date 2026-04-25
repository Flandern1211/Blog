package ginblog

import (
	"gin-blog/docs"
	"gin-blog/internal/api/v1/article"
	"gin-blog/internal/api/v1/auth"
	"gin-blog/internal/api/v1/blog_info"
	"gin-blog/internal/api/v1/front"
	"gin-blog/internal/api/v1/interaction"
	"gin-blog/internal/api/v1/permission"
	"gin-blog/internal/api/v1/system"
	"gin-blog/internal/api/v1/upload"
	"gin-blog/internal/api/v1/user"
	"gin-blog/internal/middleware"
	"gin-blog/internal/repository"
	"gin-blog/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	// 后台管理系统接口 (新 MVC 模式)
	categoryCtrl *article.CategoryController
	tagCtrl      *article.TagController
	articleCtrl  *article.ArticleController

	uploadCtrl *upload.UploadController
	userCtrl   *user.UserController
	authCtrl   *auth.AuthController

	commentCtrl *interaction.CommentController
	messageCtrl *interaction.MessageController

	roleCtrl     *permission.RoleController
	resourceCtrl *permission.ResourceController
	menuCtrl     *permission.MenuController

	blogInfoCtrl *blog_info.BlogInfoController
	linkCtrl     *system.LinkController
	logCtrl      *system.OperationLogController

	// 博客前台接口 (新 MVC 模式)
	frontCtrl *front.FrontController
)

func init() {
	// 初始化仓储
	articleRepo := repository.NewArticleRepository()
	authRepo := repository.NewAuthRepository()
	userRepo := repository.NewUserRepository()
	interactRepo := repository.NewInteractionRepository()
	blogInfoRepo := repository.NewBlogInfoRepository()
	systemRepo := repository.NewSystemRepository()
	permissionRepo := repository.NewPermissionRepository()

	// 初始化服务
	articleSvc := service.NewArticleService(articleRepo, interactRepo)
	authSvc := service.NewAuthService(authRepo)
	userSvc := service.NewUserService(userRepo)
	interactSvc := service.NewInteractionService(interactRepo, blogInfoRepo)
	blogInfoSvc := service.NewBlogInfoService(blogInfoRepo)
	systemSvc := service.NewSystemService(systemRepo)
	permissionSvc := service.NewPermissionService(permissionRepo)
	frontSvc := service.NewFrontService()

	// 初始化控制器
	articleCtrl = article.NewArticleController(articleSvc)
	categoryCtrl = article.NewCategoryController(articleSvc)
	tagCtrl = article.NewTagController(articleSvc)

	authCtrl = auth.NewAuthController(authSvc)
	uploadCtrl = upload.NewUploadController(service.NewUploadService())
	userCtrl = user.NewUserController(userSvc)

	commentCtrl = interaction.NewCommentController(interactSvc)
	messageCtrl = interaction.NewMessageController(interactSvc)

	roleCtrl = permission.NewRoleController(permissionSvc)
	resourceCtrl = permission.NewResourceController(permissionSvc)
	menuCtrl = permission.NewMenuController(permissionSvc)

	blogInfoCtrl = blog_info.NewBlogInfoController(blogInfoSvc)
	linkCtrl = system.NewLinkController(systemSvc)
	logCtrl = system.NewOperationLogController(systemSvc)

	frontCtrl = front.NewFrontController(frontSvc, articleSvc, interactSvc, blogInfoSvc, systemSvc)
}

// TODO: 前端修改 PUT 和 PATCH 请求
func RegisterHandlers(r *gin.Engine) {
	// Swagger
	docs.SwaggerInfo.BasePath = "/api"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	registerBaseHandler(r)
	registerAdminHandler(r)
	registerBlogHandler(r)
}

// 通用接口: 全部不需要 登录 + 鉴权
func registerBaseHandler(r *gin.Engine) {
	base := r.Group("/api")

	base.POST("/login", authCtrl.Login)            // 登录
	base.POST("/register", authCtrl.Register)      // 注册
	base.POST("/code", authCtrl.SendCode)          // 发送邮箱验证码
	base.GET("/logout", authCtrl.Logout)           // 退出登录
	base.GET("/config", blogInfoCtrl.GetConfigMap) // 获取配置
}

// 后台管理系统的接口: 全部需要 登录 + 鉴权
func registerAdminHandler(r *gin.Engine) {
	auth := r.Group("/api")

	// !注意使用中间件的顺序
	auth.Use(middleware.JWTAuth())
	auth.Use(middleware.PermissionCheck())
	auth.Use(middleware.OperationLog())
	auth.Use(middleware.ListenOnline())

	auth.GET("/home", blogInfoCtrl.GetHomeInfo)         // 后台首页信息
	auth.POST("/upload", uploadCtrl.UploadFile)         // 文件上传
	auth.PATCH("/config", blogInfoCtrl.UpdateConfigMap) // 更新配置 (Issue #6)

	// 博客设置
	setting := auth.Group("/setting")
	{
		setting.GET("/about", blogInfoCtrl.GetAbout)    // 获取关于我
		setting.PUT("/about", blogInfoCtrl.UpdateAbout) // 编辑关于我
	}
	// 用户模块
	user := auth.Group("/user")
	{
		user.GET("/list", userCtrl.GetList)          // 用户列表
		user.PUT("", userCtrl.Update)                // 更新用户信息
		user.PUT("/disable", userCtrl.UpdateDisable) // 修改用户禁用状态
		// user.PUT("/password", userCtrl.UpdatePassword)                // 修改普通用户密码
		// user.PUT("/current/password", userCtrl.UpdateCurrentPassword) // 修改管理员密码
		user.GET("/info", userCtrl.GetInfo)          // 获取当前用户信息
		user.PUT("/current", userCtrl.UpdateCurrent) // 修改当前用户信息
		// user.GET("/online", userCtrl.GetOnlineList)                   // 获取在线用户
		// user.POST("/offline/:id", userCtrl.ForceOffline)              // 强制用户下线
	}
	// 分类模块
	category := auth.Group("/category")
	{
		category.GET("/list", categoryCtrl.GetList)     // 分类列表
		category.POST("", categoryCtrl.SaveOrUpdate)    // 新增/编辑分类
		category.DELETE("", categoryCtrl.Delete)        // 删除分类
		category.GET("/option", categoryCtrl.GetOption) // 分类选项列表
	}
	// 标签模块
	tag := auth.Group("/tag")
	{
		tag.GET("/list", tagCtrl.GetList)     // 标签列表
		tag.POST("", tagCtrl.SaveOrUpdate)    // 新增/编辑标签
		tag.DELETE("", tagCtrl.Delete)        // 删除标签
		tag.GET("/option", tagCtrl.GetOption) // 标签选项列表
	}
	// 文章模块
	articles := auth.Group("/article")
	{
		articles.GET("/list", articleCtrl.GetList)           // 文章列表
		articles.POST("", articleCtrl.SaveOrUpdate)          // 新增/编辑文章
		articles.PUT("/top", articleCtrl.UpdateTop)          // 更新文章置顶
		articles.GET("/:id", articleCtrl.GetById)            // 文章详情
		articles.PUT("/soft-delete", articleCtrl.SoftDelete) // 软删除文章
		articles.DELETE("", articleCtrl.Delete)              // 物理删除文章
		// articles.POST("/export", articleCtrl.Export)               // 导出文章
		// articles.POST("/import", articleCtrl.Import)               // 导入文章
	}
	// 评论模块
	comment := auth.Group("/comment")
	{
		comment.GET("/list", commentCtrl.GetList)        // 评论列表
		comment.DELETE("", commentCtrl.Delete)           // 删除评论
		comment.PUT("/review", commentCtrl.UpdateReview) // 修改评论审核
	}
	// 留言模块
	message := auth.Group("/message")
	{
		message.GET("/list", messageCtrl.GetList)        // 留言列表
		message.DELETE("", messageCtrl.Delete)           // 删除留言
		message.PUT("/review", messageCtrl.UpdateReview) // 审核留言
	}
	// 资源模块
	resource := auth.Group("/resource")
	{
		resource.GET("/list", resourceCtrl.GetTreeList)          // 资源列表(树形)
		resource.POST("", resourceCtrl.SaveOrUpdate)             // 新增/编辑资源
		resource.DELETE("/:id", resourceCtrl.Delete)             // 删除资源
		resource.PUT("/anonymous", resourceCtrl.UpdateAnonymous) // 修改资源匿名访问
		resource.GET("/option", resourceCtrl.GetOption)          // 资源选项列表(树形)
	}
	// 菜单模块
	menu := auth.Group("/menu")
	{
		menu.GET("/list", menuCtrl.GetTreeList)      // 菜单列表
		menu.POST("", menuCtrl.SaveOrUpdate)         // 新增/编辑菜单
		menu.DELETE("/:id", menuCtrl.Delete)         // 删除菜单
		menu.GET("/user/list", menuCtrl.GetUserMenu) // 获取当前用户的菜单
		menu.GET("/option", menuCtrl.GetOption)      // 菜单选项列表(树形)
	}
	// 角色模块
	role := auth.Group("/role")
	{
		role.GET("/list", roleCtrl.GetTreeList) // 角色列表(树形)
		role.POST("", roleCtrl.SaveOrUpdate)    // 新增/编辑菜单
		role.DELETE("", roleCtrl.Delete)        // 删除角色
		role.GET("/option", roleCtrl.GetOption) // 角色选项列表(树形)
	}
	// 操作日志模块
	operationLog := auth.Group("/operation/log")
	{
		operationLog.GET("/list", logCtrl.GetList) // 操作日志列表
		operationLog.DELETE("", logCtrl.Delete)    // 删除操作日志
	}
	// 页面模块
	// page := auth.Group("/page")
	// {
	// 	page.GET("/list", pageCtrl.GetList)  // 页面列表
	// 	page.POST("", pageCtrl.SaveOrUpdate) // 新增/编辑页面
	// 	page.DELETE("", pageCtrl.Delete)     // 删除页面
	// }
}

// 博客前台的接口: 大部分不需要登录, 部分需要登录
func registerBlogHandler(r *gin.Engine) {
	base := r.Group("/api/front")

	base.GET("/about", blogInfoCtrl.GetAbout) // 获取关于我
	// base.GET("/page", pageCtrl.GetList)       // 前台页面

	// 使用新的 FrontController 注册前台路由
	front.RegisterFrontRouter(base, frontCtrl)

	// 需要登录才能进行的操作
	base.Use(middleware.JWTAuth())
	{
		// base.POST("/upload", uploadCtrl.UploadFile)    // 文件上传
		base.GET("/user/info", userCtrl.GetInfo)       // 根据 Token 获取用户信息
		base.PUT("/user/info", userCtrl.UpdateCurrent) // 根据 Token 更新当前用户信息

		// 使用新的 FrontController 注册需要登录的前台路由
		front.RegisterFrontAuthRouter(base, frontCtrl)
	}
}
