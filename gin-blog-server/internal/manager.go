package ginblog

import (
	"gin-blog/docs"
	"gin-blog/internal/api/v1/article"
	"gin-blog/internal/api/v1/auth"
	"gin-blog/internal/api/v1/blog_info"
	"gin-blog/internal/api/v1/category"
	"gin-blog/internal/api/v1/comment"
	"gin-blog/internal/api/v1/config"
	"gin-blog/internal/api/v1/front"
	"gin-blog/internal/api/v1/message"
	"gin-blog/internal/api/v1/operation_log"
	"gin-blog/internal/api/v1/page"
	"gin-blog/internal/api/v1/permissionManage/menu"
	"gin-blog/internal/api/v1/permissionManage/resource"
	"gin-blog/internal/api/v1/permissionManage/role"
	"gin-blog/internal/api/v1/system"
	"gin-blog/internal/api/v1/tag"
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
	categoryCtrl *category.CategoryController
	tagCtrl      *tag.TagController
	articleCtrl  *article.ArticleController

	uploadCtrl *upload.UploadController
	userCtrl   *user.UserController
	authCtrl   *auth.AuthController
	configCtrl *config.ConfigController

	commentCtrl *comment.CommentController
	messageCtrl *message.MessageController

	roleCtrl     *role.RoleController
	resourceCtrl *resource.ResourceController
	menuCtrl     *menu.MenuController

	blogInfoCtrl *blog_info.BlogInfoController
	pageCtrl     *page.PageController
	linkCtrl     *system.LinkController
	logCtrl      *operation_log.OperationLogController

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
	categoryCtrl = category.NewCategoryController(articleSvc)
	tagCtrl = tag.NewTagController(articleSvc)

	authCtrl = auth.NewAuthController(authSvc)
	configCtrl = config.NewConfigController(blogInfoSvc)
	uploadCtrl = upload.NewUploadController(service.NewUploadService())
	userCtrl = user.NewUserController(userSvc)

	commentCtrl = comment.NewCommentController(interactSvc)
	messageCtrl = message.NewMessageController(interactSvc)

	roleCtrl = role.NewRoleController(permissionSvc)
	resourceCtrl = resource.NewResourceController(permissionSvc)
	menuCtrl = menu.NewMenuController(permissionSvc)

	blogInfoCtrl = blog_info.NewBlogInfoController(blogInfoSvc)
	pageCtrl = page.NewPageController(blogInfoSvc)
	linkCtrl = system.NewLinkController(systemSvc)
	logCtrl = operation_log.NewOperationLogController(systemSvc)

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

	auth.RegisterAuthRouter(base, authCtrl)
	// config.RegisterConfigRouter(base, configCtrl) // moved to admin handler for auth
}

// 后台管理系统的接口: 全部需要 登录 + 鉴权
func registerAdminHandler(r *gin.Engine) {
	admin := r.Group("/api")

	// 管理员登录路由（无需鉴权，登录时校验后台权限）
	auth.RegisterAdminAuthRouter(admin, authCtrl)

	// !注意使用中间件的顺序
	admin.Use(middleware.JWTAuth())
	admin.Use(middleware.PermissionCheck())
	admin.Use(middleware.OperationLog())
	admin.Use(middleware.ListenOnline())

	blog_info.RegisterBlogInfoRouter(admin, blogInfoCtrl)
	blog_info.RegisterSettingRouter(admin, blogInfoCtrl)
	upload.RegisterUploadRouter(admin, uploadCtrl)
	user.RegisterUserRouter(admin, userCtrl)
	category.RegisterCategoryRouter(admin, categoryCtrl)
	tag.RegisterTagRouter(admin, tagCtrl)
	article.RegisterArticleRouter(admin, articleCtrl)
	comment.RegisterCommentRouter(admin, commentCtrl)
	message.RegisterMessageRouter(admin, messageCtrl)
	resource.RegisterResourceRouter(admin, resourceCtrl)
	menu.RegisterMenuRouter(admin, menuCtrl)
	role.RegisterRoleRouter(admin, roleCtrl)
	operation_log.RegisterOperationLogRouter(admin, logCtrl)
	page.RegisterPageRouter(admin, pageCtrl)
	system.RegisterLinkRouter(admin, linkCtrl)
	config.RegisterConfigRouter(admin, configCtrl)
}

// 博客前台的接口: 大部分不需要登录, 部分需要登录
func registerBlogHandler(r *gin.Engine) {
	base := r.Group("/api/front")

	base.GET("/about", blogInfoCtrl.GetAbout) // 获取关于我
	base.GET("/page", pageCtrl.GetList)       // 前台页面

	// 使用新的 FrontController 注册前台路由
	front.RegisterFrontRouter(base, frontCtrl)

	// 需要登录才能进行的操作
	base.Use(middleware.JWTAuth())
	{
		// base.POST("/upload", uploadCtrl.UploadFile)    // 文件上传
		base.GET("/user/info", userCtrl.GetInfo)                  // 根据 Token 获取用户信息
		base.PUT("/user/info", userCtrl.UpdateCurrent)            // 根据 Token 更新当前用户信息
		base.PUT("/user/password", userCtrl.UpdatePasswordByCode) // 前台用户通过验证码修改密码

		// 使用新的 FrontController 注册需要登录的前台路由
		front.RegisterFrontAuthRouter(base, frontCtrl)
	}
}
