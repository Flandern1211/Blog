package repository

import (
	"gin-blog/internal/model/entity"

	"gorm.io/gorm"
)

type SystemRepository interface {
	// FriendLink
	GetLinkList(db *gorm.DB, page, size int, keyword string) ([]entity.FriendLink, int64, error)
	SaveOrUpdateLink(db *gorm.DB, link *entity.FriendLink) error
	DeleteLinks(db *gorm.DB, ids []int) error

	// OperationLog
	GetOperationLogList(db *gorm.DB, page, size int, keyword string) ([]entity.OperationLog, int64, error)
	DeleteOperationLogs(db *gorm.DB, ids []int) error
	CreateOperationLog(db *gorm.DB, log *entity.OperationLog) error
}

type systemRepository struct{}

func NewSystemRepository() SystemRepository {
	return &systemRepository{}
}

// FriendLink implementations
func (r *systemRepository) GetLinkList(db *gorm.DB, page, size int, keyword string) ([]entity.FriendLink, int64, error) {
	var list []entity.FriendLink
	var total int64
	query := db.Model(&entity.FriendLink{})

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%").
			Or("address LIKE ?", "%"+keyword+"%").
			Or("intro LIKE ?", "%"+keyword+"%")
	}

	err := query.Count(&total).Order("created_at DESC").Scopes(Paginate(page, size)).Find(&list).Error
	return list, total, err
}

func (r *systemRepository) SaveOrUpdateLink(db *gorm.DB, link *entity.FriendLink) error {
	if link.ID > 0 {
		return db.Updates(link).Error
	}
	return db.Create(link).Error
}

func (r *systemRepository) DeleteLinks(db *gorm.DB, ids []int) error {
	return db.Where("id IN ?", ids).Delete(&entity.FriendLink{}).Error
}

// OperationLog implementations
func (r *systemRepository) GetOperationLogList(db *gorm.DB, page, size int, keyword string) ([]entity.OperationLog, int64, error) {
	var list []entity.OperationLog
	var total int64
	query := db.Model(&entity.OperationLog{})

	if keyword != "" {
		query = query.Where("opt_module LIKE ?", "%"+keyword+"%").
			Or("opt_desc LIKE ?", "%"+keyword+"%")
	}

	err := query.Count(&total).Order("created_at DESC").Scopes(Paginate(page, size)).Find(&list).Error
	return list, total, err
}

func (r *systemRepository) DeleteOperationLogs(db *gorm.DB, ids []int) error {
	return db.Where("id IN ?", ids).Delete(&entity.OperationLog{}).Error
}

func (r *systemRepository) CreateOperationLog(db *gorm.DB, log *entity.OperationLog) error {
	return db.Create(log).Error
}
