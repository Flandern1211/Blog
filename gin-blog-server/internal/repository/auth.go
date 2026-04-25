package repository

import (
	"gin-blog/internal/model/entity"
	"strconv"

	"gorm.io/gorm"
)

type AuthRepository interface {
	GetUserAuthInfoByName(db *gorm.DB, username string) (*entity.UserAuth, error)
	GetUserInfoById(db *gorm.DB, id int) (*entity.UserInfo, error)
	GetRoleIdsByUserId(db *gorm.DB, userId int) ([]int, error)
	UpdateUserLoginInfo(db *gorm.DB, userId int, ipAddress, ipSource string) error
	CreateNewUser(db *gorm.DB, username, email, password string) (*entity.UserAuth, *entity.UserInfo, *entity.UserAuthRole, error)
	GetUserAuthInfoById(db *gorm.DB, id int) (*entity.UserAuth, error)
	GetResource(db *gorm.DB, url, method string) (*entity.Resource, error)
	CheckRoleAuth(db *gorm.DB, roleId int, url, method string) (bool, error)
}

type authRepository struct{}

func NewAuthRepository() AuthRepository {
	return &authRepository{}
}

func (r *authRepository) GetUserAuthInfoByName(db *gorm.DB, username string) (*entity.UserAuth, error) {
	var userAuth entity.UserAuth
	result := db.Where("username = ?", username).First(&userAuth)
	return &userAuth, result.Error
}

func (r *authRepository) GetUserInfoById(db *gorm.DB, id int) (*entity.UserInfo, error) {
	var userInfo entity.UserInfo
	result := db.First(&userInfo, id)
	return &userInfo, result.Error
}

func (r *authRepository) GetRoleIdsByUserId(db *gorm.DB, userId int) ([]int, error) {
	var ids []int
	result := db.Model(&entity.UserAuthRole{UserAuthId: userId}).Pluck("role_id", &ids)
	return ids, result.Error
}

func (r *authRepository) UpdateUserLoginInfo(db *gorm.DB, userId int, ipAddress, ipSource string) error {
	return db.Model(&entity.UserAuth{Model: entity.Model{ID: userId}}).
		Updates(map[string]interface{}{
			"ip_address": ipAddress,
			"ip_source":  ipSource,
		}).Error
}

func (r *authRepository) CreateNewUser(db *gorm.DB, username, email, password string) (*entity.UserAuth, *entity.UserInfo, *entity.UserAuthRole, error) {
	var num int64
	db.Model(&entity.UserInfo{}).Count(&num)
	number := strconv.FormatInt(num, 10)

	userInfo := &entity.UserInfo{
		Email:    email,
		Nickname: "游客" + number,
		Avatar:   "https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png",
		Intro:    "我是这个程序的第" + number + "个用户",
	}

	var userAuth *entity.UserAuth
	var userRole *entity.UserAuthRole

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(userInfo).Error; err != nil {
			return err
		}

		userAuth = &entity.UserAuth{
			Username:   username,
			Password:   password,
			UserInfoId: userInfo.ID,
		}
		if err := tx.Create(userAuth).Error; err != nil {
			return err
		}

		userRole = &entity.UserAuthRole{
			UserAuthId: userAuth.ID,
			RoleId:     2, // 默认身份为游客
		}
		if err := tx.Create(userRole).Error; err != nil {
			return err
		}

		return nil
	})

	return userAuth, userInfo, userRole, err
}

func (r *authRepository) GetUserAuthInfoById(db *gorm.DB, id int) (*entity.UserAuth, error) {
	var userAuth entity.UserAuth
	err := db.Preload("Roles").Preload("UserInfo").First(&userAuth, id).Error
	return &userAuth, err
}

func (r *authRepository) GetResource(db *gorm.DB, url, method string) (*entity.Resource, error) {
	var resource entity.Resource
	err := db.Where("url = ? AND method = ?", url, method).First(&resource).Error
	return &resource, err
}

func (r *authRepository) CheckRoleAuth(db *gorm.DB, roleId int, url, method string) (bool, error) {
	var role entity.Role
	if err := db.Preload("Resources").First(&role, roleId).Error; err != nil {
		return false, err
	}

	for _, res := range role.Resources {
		if res.Anonymous || (res.Url == url && res.Method == method) {
			return true, nil
		}
	}

	return false, nil
}
