package response

import "gin-blog/internal/model/entity"

type LoginVO struct {
	entity.UserInfo

	ArticleLikeSet []string `json:"article_like_set"`
	CommentLikeSet []string `json:"comment_like_set"`
	Token          string   `json:"token"`
	IsSuper        bool     `json:"is_super"`
}
