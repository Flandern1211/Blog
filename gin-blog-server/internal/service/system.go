package service

import (
	global "gin-blog/internal/global"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/model/entity"
	"gin-blog/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SystemService interface {
	// FriendLink
	GetLinkList(c *gin.Context, query request.FriendLinkQuery) ([]entity.FriendLink, int64, error)
	SaveOrUpdateLink(c *gin.Context, req request.AddOrEditLinkReq) (*entity.FriendLink, error)
	DeleteLinks(c *gin.Context, ids []int) error

	// OperationLog
	GetOperationLogList(c *gin.Context, query request.OperationLogQuery) ([]entity.OperationLog, int64, error)
	DeleteOperationLogs(c *gin.Context, ids []int) error
}

type systemService struct {
	repo repository.SystemRepository
}

func NewSystemService(repo repository.SystemRepository) SystemService {
	return &systemService{repo: repo}
}

// FriendLink implementations
func (s *systemService) GetLinkList(c *gin.Context, query request.FriendLinkQuery) ([]entity.FriendLink, int64, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.GetLinkList(db, query.Page, query.Size, query.Keyword)
}

func (s *systemService) SaveOrUpdateLink(c *gin.Context, req request.AddOrEditLinkReq) (*entity.FriendLink, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	link := &entity.FriendLink{
		Model:   entity.Model{ID: req.ID},
		Name:    req.Name,
		Avatar:  req.Avatar,
		Address: req.Address,
		Intro:   req.Intro,
	}
	err := s.repo.SaveOrUpdateLink(db, link)
	return link, err
}

func (s *systemService) DeleteLinks(c *gin.Context, ids []int) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.DeleteLinks(db, ids)
}

// OperationLog implementations
func (s *systemService) GetOperationLogList(c *gin.Context, query request.OperationLogQuery) ([]entity.OperationLog, int64, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.GetOperationLogList(db, query.Page, query.Size, query.Keyword)
}

func (s *systemService) DeleteOperationLogs(c *gin.Context, ids []int) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.DeleteOperationLogs(db, ids)
}
