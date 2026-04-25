package repository

import (
	"gin-blog/internal/model/entity"
	"strconv"

	"gorm.io/gorm"
)

type BlogInfoRepository interface {
	// Config
	GetConfigMap(db *gorm.DB) (map[string]string, error)
	UpdateConfigMap(db *gorm.DB, m map[string]string) error
	GetConfig(db *gorm.DB, key string) (string, error)
	GetConfigBool(db *gorm.DB, key string) bool
	GetConfigInt(db *gorm.DB, key string) int
	UpdateConfig(db *gorm.DB, key, value string) error

	// Page
	GetPageList(db *gorm.DB) ([]entity.Page, int64, error)
	SaveOrUpdatePage(db *gorm.DB, page *entity.Page) error
	DeletePages(db *gorm.DB, ids []int) error
}

type blogInfoRepository struct{}

func NewBlogInfoRepository() BlogInfoRepository {
	return &blogInfoRepository{}
}

// Config implementations
func (r *blogInfoRepository) GetConfigMap(db *gorm.DB) (map[string]string, error) {
	var configs []entity.Config
	if err := db.Find(&configs).Error; err != nil {
		return nil, err
	}
	m := make(map[string]string)
	for _, config := range configs {
		m[config.Key] = config.Value
	}
	return m, nil
}

func (r *blogInfoRepository) UpdateConfigMap(db *gorm.DB, m map[string]string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for k, v := range m {
			if err := tx.Model(&entity.Config{}).Where("`key` = ?", k).Update("value", v).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *blogInfoRepository) GetConfig(db *gorm.DB, key string) (string, error) {
	var config entity.Config
	if err := db.Where("`key` = ?", key).First(&config).Error; err != nil {
		return "", err
	}
	return config.Value, nil
}

func (r *blogInfoRepository) GetConfigBool(db *gorm.DB, key string) bool {
	val, err := r.GetConfig(db, key)
	if err != nil {
		return false
	}
	return val == "true"
}

func (r *blogInfoRepository) GetConfigInt(db *gorm.DB, key string) int {
	val, err := r.GetConfig(db, key)
	if err != nil {
		return 0
	}
	result, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return result
}

func (r *blogInfoRepository) UpdateConfig(db *gorm.DB, key, value string) error {
	return db.Where(&entity.Config{Key: key}).Assign(&entity.Config{Value: value}).FirstOrCreate(&entity.Config{}).Error
}

// Page implementations
func (r *blogInfoRepository) GetPageList(db *gorm.DB) ([]entity.Page, int64, error) {
	var pages []entity.Page
	var total int64
	if err := db.Model(&entity.Page{}).Count(&total).Find(&pages).Error; err != nil {
		return nil, 0, err
	}
	return pages, total, nil
}

func (r *blogInfoRepository) SaveOrUpdatePage(db *gorm.DB, page *entity.Page) error {
	if page.ID > 0 {
		return db.Updates(page).Error
	}
	return db.Create(page).Error
}

func (r *blogInfoRepository) DeletePages(db *gorm.DB, ids []int) error {
	return db.Delete(&entity.Page{}, ids).Error
}
