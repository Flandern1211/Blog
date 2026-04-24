package service

import (
	global "gin-blog/internal/global"
	"gin-blog/internal/model/dto/response"
	"gin-blog/internal/model/entity"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

type FrontService interface {
	GetHomeInfo(c *gin.Context) (response.FrontHomeVO, error)
	LikeArticle(c *gin.Context, articleId int, authId int) error
	LikeComment(c *gin.Context, commentId int, authId int) error
	SearchArticle(c *gin.Context, keyword string) ([]response.ArticleSearchVO, error)
}

type frontService struct {
}

func NewFrontService() FrontService {
	return &frontService{}
}

func (s *frontService) GetHomeInfo(c *gin.Context) (response.FrontHomeVO, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()

	var data response.FrontHomeVO

	if err := db.Model(&entity.Article{}).Where("status = ? AND is_delete = ?", entity.ARTICLE_STATUS_PUBLIC, false).Count(&data.ArticleCount).Error; err != nil {
		return data, err
	}
	if err := db.Model(&entity.UserInfo{}).Count(&data.UserCount).Error; err != nil {
		return data, err
	}
	if err := db.Model(&entity.Message{}).Count(&data.MessageCount).Error; err != nil {
		return data, err
	}
	if err := db.Model(&entity.Category{}).Count(&data.CategoryCount).Error; err != nil {
		return data, err
	}
	if err := db.Model(&entity.Tag{}).Count(&data.TagCount).Error; err != nil {
		return data, err
	}

	// Get config
	var configs []entity.Config
	if err := db.Find(&configs).Error; err != nil {
		return data, err
	}
	data.Config = make(map[string]string)
	for _, conf := range configs {
		data.Config[conf.Key] = conf.Value
	}

	// Get view count
	viewCount, err := rdb.Get(rctx, global.VIEW_COUNT).Int64()
	if err != nil && err != redis.Nil {
		return data, err
	}
	data.ViewCount = viewCount

	return data, nil
}

func (s *frontService) LikeArticle(c *gin.Context, articleId int, authId int) error {
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()

	articleLikeUserKey := global.ARTICLE_USER_LIKE_SET + strconv.Itoa(authId)
	if rdb.SIsMember(rctx, articleLikeUserKey, articleId).Val() {
		rdb.SRem(rctx, articleLikeUserKey, articleId)
		rdb.HIncrBy(rctx, global.ARTICLE_LIKE_COUNT, strconv.Itoa(articleId), -1)
	} else {
		rdb.SAdd(rctx, articleLikeUserKey, articleId)
		rdb.HIncrBy(rctx, global.ARTICLE_LIKE_COUNT, strconv.Itoa(articleId), 1)
	}
	return nil
}

func (s *frontService) LikeComment(c *gin.Context, commentId int, authId int) error {
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()

	commentLikeUserKey := global.COMMENT_USER_LIKE_SET + strconv.Itoa(authId)
	if rdb.SIsMember(rctx, commentLikeUserKey, commentId).Val() {
		rdb.SRem(rctx, commentLikeUserKey, commentId)
		rdb.HIncrBy(rctx, global.COMMENT_LIKE_COUNT, strconv.Itoa(commentId), -1)
	} else {
		rdb.SAdd(rctx, commentLikeUserKey, commentId)
		rdb.HIncrBy(rctx, global.COMMENT_LIKE_COUNT, strconv.Itoa(commentId), 1)
	}
	return nil
}

func (s *frontService) SearchArticle(c *gin.Context, keyword string) ([]response.ArticleSearchVO, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	var result []response.ArticleSearchVO

	if keyword == "" {
		return result, nil
	}

	var articles []entity.Article
	err := db.Where("is_delete = ? AND status = ? AND (title LIKE ? OR content LIKE ?)", false, entity.ARTICLE_STATUS_PUBLIC, "%"+keyword+"%", "%"+keyword+"%").Find(&articles).Error
	if err != nil {
		return nil, err
	}

	for _, article := range articles {
		title := strings.ReplaceAll(article.Title, keyword, "<span style='color:#f47466'>"+keyword+"</span>")
		content := article.Content

		keywordStartIndex := unicodeIndex(content, keyword)
		if keywordStartIndex != -1 {
			preIndex, afterIndex := 0, 0
			if keywordStartIndex > 25 {
				preIndex = keywordStartIndex - 25
			}
			preText := substring(content, preIndex, keywordStartIndex)

			keywordEndIndex := keywordStartIndex + unicodeLen(keyword)
			afterLength := len([]rune(content)) - keywordEndIndex
			if afterLength > 175 {
				afterIndex = keywordEndIndex + 175
			} else {
				afterIndex = keywordEndIndex + afterLength
			}
			afterText := substring(content, keywordStartIndex, afterIndex)
			content = strings.ReplaceAll(preText+afterText, keyword, "<span style='color:#f47466'>"+keyword+"</span>")
		}

		result = append(result, response.ArticleSearchVO{
			ID:      article.ID,
			Title:   title,
			Content: content,
		})
	}

	return result, nil
}

func unicodeIndex(str, substr string) int {
	result := strings.Index(str, substr)
	if result > 0 {
		prefix := []byte(str)[0:result]
		rs := []rune(string(prefix))
		result = len(rs)
	}
	return result
}

func unicodeLen(str string) int {
	var r = []rune(str)
	return len(r)
}

func substring(source string, start int, end int) string {
	var unicodeStr = []rune(source)
	length := len(unicodeStr)
	if start >= end {
		return ""
	}
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	if start <= 0 && end >= length {
		return source
	}
	var substring = ""
	for i := start; i < end; i++ {
		substring += string(unicodeStr[i])
	}
	return substring
}
