package service

import (
	"context"
	global "gin-blog/internal/global"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/model/dto/response"
	"gin-blog/internal/model/entity"
	"gin-blog/internal/repository"
	"gin-blog/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

type BlogInfoService interface {
	// BlogInfo
	GetHomeInfo(c *gin.Context) (response.BlogHomeVO, error)
	GetAbout(c *gin.Context) (string, error)
	UpdateAbout(c *gin.Context, req request.AboutReq) error
	Report(c *gin.Context) error

	// Config
	GetConfigMap(c *gin.Context) (map[string]string, error)
	UpdateConfigMap(c *gin.Context, m map[string]string) error

	// Page
	GetPageList(c *gin.Context) ([]entity.Page, int64, error)
	SaveOrUpdatePage(c *gin.Context, req request.AddOrEditPageReq) (*entity.Page, error)
	DeletePages(c *gin.Context, ids []int) error
}

type blogInfoService struct {
	repo repository.BlogInfoRepository
}

func NewBlogInfoService(repo repository.BlogInfoRepository) BlogInfoService {
	return &blogInfoService{repo: repo}
}

// BlogInfo implementations
func (s *blogInfoService) GetHomeInfo(c *gin.Context) (response.BlogHomeVO, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()

	var articleCount int64
	if err := db.Model(&entity.Article{}).Where("status = ? AND is_delete = ?", entity.ARTICLE_STATUS_PUBLIC, false).Count(&articleCount).Error; err != nil {
		return response.BlogHomeVO{}, err
	}

	var userCount int64
	if err := db.Model(&entity.UserInfo{}).Count(&userCount).Error; err != nil {
		return response.BlogHomeVO{}, err
	}

	var messageCount int64
	if err := db.Table("message").Count(&messageCount).Error; err != nil {
		return response.BlogHomeVO{}, err
	}

	viewCount, err := rdb.Get(rctx, global.VIEW_COUNT).Int()
	if err != nil && err != redis.Nil {
		return response.BlogHomeVO{}, err
	}

	return response.BlogHomeVO{
		ArticleCount: int(articleCount),
		UserCount:    int(userCount),
		MessageCount: int(messageCount),
		ViewCount:    viewCount,
	}, nil
}

func (s *blogInfoService) GetAbout(c *gin.Context) (string, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.GetConfig(db, global.CONFIG_ABOUT)
}

func (s *blogInfoService) UpdateAbout(c *gin.Context, req request.AboutReq) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.UpdateConfig(db, global.CONFIG_ABOUT, req.Content)
}

func (s *blogInfoService) Report(c *gin.Context) error {
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	ipAddress := utils.IP.GetIpAddress(c)
	userAgent := utils.IP.GetUserAgent(c)
	var uuid string
	if userAgent != nil {
		uuid = utils.MD5(ipAddress + userAgent.Name + " " + userAgent.Version.String() + userAgent.OS + " " + userAgent.OSVersion.String())
	} else {
		uuid = utils.MD5(ipAddress)
	}

	ctx := context.Background()

	if !rdb.SIsMember(ctx, global.KEY_UNIQUE_VISITOR_SET, uuid).Val() {
		ipSource := utils.IP.GetIpSource(ipAddress)
		if ipSource != "" {
			address := strings.Split(ipSource, "|")
			if len(address) > 2 {
				province := strings.ReplaceAll(address[2], "省", "")
				rdb.HIncrBy(ctx, global.VISITOR_AREA, province, 1)
			} else {
				rdb.HIncrBy(ctx, global.VISITOR_AREA, "未知", 1)
			}
		} else {
			rdb.HIncrBy(ctx, global.VISITOR_AREA, "未知", 1)
		}
		rdb.Incr(ctx, global.VIEW_COUNT)
		rdb.SAdd(ctx, global.KEY_UNIQUE_VISITOR_SET, uuid)
	}

	return nil
}

// Config implementations
func (s *blogInfoService) GetConfigMap(c *gin.Context) (map[string]string, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.GetConfigMap(db)
}

func (s *blogInfoService) UpdateConfigMap(c *gin.Context, m map[string]string) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.UpdateConfigMap(db, m)
}

// Page implementations
func (s *blogInfoService) GetPageList(c *gin.Context) ([]entity.Page, int64, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.GetPageList(db)
}

func (s *blogInfoService) SaveOrUpdatePage(c *gin.Context, req request.AddOrEditPageReq) (*entity.Page, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	page := &entity.Page{
		Model: entity.Model{ID: req.ID},
		Name:  req.Name,
		Label: req.Label,
		Cover: req.Cover,
	}
	err := s.repo.SaveOrUpdatePage(db, page)
	return page, err
}

func (s *blogInfoService) DeletePages(c *gin.Context, ids []int) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.DeletePages(db, ids)
}
