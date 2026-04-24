package response

import (
	"gin-blog/internal/model/entity"
	"time"
)

type ArticleVO struct {
	entity.Article
	LikeCount    int `json:"like_count"`
	ViewCount    int `json:"view_count"`
	CommentCount int `json:"comment_count"`
}

type CategoryVO struct {
	entity.Category
	ArticleCount int `json:"article_count"`
}

type TagVO struct {
	entity.Tag
	ArticleCount int `json:"article_count"`
}

type BlogArticleVO struct {
	entity.Article
	CommentCount int64 `json:"comment_count"`
	LikeCount    int64 `json:"like_count"`
	ViewCount    int64 `json:"view_count"`

	LastArticle       ArticlePaginationVO  `json:"last_article"`
	NextArticle       ArticlePaginationVO  `json:"next_article"`
	RecommendArticles []RecommendArticleVO `json:"recommend_articles"`
	NewestArticles    []RecommendArticleVO `json:"newest_articles"`
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
