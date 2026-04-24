package service

import (
	"errors"
	global "gin-blog/internal/global"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/model/dto/response"
	"gin-blog/internal/repository"
	"gin-blog/internal/utils"
	"gin-blog/internal/utils/jwt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(c *gin.Context, req request.LoginReq) (*response.LoginVO, error)
	Register(c *gin.Context, req request.RegisterReq) error
	VerifyCode(c *gin.Context, code string) error
	Logout(c *gin.Context, authId int) error
	SendCode(c *gin.Context, email string) error
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Login(c *gin.Context, req request.LoginReq) (*response.LoginVO, error) {
	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)

	userAuth, err := s.repo.GetUserAuthInfoByName(db, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.ErrUserNotExist
		}
		return nil, global.ErrDbOp
	}

	if !utils.BcryptCheck(req.Password, userAuth.Password) {
		return nil, global.ErrPassword
	}

	ipAddress := utils.IP.GetIpAddress(c)
	ipSource := utils.IP.GetIpSourceSimpleIdle(ipAddress)

	userInfo, err := s.repo.GetUserInfoById(db, userAuth.UserInfoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, global.ErrUserNotExist
		}
		return nil, global.ErrDbOp
	}

	roleIds, err := s.repo.GetRoleIdsByUserId(db, userAuth.ID)
	if err != nil {
		return nil, global.ErrDbOp
	}

	rctx := c.Request.Context()
	articleLikeSet, err := rdb.SMembers(rctx, global.ARTICLE_USER_LIKE_SET+strconv.Itoa(userAuth.ID)).Result()
	if err != nil {
		return nil, global.ErrRedisOp
	}
	commentLikeSet, err := rdb.SMembers(rctx, global.COMMENT_USER_LIKE_SET+strconv.Itoa(userAuth.ID)).Result()
	if err != nil {
		return nil, global.ErrRedisOp
	}

	conf := global.Conf.JWT
	token, err := jwt.GenToken(conf.Secret, conf.Issuer, int(conf.Expire), userAuth.ID, roleIds)
	if err != nil {
		return nil, global.ErrTokenCreate
	}

	err = s.repo.UpdateUserLoginInfo(db, userAuth.ID, ipAddress, ipSource)
	if err != nil {
		return nil, global.ErrDbOp
	}

	offlineKey := global.OFFLINE_USER + strconv.Itoa(userAuth.ID)
	rdb.Del(rctx, offlineKey).Result()

	return &response.LoginVO{
		UserInfo:       *userInfo,
		ArticleLikeSet: articleLikeSet,
		CommentLikeSet: commentLikeSet,
		Token:          token,
	}, nil
}

func (s *authService) Logout(c *gin.Context, authId int) error {
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()
	onlineKey := global.ONLINE_USER + strconv.Itoa(authId)
	rdb.Del(rctx, onlineKey)
	return nil
}

func (s *authService) SendCode(c *gin.Context, email string) error {
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()

	code := utils.RandomCode(6)
	// 存入 redis，有效期 15 分钟
	err := rdb.Set(rctx, global.EMAIL_CODE+email, code, 15*time.Minute).Err()
	if err != nil {
		return global.ErrRedisOp
	}

	// 发送邮件
	err = utils.SendCodeEmail(email, &utils.EmailData{
		UserName: email,
		Subject:  "注册验证码",
		Code:     code,
	})
	if err != nil {
		return global.ErrSendEmail
	}

	return nil
}

func (s *authService) Register(c *gin.Context, req request.RegisterReq) error {
	req.Username = utils.Format(req.Username)
	db := c.MustGet(global.CTX_DB).(*gorm.DB)

	auth, err := s.repo.GetUserAuthInfoByName(db, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return global.ErrDbOp
	}
	if auth != nil {
		return global.ErrUserExist
	}

	info := utils.GenEmailVerificationInfo(req.Username, req.Password)
	// Wait, original code uses SetMailInfo, we need to adapt that, maybe via rdb directly?
	// For simplicity, assuming SetMailInfo is available in utils or we rewrite it here:
	// utils.SetMailInfo(rdb, info, 15*time.Minute)
	// But let's check utils to see if it's there. The original code uses `SetMailInfo(GetRDB(c), info, 15*time.Minute)`.
	// For now, I'll use the original `SetMailInfo` if it's in `handle` package, I need to copy it or adapt.
	// Let's implement it here.
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()
	rdb.Set(rctx, info, info, 15*time.Minute)

	EmailData := utils.GetEmailData(req.Username, info)
	err = utils.SendEmail(req.Username, EmailData)
	if err != nil {
		return global.ErrSendEmail
	}
	return nil
}

func (s *authService) VerifyCode(c *gin.Context, code string) error {
	rdb := c.MustGet(global.CTX_RDB).(*global.RedisClient)
	rctx := c.Request.Context()

	val, err := rdb.Get(rctx, code).Result()
	if err != nil || val == "" {
		return errors.New("code not exist")
	}
	rdb.Del(rctx, code)

	username, password, err := utils.ParseEmailVerificationInfo(code)
	if err != nil {
		return err
	}

	db := c.MustGet(global.CTX_DB).(*gorm.DB)
	_, _, _, err = s.repo.CreateNewUser(db, username, password)
	if err != nil {
		return err
	}

	return nil
}
