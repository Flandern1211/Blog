package service

import (
	global "gin-blog/internal/global"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/model/dto/response"
	"gin-blog/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService interface {
	GetInfo(c *gin.Context, authId int) (*response.UserInfoVO, error)
	UpdateCurrent(c *gin.Context, authId int, req request.UpdateCurrentUserReq) error
	Update(c *gin.Context, req request.UpdateUserReq) error
	UpdateDisable(c *gin.Context, req request.UpdateUserDisableReq) error
	GetList(c *gin.Context, query request.UserQuery) ([]response.UserVO, int64, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetInfo(c *gin.Context, authId int) (*response.UserInfoVO, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()

	userAuth, err := s.repo.GetInfoById(db, authId)
	if err != nil {
		return nil, err
	}

	articleLikeSet, _ := rdb.SMembers(rctx, global.ARTICLE_USER_LIKE_SET+strconv.Itoa(authId)).Result()
	commentLikeSet, _ := rdb.SMembers(rctx, global.COMMENT_USER_LIKE_SET+strconv.Itoa(authId)).Result()

	return &response.UserInfoVO{
		UserInfo:       *userAuth.UserInfo,
		ArticleLikeSet: articleLikeSet,
		CommentLikeSet: commentLikeSet,
	}, nil
}

func (s *userService) UpdateCurrent(c *gin.Context, authId int, req request.UpdateCurrentUserReq) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	userAuth, err := s.repo.GetInfoById(db, authId)
	if err != nil {
		return err
	}
	return s.repo.UpdateUserInfo(db, userAuth.UserInfoId, req.Nickname, req.Avatar, req.Intro, req.Website)
}

func (s *userService) Update(c *gin.Context, req request.UpdateUserReq) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.UpdateUserNicknameAndRole(db, req.UserAuthId, req.Nickname, req.RoleIds)
}

func (s *userService) UpdateDisable(c *gin.Context, req request.UpdateUserDisableReq) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	return s.repo.UpdateUserDisable(db, req.UserAuthId, req.IsDisable)
}

func (s *userService) GetList(c *gin.Context, query request.UserQuery) ([]response.UserVO, int64, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	list, total, err := s.repo.GetList(db, query.Page, query.Size, query.LoginType, query.Nickname, query.Username)
	if err != nil {
		return nil, 0, err
	}

	var res []response.UserVO
	for _, user := range list {
		res = append(res, response.UserVO{
			ID:            user.ID,
			UserInfoId:    user.UserInfoId,
			Avatar:        user.UserInfo.Avatar,
			Nickname:      user.UserInfo.Nickname,
			Roles:         user.Roles,
			LoginType:     user.LoginType,
			IpAddress:     user.IpAddress,
			IpSource:      user.IpSource,
			CreatedAt:     user.CreatedAt,
			LastLoginTime: user.LastLoginTime,
			IsDisable:     user.IsDisable,
		})
	}
	return res, total, nil
}
