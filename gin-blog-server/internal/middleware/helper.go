package middleware

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

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
func CurrentUserAuth(c *gin.Context) (*model.UserAuth, error) {
	key := g.CTX_USER_AUTH

	// 1. 从 gin context 中获取
	if cache, exist := c.Get(key); exist && cache != nil {
		return cache.(*model.UserAuth), nil
	}

	// 2. 从 session 中获取 id
	session := sessions.Default(c)
	id := session.Get(key)
	if id == nil {
		return nil, errors.New("session 中没有 user_auth_id")
	}

	// 3. 根据 id 从数据库获取
	db := GetDB(c)
	user, err := model.GetUserAuthInfoById(db, id.(int))
	if err != nil {
		return nil, err
	}

	c.Set(key, user)
	return user, nil
}
