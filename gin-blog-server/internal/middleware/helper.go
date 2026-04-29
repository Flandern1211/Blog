package middleware

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/model/entity"
	"gin-blog/internal/repository"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

// 获取 *gorm.DB
func GetDB(c *gin.Context) *gorm.DB {
	return c.MustGet(g.CTX_DB).(*gorm.DB)
}

// 获取 *redis.Client
func GetRDB(c *gin.Context) *redis.Client {
	return c.MustGet(g.CTX_RDB).(*redis.Client)
}

// 获取当前登录用户信息
func CurrentUserAuth(c *gin.Context) (*entity.UserAuth, error) {
	key := g.CTX_USER_AUTH
	authRepo := repository.NewAuthRepository()

	// 1. 从 gin context 获取完整用户
	if cache, exist := c.Get(key); exist && cache != nil {
		if user, ok := cache.(*entity.UserAuth); ok {
			return user, nil
		}
	}

	// 2. 从 gin context 获取 UserID, 查库(兼容 JWTAuth 仅设 int 的场景)
	if idVal, exist := c.Get(key); exist {
		if id, ok := idVal.(int); ok && id > 0 {
			db := GetDB(c)
			user, err := authRepo.GetUserAuthInfoById(db, id)
			if err != nil {
				return nil, err
			}
			c.Set(key, user)
			return user, nil
		}
	}

	// 3. 从 session 获取 UserID (兼容 /api/logout 等无 JWT 的路由)
	session := sessions.Default(c)
	val := session.Get(key)
	id, ok := val.(int)
	if !ok {
		return nil, errors.New("session 中没有有效的 user_auth_id")
	}

	db := GetDB(c)
	user, err := authRepo.GetUserAuthInfoById(db, id)
	if err != nil {
		return nil, err
	}

	c.Set(key, user)
	return user, nil
}

// 获取当前登录用户 ID
func GetUserID(c *gin.Context) int {
	// 1. 从 gin context 获取
	if cache, exist := c.Get(g.CTX_USER_AUTH); exist && cache != nil {
		switch v := cache.(type) {
		case *entity.UserAuth:
			return v.ID
		case int:
			return v
		}
	}

	// 2. 从 session 获取
	session := sessions.Default(c)
	if id, ok := session.Get(g.CTX_USER_AUTH).(int); ok {
		return id
	}

	return 0
}

// 判断当前登录用户是否为超级管理员
func IsSuper(c *gin.Context) bool {
	// 1. 从 gin context 中获取
	if val, exist := c.Get(g.CTX_IS_SUPER); exist && val != nil {
		if isSuper, ok := val.(bool); ok {
			return isSuper
		}
	}

	// 2. 从 session 中获取
	session := sessions.Default(c)
	if isSuper, ok := session.Get(g.CTX_IS_SUPER).(bool); ok {
		return isSuper
	}

	return false
}
