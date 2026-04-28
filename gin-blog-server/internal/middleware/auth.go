package middleware

import (
	"errors"
	"fmt"
	g "gin-blog/internal/global"
	"gin-blog/internal/repository"
	pkgErrors "gin-blog/pkg/errors"
	"gin-blog/pkg/jwt"
	"gin-blog/pkg/response"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 基于 JWT 的授权
// 从 Authorization 中获取 token, 并解析 token 获取用户信息
// 解析成功后, 将用户信息设置到 session 和 gin context 中
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := c.MustGet(g.CTX_DB).(*gorm.DB)
		authRepo := repository.NewAuthRepository()

		// 修正 URL 提取逻辑 (Issue #4)
		fullPath := c.FullPath()
		url := fullPath
		if strings.HasPrefix(fullPath, "/api/front") {
			url = fullPath[10:]
		} else if strings.HasPrefix(fullPath, "/api") {
			url = fullPath[4:]
		}
		method := c.Request.Method

		slog.Info(fmt.Sprintf("[middleware-JWTAuth] checking: %s %s (raw: %s)", method, url, fullPath))
		resource, err := authRepo.GetResource(db, url, method)

		authRequired := true
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				slog.Info("[middleware-JWTAuth] resource not exist, skip mandatory auth")
				authRequired = false
			} else {
				response.Error(c, pkgErrors.CodeDbOpError, pkgErrors.GetMessage(pkgErrors.CodeDbOpError))
				c.Abort()
				return
			}
		} else if resource.Anonymous {
			slog.Debug(fmt.Sprintf("[middleware-JWTAuth] resource: %s %s is anonymous, skip mandatory auth!", url, method))
			authRequired = false
		}

		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			if authRequired {
				response.Error(c, pkgErrors.CodeTokenNotExist, pkgErrors.GetMessage(pkgErrors.CodeTokenNotExist))
				c.Abort()
				return
			}
			c.Set("skip_check", true)
			c.Next()
			c.Set("skip_check", false)
			return
		}

		// token 的正确格式: `Bearer [tokenString]`
		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			if authRequired {
				response.Error(c, pkgErrors.CodeTokenTypeErr, pkgErrors.GetMessage(pkgErrors.CodeTokenTypeErr))
				c.Abort()
				return
			}
			c.Set("skip_check", true)
			c.Next()
			c.Set("skip_check", false)
			return
		}

		claims, err := jwt.ParseToken(g.Conf.JWT.Secret, parts[1])
		if err != nil {
			if authRequired {
				response.Error(c, pkgErrors.CodeInvalidToken, pkgErrors.GetMessage(pkgErrors.CodeInvalidToken))
				c.Abort()
				return
			}
			c.Set("skip_check", true)
			c.Next()
			c.Set("skip_check", false)
			return
		}

		// 判断 token 已过期
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			if authRequired {
				response.Error(c, pkgErrors.CodeTokenExpired, pkgErrors.GetMessage(pkgErrors.CodeTokenExpired))
				c.Abort()
				return
			}
			c.Set("skip_check", true)
			c.Next()
			c.Set("skip_check", false)
			return
		}

		user, err := authRepo.GetUserAuthInfoById(db, claims.UserID)
		if err != nil {
			if authRequired {
				response.Error(c, pkgErrors.CodeUserNotFound, pkgErrors.GetMessage(pkgErrors.CodeUserNotFound))
				c.Abort()
				return
			}
			c.Set("skip_check", true)
			c.Next()
			c.Set("skip_check", false)
			return
		}

		// session
		session := sessions.Default(c)
		session.Set(g.CTX_USER_AUTH, claims.UserID)
		session.Set(g.CTX_IS_SUPER, user.IsSuper)
		session.Save()

		// gin context
		c.Set(g.CTX_USER_AUTH, user)
		c.Set(g.CTX_IS_SUPER, user.IsSuper)

		if !authRequired {
			c.Set("skip_check", true)
		}
		c.Next()
		if !authRequired {
			c.Set("skip_check", false)
		}
	}
}

// 资源访问权限验证
func PermissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool("skip_check") {
			c.Next()
			return
		}

		db := c.MustGet(g.CTX_DB).(*gorm.DB)
		auth, err := CurrentUserAuth(c)
		if err != nil {
			response.Error(c, pkgErrors.CodeUserNotFound, pkgErrors.GetMessage(pkgErrors.CodeUserNotFound))
			c.Abort()
			return
		}

		if auth.IsSuper {
			slog.Debug("[middleware-PermissionCheck]: super admin no need to check, pass!")
			c.Next()
			return
		}

		// 修正 URL 提取逻辑 (Issue #4, #12)
		fullPath := c.FullPath()
		url := fullPath
		if strings.HasPrefix(fullPath, "/api/front") {
			url = fullPath[10:]
		} else if strings.HasPrefix(fullPath, "/api") {
			url = fullPath[4:]
		}
		method := c.Request.Method
		authRepo := repository.NewAuthRepository()

		slog.Debug(fmt.Sprintf("[middleware-PermissionCheck] %v, %v, %v\n", auth.Username, url, method))
		pass := false
		for _, role := range auth.Roles {
			slog.Debug(fmt.Sprintf("[middleware-PermissionCheck] check role: %v\n", role.Name))
			p, err := authRepo.CheckRoleAuth(db, role.ID, url, method)
			if err != nil {
				response.Error(c, pkgErrors.CodeDbOpError, pkgErrors.GetMessage(pkgErrors.CodeDbOpError))
				c.Abort()
				return
			}
			if p {
				pass = true
				break
			}
		}

		if !pass {
			response.Error(c, pkgErrors.CodePermissionErr, pkgErrors.GetMessage(pkgErrors.CodePermissionErr))
			c.Abort()
			return
		}

		slog.Debug("[middleware-PermissionCheck]: pass")
		c.Next()
	}
}
