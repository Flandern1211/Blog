package request

type FAddMessageReq struct {
	Nickname string `json:"nickname" binding:"required"`
	Avatar   string `json:"avatar"`
	Content  string `json:"content" binding:"required"`
	Speed    int    `json:"speed"`
}

type FAddCommentReq struct {
	ReplyUserId int    `json:"reply_user_id" form:"reply_user_id"`
	TopicId     int    `json:"topic_id" form:"topic_id"`
	Content     string `json:"content" form:"content"`
	ParentId    int    `json:"parent_id" form:"parent_id"`
	Type        int    `json:"type" form:"type" validate:"required,min=1,max=3" label:"评论类型"`
}

type FCommentQuery struct {
	PageQuery
	ReplyUserId int    `json:"reply_user_id" form:"reply_user_id"`
	TopicId     int    `json:"topic_id" form:"topic_id"`
	Content     string `json:"content" form:"content"`
	ParentId    int    `json:"parent_id" form:"parent_id"`
	Type        int    `json:"type" form:"type"`
}

type FArticleQuery struct {
	PageQuery
	CategoryId int `form:"category_id"`
	TagId      int `form:"tag_id"`
}
