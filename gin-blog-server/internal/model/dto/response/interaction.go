package response

import "gin-blog/internal/model/entity"

type CommentVO struct {
	entity.Comment
	LikeCount  int         `json:"like_count" gorm:"-"`
	ReplyCount int         `json:"reply_count" gorm:"-"`
	ReplyList  []CommentVO `json:"reply_list" gorm:"-"`
}
