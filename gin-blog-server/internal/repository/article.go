package repository

import (
	"gin-blog/internal/model/entity"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	// Article
	GetList(db *gorm.DB, page, size int, title string, categoryId, tagId, artType, status int, isDelete *bool) ([]entity.Article, int64, error)
	GetById(db *gorm.DB, id int) (*entity.Article, error)
	SaveOrUpdate(db *gorm.DB, article *entity.Article, categoryName string, tagNames []string) error
	UpdateTop(db *gorm.DB, id int, isTop bool) error
	SoftDelete(db *gorm.DB, ids []int, isDelete bool) error
	Delete(db *gorm.DB, ids []int) error

	// Front-end specific Article methods
	GetBlogArticle(db *gorm.DB, id int) (*entity.Article, error)
	GetBlogArticleList(db *gorm.DB, page, size, categoryId, tagId int) ([]entity.Article, int64, error)
	GetRecommendList(db *gorm.DB, id, n int) ([]entity.RecommendArticleVO, error)
	GetLastArticle(db *gorm.DB, id int) (entity.ArticlePaginationVO, error)
	GetNextArticle(db *gorm.DB, id int) (entity.ArticlePaginationVO, error)
	GetNewestList(db *gorm.DB, n int) ([]entity.RecommendArticleVO, error)

	// Category
	GetCategoryList(db *gorm.DB, page, size int, keyword string) ([]entity.Category, int64, error)
	SaveOrUpdateCategory(db *gorm.DB, id int, name string) error
	DeleteCategories(db *gorm.DB, ids []int) error
	GetCategoryOption(db *gorm.DB) ([]entity.Category, error)

	// Tag
	GetTagList(db *gorm.DB, page, size int, keyword string) ([]entity.Tag, int64, error)
	SaveOrUpdateTag(db *gorm.DB, id int, name string) error
	DeleteTags(db *gorm.DB, ids []int) error
	GetTagOption(db *gorm.DB) ([]entity.Tag, error)
}

type articleRepository struct{}

func NewArticleRepository() ArticleRepository {
	return &articleRepository{}
}

// Article implementations
func (r *articleRepository) GetList(db *gorm.DB, page, size int, title string, categoryId, tagId, artType, status int, isDelete *bool) ([]entity.Article, int64, error) {
	var list []entity.Article
	var total int64
	query := db.Model(&entity.Article{}).Preload("Category").Preload("Tags")

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if categoryId != 0 {
		query = query.Where("category_id = ?", categoryId)
	}
	if tagId != 0 {
		query = query.Joins("JOIN article_tag ON article.id = article_tag.article_id").Where("article_tag.tag_id = ?", tagId)
	}
	if artType != 0 {
		query = query.Where("type = ?", artType)
	}
	if status != 0 {
		query = query.Where("status = ?", status)
	}
	if isDelete != nil {
		query = query.Where("is_delete = ?", *isDelete)
	}

	err := query.Count(&total).Offset((page - 1) * size).Limit(size).Order("is_top DESC, id DESC").Find(&list).Error
	return list, total, err
}

func (r *articleRepository) GetById(db *gorm.DB, id int) (*entity.Article, error) {
	var article entity.Article
	err := db.Preload("Category").Preload("Tags").First(&article, id).Error
	return &article, err
}

func (r *articleRepository) SaveOrUpdate(db *gorm.DB, article *entity.Article, categoryName string, tagNames []string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Handle Category
		if categoryName != "" {
			var category entity.Category
			if err := tx.Where("name = ?", categoryName).First(&category).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					category.Name = categoryName
					if err := tx.Create(&category).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}
			article.CategoryId = category.ID
		}

		// Save or Update Article
		if article.ID == 0 {
			if err := tx.Create(article).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(article).Updates(article).Error; err != nil {
				return err
			}
			// Clear tags for update
			if err := tx.Delete(&entity.ArticleTag{}, "article_id = ?", article.ID).Error; err != nil {
				return err
			}
		}

		// Handle Tags
		if len(tagNames) > 0 {
			var tags []entity.Tag
			for _, name := range tagNames {
				var tag entity.Tag
				if err := tx.Where("name = ?", name).First(&tag).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						tag.Name = name
						if err := tx.Create(&tag).Error; err != nil {
							return err
						}
					} else {
						return err
					}
				}
				tags = append(tags, tag)
			}
			var articleTags []entity.ArticleTag
			for _, tag := range tags {
				articleTags = append(articleTags, entity.ArticleTag{ArticleId: article.ID, TagId: tag.ID})
			}
			if err := tx.Create(&articleTags).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *articleRepository) UpdateTop(db *gorm.DB, id int, isTop bool) error {
	return db.Model(&entity.Article{Model: entity.Model{ID: id}}).Update("is_top", isTop).Error
}

func (r *articleRepository) SoftDelete(db *gorm.DB, ids []int, isDelete bool) error {
	return db.Model(&entity.Article{}).Where("id IN ?", ids).Update("is_delete", isDelete).Error
}

func (r *articleRepository) Delete(db *gorm.DB, ids []int) error {
	return db.Delete(&entity.Article{}, ids).Error
}

// Front-end implementations
func (r *articleRepository) GetBlogArticle(db *gorm.DB, id int) (*entity.Article, error) {
	var data entity.Article
	result := db.Preload("Category").Preload("Tags").
		Where(entity.Article{Model: entity.Model{ID: id}}).
		Where("is_delete = 0 AND status = 1").
		First(&data)
	return &data, result.Error
}

func (r *articleRepository) GetBlogArticleList(db *gorm.DB, page, size, categoryId, tagId int) ([]entity.Article, int64, error) {
	var data []entity.Article
	var total int64
	query := db.Model(&entity.Article{}).Where("is_delete = 0 AND status = 1")

	if categoryId != 0 {
		query = query.Where("category_id = ?", categoryId)
	}
	if tagId != 0 {
		query = query.Where("id IN (SELECT article_id FROM article_tag WHERE tag_id = ?)", tagId)
	}

	query.Count(&total)
	result := query.Preload("Tags").Preload("Category").
		Order("is_top DESC, id DESC").
		Offset((page - 1) * size).Limit(size).
		Find(&data)

	return data, total, result.Error
}

func (r *articleRepository) GetRecommendList(db *gorm.DB, id, n int) ([]entity.RecommendArticleVO, error) {
	var list []entity.RecommendArticleVO
	sub1 := db.Table("article_tag").Select("tag_id").Where("article_id = ?", id)
	sub2 := db.Table("(?) t1", sub1).
		Select("DISTINCT article_id").
		Joins("JOIN article_tag t ON t.tag_id = t1.tag_id").
		Where("article_id != ?", id)
	result := db.Table("(?) t2", sub2).
		Select("id, title, img, created_at").
		Joins("JOIN article a ON t2.article_id = a.id").
		Where("a.is_delete = 0 AND a.status = 1").
		Order("is_top DESC, id DESC").
		Limit(n).
		Find(&list)
	return list, result.Error
}

func (r *articleRepository) GetLastArticle(db *gorm.DB, id int) (entity.ArticlePaginationVO, error) {
	var val entity.ArticlePaginationVO
	sub := db.Table("article").Select("max(id)").Where("id < ?", id)
	result := db.Table("article").
		Select("id, title, img").
		Where("is_delete = 0 AND status = 1 AND id = (?)", sub).
		Limit(1).
		Find(&val)
	return val, result.Error
}

func (r *articleRepository) GetNextArticle(db *gorm.DB, id int) (entity.ArticlePaginationVO, error) {
	var data entity.ArticlePaginationVO
	result := db.Model(&entity.Article{}).
		Select("id, title, img").
		Where("is_delete = 0 AND status = 1 AND id > ?", id).
		Limit(1).
		Find(&data)
	return data, result.Error
}

func (r *articleRepository) GetNewestList(db *gorm.DB, n int) ([]entity.RecommendArticleVO, error) {
	var data []entity.RecommendArticleVO
	result := db.Model(&entity.Article{}).
		Select("id, title, img, created_at").
		Where("is_delete = 0 AND status = 1").
		Order("created_at DESC, id ASC").
		Limit(n).
		Find(&data)
	return data, result.Error
}

// Category implementations
func (r *articleRepository) GetCategoryList(db *gorm.DB, page, size int, keyword string) ([]entity.Category, int64, error) {
	var list []entity.Category
	var total int64
	db = db.Model(&entity.Category{})
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	err := db.Count(&total).Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *articleRepository) SaveOrUpdateCategory(db *gorm.DB, id int, name string) error {
	if id == 0 {
		return db.Create(&entity.Category{Name: name}).Error
	}
	return db.Model(&entity.Category{Model: entity.Model{ID: id}}).Update("name", name).Error
}

func (r *articleRepository) DeleteCategories(db *gorm.DB, ids []int) error {
	return db.Delete(&entity.Category{}, ids).Error
}

func (r *articleRepository) GetCategoryOption(db *gorm.DB) ([]entity.Category, error) {
	var list []entity.Category
	err := db.Model(&entity.Category{}).Select("id", "name").Find(&list).Error
	return list, err
}

// Tag implementations
func (r *articleRepository) GetTagList(db *gorm.DB, page, size int, keyword string) ([]entity.Tag, int64, error) {
	var list []entity.Tag
	var total int64
	db = db.Model(&entity.Tag{})
	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}
	err := db.Count(&total).Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *articleRepository) SaveOrUpdateTag(db *gorm.DB, id int, name string) error {
	if id == 0 {
		return db.Create(&entity.Tag{Name: name}).Error
	}
	return db.Model(&entity.Tag{Model: entity.Model{ID: id}}).Update("name", name).Error
}

func (r *articleRepository) DeleteTags(db *gorm.DB, ids []int) error {
	return db.Delete(&entity.Tag{}, ids).Error
}

func (r *articleRepository) GetTagOption(db *gorm.DB) ([]entity.Tag, error) {
	var list []entity.Tag
	err := db.Model(&entity.Tag{}).Select("id", "name").Find(&list).Error
	return list, err
}
