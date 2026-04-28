package response

import (
	"gin-blog/internal/model/entity"
	"time"
)

type UserInfoVO struct {
	entity.UserInfo
	ArticleLikeSet []string `json:"article_like_set"`
	CommentLikeSet []string `json:"comment_like_set"`
}

type UserVO struct {
	ID            int              `json:"id"`
	UserInfoId    int              `json:"user_info_id"`
	Info          *entity.UserInfo `json:"info"`
	Roles         []*entity.Role   `json:"roles"`
	LoginType     int              `json:"login_type"`
	IpAddress     string           `json:"ip_address"`
	IpSource      string           `json:"ip_source"`
	CreatedAt     time.Time        `json:"created_at"`
	LastLoginTime *time.Time       `json:"last_login_time"`
	IsDisable     bool             `json:"is_disable"`
}
