package auth

import (
	"cloudque/internal/model/dto/request"
	"cloudque/internal/service"
	"cloudque/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	authService service.AuthService
	userService service.UserService
}

// NewAuthController 创建认证控制器
func NewAuthController(authService service.AuthService, userService service.UserService) *AuthController {
	return &AuthController{
		authService: authService,
		userService: userService,
	}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := ctrl.userService.Register(&req); err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, nil)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取 Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=response.LoginResponse}
// @Router /api/v1/auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := ctrl.authService.Login(&req)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, resp)
}

// RefreshToken 刷新 Token
// @Summary 刷新 Token
// @Description 使用旧 Token 获取新 Token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body request.RefreshTokenRequest true "Token"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/refresh [post]
func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	newToken, err := ctrl.authService.RefreshToken(req.Token)
	if err != nil {
		response.BizError(c, err)
		return
	}

	response.Success(c, map[string]string{
		"token": newToken,
	})
}
