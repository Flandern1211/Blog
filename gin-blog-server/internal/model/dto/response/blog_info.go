package response

type BlogHomeVO struct {
	ArticleCount int `json:"article_count"`
	UserCount    int `json:"user_count"`
	MessageCount int `json:"message_count"`
	ViewCount    int `json:"view_count"`
}
