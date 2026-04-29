package middleware

import (
	"errors"
	"fmt"
	g "gin-blog/internal/global"
	"gin-blog/internal/repository"
	"gin-blog/internal/utils"
	pkgErrors "gin-blog/pkg/errors"
	"gin-blog/pkg/jwt"
	"gin-blog/pkg/response"
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

// JWTAuth 基于 JWT 的身份认证中间件
// 只做身份认证: JWT 解析 + Redis 白名单检查
// 通过后将 UserID 写入 context, 不查 DB, 不设 session
// - 无 token → 放行(后续中间件决定是否拦截)
// - token 无效/过期 → 拒绝
// - token 有效 → c.Set(g.CTX_USER_AUTH, claims.UserID)
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			return
		}

		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, pkgErrors.CodeTokenTypeErr, pkgErrors.GetMessage(pkgErrors.CodeTokenTypeErr))
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(g.Conf.JWT.Secret, parts[1])
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				response.Error(c, pkgErrors.CodeTokenExpired, pkgErrors.GetMessage(pkgErrors.CodeTokenExpired))
			} else {
				response.Error(c, pkgErrors.CodeInvalidToken, pkgErrors.GetMessage(pkgErrors.CodeInvalidToken))
			}
			c.Abort()
			return
		}
		//从redis中
		rdb := c.MustGet(g.CTX_RDB).(*redis.Client)
		tokenKey := g.TOKEN_WHITELIST + utils.MD5(parts[1])
		if rdb.Exists(c.Request.Context(), tokenKey).Val() == 0 {
			response.Error(c, pkgErrors.CodeInvalidToken, pkgErrors.GetMessage(pkgErrors.CodeInvalidToken))
			c.Abort()
			return
		}

		c.Set(g.CTX_USER_AUTH, claims.UserID)
		c.Next()
	}
}

// PermissionCheck 资源访问权限验证中间件
// 独立判断当前请求是否需要鉴权, 以及当前用户是否有权限访问
func PermissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, err := CurrentUserAuth(c)
		if err != nil {
			response.Error(c, pkgErrors.CodeTokenNotExist, pkgErrors.GetMessage(pkgErrors.CodeTokenNotExist))
			c.Abort()
			return
		}

		db := c.MustGet(g.CTX_DB).(*gorm.DB)
		authRepo := repository.NewAuthRepository()

		// 修正 URL 提取逻辑 (Issue #4, #12)
		fullPath := c.FullPath()
		url := fullPath
		if strings.HasPrefix(fullPath, "/api/front") {
			url = fullPath[10:]
		} else if strings.HasPrefix(fullPath, "/api") {
			url = fullPath[4:]
		}
		method := c.Request.Method

		slog.Debug(fmt.Sprintf("[middleware-PermissionCheck] checking: %s %s", method, url))

		// 查资源表, 判断该接口是否需要鉴权
		resource, err := authRepo.GetResource(db, url, method)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				slog.Debug("[middleware-PermissionCheck] resource not registered, pass")
				c.Next()
				return
			}
			response.Error(c, pkgErrors.CodeDbOpError, pkgErrors.GetMessage(pkgErrors.CodeDbOpError))
			c.Abort()
			return
		}

		if resource.Anonymous {
			slog.Debug("[middleware-PermissionCheck] anonymous resource, pass")
			c.Next()
			return
		}

		if auth.IsSuper {
			slog.Debug("[middleware-PermissionCheck] super admin, pass")
			c.Next()
			return
		}

		slog.Debug(fmt.Sprintf("[middleware-PermissionCheck] %v, %v, %v\n", auth.Username, url, method))
		for _, role := range auth.Roles {
			slog.Debug(fmt.Sprintf("[middleware-PermissionCheck] check role: %v\n", role.Name))
			p, err := authRepo.CheckRoleAuth(db, role.ID, url, method)
			if err != nil {
				response.Error(c, pkgErrors.CodeDbOpError, pkgErrors.GetMessage(pkgErrors.CodeDbOpError))
				c.Abort()
				return
			}
			if p {
				slog.Debug("[middleware-PermissionCheck] pass")
				c.Next()
				return
			}
		}

		response.Error(c, pkgErrors.CodePermissionErr, pkgErrors.GetMessage(pkgErrors.CodePermissionErr))
		c.Abort()
	}
}

// AdminOnly 验证用户是否有后台管理权限
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, err := CurrentUserAuth(c)
		if err != nil {
			response.Error(c, pkgErrors.CodeTokenNotExist, pkgErrors.GetMessage(pkgErrors.CodeTokenNotExist))
			c.Abort()
			return
		}

		if auth.IsSuper {
			c.Next()
			return
		}

		db := c.MustGet(g.CTX_DB).(*gorm.DB)
		authRepo := repository.NewAuthRepository()
		hasResource, err := authRepo.CheckUserHasResource(db, auth.ID, g.RESOURCE_BACKEND_LOGIN, g.METHOD_BACKEND_LOGIN)
		if err != nil {
			response.Error(c, pkgErrors.CodeDbOpError, pkgErrors.GetMessage(pkgErrors.CodeDbOpError))
			c.Abort()
			return
		}
		if !hasResource {
			response.Error(c, pkgErrors.CodePermissionErr, pkgErrors.GetMessage(pkgErrors.CodePermissionErr))
			c.Abort()
			return
		}

		c.Next()
	}
}
