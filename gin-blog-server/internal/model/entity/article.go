package entity

import "time"

const (
	ARTICLE_STATUS_PUBLIC = iota + 1 // 公开
	ARTICLE_STATUS_SECRET            // 私密
	ARTICLE_STATUS_DRAFT             // 草稿
)

const (
	ARTICLE_TYPE_ORIGINAL  = iota + 1 // 原创
	ARTICLE_TYPE_REPRINT              // 转载
	ARTICLE_TYPE_TRANSLATE            // 翻译
)

type Article struct {
	Model
	Title       string `gorm:"type:varchar(100);not null" json:"title"`
	Desc        string `json:"desc"`
	Content     string `json:"content"`
	Img         string `json:"img"`
	Type        int    `gorm:"type:tinyint;comment:类型(1-原创 2-转载 3-翻译)" json:"type"`
	Status      int    `gorm:"type:tinyint;comment:状态(1-公开 2-私密)" json:"status"`
	IsTop       bool   `json:"is_top"`
	IsDelete    bool   `json:"is_delete"`
	OriginalUrl string `json:"original_url"`

	CategoryId int `json:"category_id"`
	UserId     int `json:"-"` // user_auth_id

	Tags     []*Tag    `gorm:"many2many:article_tag;joinForeignKey:article_id" json:"tags"`
	Category *Category `gorm:"foreignkey:CategoryId" json:"category"`
	User     *UserAuth `gorm:"foreignkey:UserId" json:"user"`
}

type ArticleTag struct {
	ArticleId int
	TagId     int
}

type Category struct {
	Model
	Name string `gorm:"unique;type:varchar(20);not null" json:"name"`

	Articles []Article `gorm:"foreignkey:CategoryId" json:"articles"`
}

type Tag struct {
	Model
	Name string `gorm:"unique;type:varchar(20);not null" json:"name"`

	Articles []Article `gorm:"many2many:article_tag" json:"articles"`
}

type CategoryVO struct {
	Category
	ArticleCount int `json:"article_count"`
}

type TagVO struct {
	Tag
	ArticleCount int `json:"article_count"`
}

type ArticlePaginationVO struct {
	ID    int    `json:"id"`
	Img   string `json:"img"`
	Title string `json:"title"`
}

type RecommendArticleVO struct {
	ID        int       `json:"id"`
	Img       string    `json:"img"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}
