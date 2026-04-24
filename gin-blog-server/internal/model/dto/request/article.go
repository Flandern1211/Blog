package request

type AddOrEditArticleReq struct {
	ID           int      `json:"id"`
	Title        string   `json:"title" binding:"required" validate:"required"`
	Desc         string   `json:"desc"`
	Content      string   `json:"content" binding:"required" validate:"required"`
	Img          string   `json:"img"`
	Type         int      `json:"type" binding:"required,min=1,max=3" validate:"required,min=1,max=3"`
	Status       int      `json:"status" binding:"required,min=1,max=3" validate:"required,min=1,max=3"`
	IsTop        bool     `json:"is_top"`
	OriginalUrl  string   `json:"original_url"`
	TagNames     []string `json:"tag_names"`
	CategoryName string   `json:"category_name"`
}

type ArticleQuery struct {
	PageQuery
	Title      string `form:"title"`
	CategoryId int    `form:"category_id"`
	TagId      int    `form:"tag_id"`
	Type       int    `form:"type"`
	Status     int    `form:"status"`
	IsDelete   *bool  `form:"is_delete"`
}

type UpdateArticleTopReq struct {
	ID    int  `json:"id"`
	IsTop bool `json:"is_top"`
}

type SoftDeleteReq struct {
	Ids      []int `json:"ids"`
	IsDelete bool  `json:"is_delete"`
}

type CategoryQuery struct {
	PageQuery
	Keyword string `form:"keyword"`
}

type AddOrEditCategoryReq struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required" validate:"required"`
}

type TagQuery struct {
	PageQuery
	Keyword string `form:"keyword"`
}

type AddOrEditTagReq struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required" validate:"required"`
}
