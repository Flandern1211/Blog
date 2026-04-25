package middleware

import (
	"context"
	"fmt"
	g "gin-blog/internal/global"
	pkgErrors "gin-blog/pkg/errors"
	"gin-blog/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

// 监听在线状态中间件
// 登录时: 移除用户的强制下线标记
// 退出登录时: 添加用户的在线标记
func ListenOnline() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		rdb := c.MustGet(g.CTX_RDB).(*redis.Client)

		auth, err := CurrentUserAuth(c)
		if err != nil {
			response.BizError(c, pkgErrors.NewWithErr(pkgErrors.CodeUserAuthError, pkgErrors.GetMessage(pkgErrors.CodeUserAuthError), err))
			return
		}

		onlineKey := g.ONLINE_USER + strconv.Itoa(auth.ID)
		offlineKey := g.OFFLINE_USER + strconv.Itoa(auth.ID)

		// 判断当前用户是否被强制下线
		if rdb.Exists(ctx, offlineKey).Val() == 1 {
			fmt.Println("用户被强制下线")
			response.Error(c, pkgErrors.CodeForceOffline, pkgErrors.GetMessage(pkgErrors.CodeForceOffline))
			c.Abort()
			return
		}

		// 每次发送请求会更新 Redis 中的在线状态: 重新计算 10 分钟
		rdb.Set(ctx, onlineKey, auth, 10*time.Minute)
		c.Next()
	}
}
