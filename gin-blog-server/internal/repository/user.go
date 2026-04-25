package repository

import (
	"gin-blog/internal/model/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetInfoById(db *gorm.DB, id int) (*entity.UserAuth, error)
	UpdateUserInfo(db *gorm.DB, id int, nickname, avatar, intro, website string) error
	UpdateUserPassword(db *gorm.DB, id int, password string) error
	UpdateUserNicknameAndRole(db *gorm.DB, authId int, nickname string, roleIds []int) error
	UpdateUserDisable(db *gorm.DB, id int, isDisable bool) error
	GetList(db *gorm.DB, page, size int, loginType int8, nickname, username string) ([]entity.UserAuth, int64, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) GetInfoById(db *gorm.DB, id int) (*entity.UserAuth, error) {
	var userAuth entity.UserAuth
	err := db.Model(&userAuth).
		Preload("Roles").Preload("UserInfo").
		First(&userAuth, id).Error
	return &userAuth, err
}

func (r *userRepository) UpdateUserInfo(db *gorm.DB, id int, nickname, avatar, intro, website string) error {
	return db.Model(&entity.UserInfo{Model: entity.Model{ID: id}}).
		Select("nickname", "avatar", "intro", "website").
		Updates(entity.UserInfo{
			Nickname: nickname,
			Avatar:   avatar,
			Intro:    intro,
			Website:  website,
		}).Error
}

func (r *userRepository) UpdateUserPassword(db *gorm.DB, id int, password string) error {
	return db.Model(&entity.UserAuth{Model: entity.Model{ID: id}}).
		Update("password", password).Error
}

func (r *userRepository) UpdateUserNicknameAndRole(db *gorm.DB, authId int, nickname string, roleIds []int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var userAuth entity.UserAuth
		if err := tx.First(&userAuth, authId).Error; err != nil {
			return err
		}

		if err := tx.Model(&entity.UserInfo{Model: entity.Model{ID: userAuth.UserInfoId}}).
			Update("nickname", nickname).Error; err != nil {
			return err
		}

		if len(roleIds) > 0 {
			if err := tx.Delete(&entity.UserAuthRole{}, "user_auth_id = ?", authId).Error; err != nil {
				return err
			}
			var userRoles []entity.UserAuthRole
			for _, rid := range roleIds {
				userRoles = append(userRoles, entity.UserAuthRole{
					UserAuthId: authId,
					RoleId:     rid,
				})
			}
			if err := tx.Create(&userRoles).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *userRepository) UpdateUserDisable(db *gorm.DB, id int, isDisable bool) error {
	return db.Model(&entity.UserAuth{Model: entity.Model{ID: id}}).
		Update("is_disable", isDisable).Error
}

func (r *userRepository) GetList(db *gorm.DB, page, size int, loginType int8, nickname, username string) ([]entity.UserAuth, int64, error) {
	var list []entity.UserAuth
	var total int64

	query := db.Model(&entity.UserAuth{}).
		Joins("LEFT JOIN user_info ON user_info.id = user_auth.user_info_id")

	if loginType != 0 {
		query = query.Where("user_auth.login_type = ?", loginType)
	}
	if username != "" {
		query = query.Where("user_auth.username LIKE ?", "%"+username+"%")
	}
	if nickname != "" {
		query = query.Where("user_info.nickname LIKE ?", "%"+nickname+"%")
	}

	err := query.Count(&total).
		Preload("UserInfo").
		Preload("Roles").
		Scopes(Paginate(page, size)).
		Find(&list).Error

	return list, total, err
}
