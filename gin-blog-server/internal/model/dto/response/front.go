package response

import "time"

type ArchiveVO struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type ArticleSearchVO struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type FrontHomeVO struct {
	ArticleCount  int64             `json:"article_count"`
	UserCount     int64             `json:"user_count"`
	MessageCount  int64             `json:"message_count"`
	CategoryCount int64             `json:"category_count"`
	TagCount      int64             `json:"tag_count"`
	ViewCount     int64             `json:"view_count"`
	Config        map[string]string `json:"blog_config"`
}

type FrontCommentVO struct {
	ID          int              `json:"id"`
	UserId      int              `json:"user_id"`
	ReplyUserId int              `json:"reply_user_id"`
	TopicId     int              `json:"topic_id"`
	ParentId    int              `json:"parent_id"`
	Content     string           `json:"content"`
	Type        int              `json:"type"`
	IsReview    bool             `json:"is_review"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Nickname    string           `json:"nickname"`
	Avatar      string           `json:"avatar"`
	ReplyUser   string           `json:"reply_user"`
	LikeCount   int              `json:"like_count"`
	ReplyCount  int              `json:"reply_count"`
	ReplyList   []FrontCommentVO `json:"reply_list"`
}
