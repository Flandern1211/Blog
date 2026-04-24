package entity

const (
	TYPE_ARTICLE = iota + 1 // 文章
	TYPE_LINK               // 友链
	TYPE_TALK               // 说说
)

type Comment struct {
	Model
	UserId      int    `json:"user_id"`       // 评论者
	ReplyUserId int    `json:"reply_user_id"` // 被回复者
	TopicId     int    `json:"topic_id"`      // 评论的文章
	ParentId    int    `json:"parent_id"`     // 父评论 被回复的评论
	Content     string `gorm:"type:varchar(500);not null" json:"content"`
	Type        int    `gorm:"type:tinyint(1);not null;comment:评论类型(1.文章 2.友链 3.说说)" json:"type"`
	IsReview    bool   `json:"is_review"`

	User      *UserAuth `gorm:"foreignKey:UserId" json:"user"`
	ReplyUser *UserAuth `gorm:"foreignKey:ReplyUserId" json:"reply_user"`
	Article   *Article  `gorm:"foreignKey:TopicId" json:"article"`
}

type Message struct {
	Model
	Nickname  string `gorm:"type:varchar(50);comment:昵称" json:"nickname"`
	Avatar    string `gorm:"type:varchar(255);comment:头像地址" json:"avatar"`
	Content   string `gorm:"type:varchar(255);comment:留言内容" json:"content"`
	IpAddress string `gorm:"type:varchar(50);comment:IP 地址" json:"ip_address"`
	IpSource  string `gorm:"type:varchar(255);comment:IP 来源" json:"ip_source"`
	Speed     int    `gorm:"type:tinyint(1);comment:弹幕速度" json:"speed"`
	IsReview  bool   `json:"is_review"`
}
