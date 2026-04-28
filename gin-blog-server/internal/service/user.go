package service

import (
	"encoding/json"
	"errors"
	global "gin-blog/internal/global"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/model/dto/response"
	"gin-blog/internal/model/entity"
	"gin-blog/internal/repository"
	"gin-blog/internal/utils"
	pkgErrors "gin-blog/pkg/errors"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService interface {
	GetInfo(c *gin.Context, authId int) (*response.UserInfoVO, error)
	UpdateCurrent(c *gin.Context, authId int, req request.UpdateCurrentUserReq) error
	Update(c *gin.Context, req request.UpdateUserReq) error
	UpdateDisable(c *gin.Context, req request.UpdateUserDisableReq) error
	GetList(c *gin.Context, query request.UserQuery) ([]response.UserVO, int64, error)
	GetOnlineList(c *gin.Context, keyword string) ([]*entity.UserAuth, error)
	ForceOffline(c *gin.Context, currentAuthId, targetUserId int) error
	UpdatePasswordByCode(c *gin.Context, authId int, req request.UpdatePasswordByCodeReq) error
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
	if err := s.repo.UpdateUserDisable(db, req.UserAuthId, req.IsDisable); err != nil {
		return err
	}

	// 禁用用户时同步强制下线
	if req.IsDisable {
		rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
		rctx := c.Request.Context()
		onlineKey := global.ONLINE_USER + strconv.Itoa(req.UserAuthId)
		offlineKey := global.OFFLINE_USER + strconv.Itoa(req.UserAuthId)
		rdb.Del(rctx, onlineKey)
		rdb.Set(rctx, offlineKey, "1", time.Hour)
	}

	return nil
}

func (s *userService) GetList(c *gin.Context, query request.UserQuery) ([]response.UserVO, int64, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	list, total, err := s.repo.GetList(db, query.GetPage(), query.GetSize(), query.LoginType, query.Nickname, query.Username)
	if err != nil {
		return nil, 0, err
	}

	var res []response.UserVO
	for _, user := range list {
		res = append(res, response.UserVO{
			ID:            user.ID,
			UserInfoId:    user.UserInfoId,
			Info:          user.UserInfo,
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

// 获取在线用户列表
func (s *userService) GetOnlineList(c *gin.Context, keyword string) ([]*entity.UserAuth, error) {
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()

	onlineList := make([]*entity.UserAuth, 0)
	keys := rdb.Keys(rctx, global.ONLINE_USER+"*").Val()

	for _, key := range keys {
		val, err := rdb.Get(rctx, key).Result()
		if err != nil || val == "" {
			continue
		}
		var auth entity.UserAuth
		if err := json.Unmarshal([]byte(val), &auth); err != nil {
			continue
		}

		if keyword != "" &&
			!strings.Contains(auth.Username, keyword) &&
			(auth.UserInfo == nil || !strings.Contains(auth.UserInfo.Nickname, keyword)) {
			continue
		}

		onlineList = append(onlineList, &auth)
	}

	// 根据最后登录时间排序
	sort.Slice(onlineList, func(i, j int) bool {
		if onlineList[i].LastLoginTime == nil || onlineList[j].LastLoginTime == nil {
			return false
		}
		return onlineList[i].LastLoginTime.Unix() > onlineList[j].LastLoginTime.Unix()
	})

	return onlineList, nil
}

// 强制用户下线
func (s *userService) ForceOffline(c *gin.Context, currentAuthId, targetUserId int) error {
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()

	onlineKey := global.ONLINE_USER + strconv.Itoa(targetUserId)
	offlineKey := global.OFFLINE_USER + strconv.Itoa(targetUserId)

	rdb.Del(rctx, onlineKey)
	rdb.Set(rctx, offlineKey, "1", time.Hour)

	return nil
}

// 前台用户通过邮箱验证码修改密码
func (s *userService) UpdatePasswordByCode(c *gin.Context, authId int, req request.UpdatePasswordByCodeReq) error {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()

	// 1. 获取当前用户信息
	userAuth, err := s.repo.GetInfoById(db, authId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return pkgErrors.NewDefault(pkgErrors.CodeUserNotFound)
		}
		return pkgErrors.NewDefault(pkgErrors.CodeDbOpError)
	}

	// 2. 验证邮箱是否匹配当前用户
	if userAuth.UserInfo == nil || userAuth.UserInfo.Email != req.Email {
		return pkgErrors.NewDefault(pkgErrors.CodeBadRequest)
	}

	// 3. 验证验证码
	codeKey := global.EMAIL_CODE + req.Email
	storedCode, err := rdb.Get(rctx, codeKey).Result()
	if err != nil || storedCode == "" {
		return pkgErrors.NewDefault(pkgErrors.CodeCodeWrong)
	}
	if storedCode != req.Code {
		return pkgErrors.NewDefault(pkgErrors.CodeCodeWrong)
	}

	// 验证码正确，删除已使用的验证码
	rdb.Del(rctx, codeKey)

	// 4. 加密新密码并更新
	hashedPassword, err := utils.BcryptHash(req.Password)
	if err != nil {
		return pkgErrors.NewDefault(pkgErrors.CodeInternalError)
	}

	return s.repo.UpdateUserPassword(db, authId, hashedPassword)
}
