package repository

import (
	"gin-blog/internal/model/entity"

	"gorm.io/gorm"
)

type InteractionRepository interface {
	// Message
	GetMessageList(db *gorm.DB, page, size int, nickname string, isReview *bool) ([]entity.Message, int64, error)
	DeleteMessages(db *gorm.DB, ids []int) error
	UpdateMessagesReview(db *gorm.DB, ids []int, isReview bool) error
	SaveMessage(db *gorm.DB, message *entity.Message) error

	// Comment
	GetCommentList(db *gorm.DB, page, size, typ int, isReview *bool, nickname string) ([]entity.Comment, int64, error)
	DeleteComments(db *gorm.DB, ids []int) error
	UpdateCommentsReview(db *gorm.DB, ids []int, isReview bool) error
	GetFrontCommentList(db *gorm.DB, page, size, topic, typ int) ([]entity.Comment, map[int][]entity.Comment, int64, error)
	GetCommentReplyList(db *gorm.DB, id, page, size int) ([]entity.Comment, error)
	AddComment(db *gorm.DB, comment *entity.Comment) error
	GetCommentById(db *gorm.DB, id int) (*entity.Comment, error)
}

type interactionRepository struct{}

func NewInteractionRepository() InteractionRepository {
	return &interactionRepository{}
}

// Message implementations
func (r *interactionRepository) GetMessageList(db *gorm.DB, page, size int, nickname string, isReview *bool) ([]entity.Message, int64, error) {
	var list []entity.Message
	var total int64
	query := db.Model(&entity.Message{})

	if nickname != "" {
		query = query.Where("nickname LIKE ?", "%"+nickname+"%")
	}
	if isReview != nil {
		query = query.Where("is_review = ?", *isReview)
	}

	err := query.Count(&total).Order("created_at DESC").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *interactionRepository) DeleteMessages(db *gorm.DB, ids []int) error {
	return db.Where("id IN ?", ids).Delete(&entity.Message{}).Error
}

func (r *interactionRepository) UpdateMessagesReview(db *gorm.DB, ids []int, isReview bool) error {
	return db.Model(&entity.Message{}).Where("id IN ?", ids).Update("is_review", isReview).Error
}

func (r *interactionRepository) SaveMessage(db *gorm.DB, message *entity.Message) error {
	return db.Create(message).Error
}

// Comment implementations
func (r *interactionRepository) GetCommentList(db *gorm.DB, page, size, typ int, isReview *bool, nickname string) ([]entity.Comment, int64, error) {
	var list []entity.Comment
	var total int64
	query := db.Model(&entity.Comment{})

	if nickname != "" {
		var uid []int
		db.Model(&entity.UserInfo{}).Where("nickname LIKE ?", "%"+nickname+"%").Pluck("id", &uid)
		if len(uid) > 0 {
			query = query.Where("user_id IN ?", uid)
		} else {
			query = query.Where("user_id = ?", 0) // no match
		}
	}

	if typ != 0 {
		query = query.Where("type = ?", typ)
	}
	if isReview != nil {
		query = query.Where("is_review = ?", *isReview)
	}

	err := query.Count(&total).
		Preload("User").Preload("User.UserInfo").
		Preload("ReplyUser").Preload("ReplyUser.UserInfo").
		Preload("Article").
		Order("id DESC").
		Offset((page - 1) * size).Limit(size).
		Find(&list).Error
	return list, total, err
}

func (r *interactionRepository) DeleteComments(db *gorm.DB, ids []int) error {
	return db.Where("id IN ?", ids).Delete(&entity.Comment{}).Error
}

func (r *interactionRepository) UpdateCommentsReview(db *gorm.DB, ids []int, isReview bool) error {
	return db.Model(&entity.Comment{}).Where("id IN ?", ids).Update("is_review", isReview).Error
}

func (r *interactionRepository) GetFrontCommentList(db *gorm.DB, page, size, topic, typ int) ([]entity.Comment, map[int][]entity.Comment, int64, error) {
	var list []entity.Comment
	var total int64

	tx := db.Model(&entity.Comment{})
	if typ != 0 {
		tx = tx.Where("type = ?", typ)
	}
	if topic != 0 {
		tx = tx.Where("topic_id = ?", topic)
	}

	err := tx.Where("parent_id = 0").
		Count(&total).
		Preload("User").Preload("User.UserInfo").
		Order("id DESC").
		Offset((page - 1) * size).Limit(size).
		Find(&list).Error

	if err != nil {
		return nil, nil, 0, err
	}

	replyMap := make(map[int][]entity.Comment)
	for i := range list {
		var replyList []entity.Comment
		db.Model(&entity.Comment{}).
			Where("parent_id = ?", list[i].ID).
			Preload("User").Preload("User.UserInfo").
			Preload("ReplyUser").Preload("ReplyUser.UserInfo").
			Order("id DESC").
			Find(&replyList)
		replyMap[list[i].ID] = replyList
	}

	return list, replyMap, total, nil
}

func (r *interactionRepository) GetCommentReplyList(db *gorm.DB, id, page, size int) ([]entity.Comment, error) {
	var data []entity.Comment
	err := db.Model(&entity.Comment{}).
		Where("parent_id = ?", id).
		Preload("User").Preload("User.UserInfo").
		Preload("ReplyUser").Preload("ReplyUser.UserInfo").
		Order("id DESC").
		Offset((page - 1) * size).Limit(size).
		Find(&data).Error
	return data, err
}

func (r *interactionRepository) AddComment(db *gorm.DB, comment *entity.Comment) error {
	return db.Create(comment).Error
}

func (r *interactionRepository) GetCommentById(db *gorm.DB, id int) (*entity.Comment, error) {
	var comment entity.Comment
	err := db.Preload("User").Preload("User.UserInfo").
		Preload("ReplyUser").Preload("ReplyUser.UserInfo").
		Preload("Article").
		First(&comment, id).Error
	return &comment, err
}
