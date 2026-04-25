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

	// 1. 从 gin context 中获取
	if cache, exist := c.Get(key); exist && cache != nil {
		return cache.(*entity.UserAuth), nil
	}

	// 2. 从 session 中获取 id
	session := sessions.Default(c)
	id := session.Get(key)
	if id == nil {
		return nil, errors.New("session 中没有 user_auth_id")
	}

	// 3. 根据 id 从数据库获取
	db := GetDB(c)
	user, err := authRepo.GetUserAuthInfoById(db, id.(int))
	if err != nil {
		return nil, err
	}

	c.Set(key, user)
	return user, nil
}

// 获取当前登录用户 ID
func GetUserID(c *gin.Context) int {
	// 1. 从 gin context 中获取 UserAuth 对象
	if cache, exist := c.Get(g.CTX_USER_AUTH); exist && cache != nil {
		if user, ok := cache.(*entity.UserAuth); ok {
			return user.ID
		}
	}

	// 2. 从 session 中获取 id
	session := sessions.Default(c)
	if id := session.Get(g.CTX_USER_AUTH); id != nil {
		return id.(int)
	}

	return 0
}

// 判断当前登录用户是否为超级管理员
func IsSuper(c *gin.Context) bool {
	// 1. 从 gin context 中获取
	if isSuper, exist := c.Get(g.CTX_IS_SUPER); exist && isSuper != nil {
		return isSuper.(bool)
	}

	// 2. 从 session 中获取
	session := sessions.Default(c)
	if isSuper := session.Get(g.CTX_IS_SUPER); isSuper != nil {
		return isSuper.(bool)
	}

	return false
}
