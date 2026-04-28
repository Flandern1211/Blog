package service

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/model/dto/request"
	"gin-blog/internal/model/dto/response"
	"gin-blog/internal/model/entity"
	"gin-blog/internal/repository"
	"gin-blog/internal/utils"
	pkgErrors "gin-blog/pkg/errors"
	"gin-blog/pkg/jwt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(c *gin.Context, req request.LoginReq) (*response.LoginVO, error)
	AdminLogin(c *gin.Context, req request.LoginReq) (*response.LoginVO, error)
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

func (s *authService) doLogin(c *gin.Context, req request.LoginReq) (*entity.UserAuth, *entity.UserInfo, []int, string, string, error) {
	db := c.MustGet(g.CTX_DB).(*gorm.DB)

	userAuth, err := s.repo.GetUserAuthInfoByName(db, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, "", "", pkgErrors.NewDefault(pkgErrors.CodeUserNotFound)
		}
		return nil, nil, nil, "", "", pkgErrors.NewDefault(pkgErrors.CodeDbOpError)
	}

	if userAuth.IsDisable {
		return nil, nil, nil, "", "", pkgErrors.NewDefault(pkgErrors.CodeUserDisabled)
	}

	if !utils.BcryptCheck(req.Password, userAuth.Password) {
		return nil, nil, nil, "", "", pkgErrors.NewDefault(pkgErrors.CodeInvalidCredentials)
	}

	ipAddress := utils.IP.GetIpAddress(c)
	ipSource := utils.IP.GetIpSourceSimpleIdle(ipAddress)

	userInfo, err := s.repo.GetUserInfoById(db, userAuth.UserInfoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil, "", "", pkgErrors.NewDefault(pkgErrors.CodeUserNotFound)
		}
		return nil, nil, nil, "", "", pkgErrors.NewDefault(pkgErrors.CodeDbOpError)
	}

	roleIds, err := s.repo.GetRoleIdsByUserId(db, userAuth.ID)
	if err != nil {
		return nil, nil, nil, "", "", pkgErrors.NewDefault(pkgErrors.CodeDbOpError)
	}

	return userAuth, userInfo, roleIds, ipAddress, ipSource, nil
}

func (s *authService) buildLoginVO(c *gin.Context, userAuth *entity.UserAuth, userInfo *entity.UserInfo, roleIds []int, ipAddress, ipSource string) (*response.LoginVO, error) {
	rdb := c.MustGet(g.CTX_RDB).(*g.RedisClient)
	rctx := c.Request.Context()

	articleLikeSet, err := rdb.SMembers(rctx, g.ARTICLE_USER_LIKE_SET+strconv.Itoa(userAuth.ID)).Result()
	if err != nil {
		return nil, pkgErrors.NewDefault(pkgErrors.CodeRedisOpError)
	}
	commentLikeSet, err := rdb.SMembers(rctx, g.COMMENT_USER_LIKE_SET+strconv.Itoa(userAuth.ID)).Result()
	if err != nil {
		return nil, pkgErrors.NewDefault(pkgErrors.CodeRedisOpError)
	}

	conf := g.Conf.JWT
	token, err := jwt.GenerateToken(conf.Secret, conf.Issuer, int(conf.Expire), userAuth.ID, roleIds)
	if err != nil {
		return nil, pkgErrors.NewDefault(pkgErrors.CodeTokenCreateErr)
	}

	db := c.MustGet(g.CTX_DB).(*gorm.DB)
	err = s.repo.UpdateUserLoginInfo(db, userAuth.ID, ipAddress, ipSource)
	if err != nil {
		return nil, pkgErrors.NewDefault(pkgErrors.CodeDbOpError)
	}

	rdb.Del(rctx, g.OFFLINE_USER+strconv.Itoa(userAuth.ID)).Result()

	return &response.LoginVO{
		UserInfo:       *userInfo,
		ArticleLikeSet: articleLikeSet,
		CommentLikeSet: commentLikeSet,
		Token:          token,
		IsSuper:        userAuth.IsSuper,
	}, nil
}

func (s *authService) Login(c *gin.Context, req request.LoginReq) (*response.LoginVO, error) {
	userAuth, userInfo, roleIds, ipAddress, ipSource, err := s.doLogin(c, req)
	if err != nil {
		return nil, err
	}

	return s.buildLoginVO(c, userAuth, userInfo, roleIds, ipAddress, ipSource)
}

func (s *authService) AdminLogin(c *gin.Context, req request.LoginReq) (*response.LoginVO, error) {
	userAuth, userInfo, roleIds, ipAddress, ipSource, err := s.doLogin(c, req)
	if err != nil {
		return nil, err
	}

	// 非超级管理员需要校验是否有后台登录权限
	if !userAuth.IsSuper {
		db := c.MustGet(g.CTX_DB).(*gorm.DB)
		hasResource, err := s.repo.CheckUserHasResource(db, userAuth.ID, g.RESOURCE_BACKEND_LOGIN, g.METHOD_BACKEND_LOGIN)
		if err != nil {
			return nil, pkgErrors.NewDefault(pkgErrors.CodeDbOpError)
		}
		if !hasResource {
			return nil, pkgErrors.NewDefault(pkgErrors.CodeNoAdminAccess)
		}
	}

	return s.buildLoginVO(c, userAuth, userInfo, roleIds, ipAddress, ipSource)
}

func (s *authService) Logout(c *gin.Context, authId int) error {
	rdb := c.MustGet(g.CTX_RDB).(*g.RedisClient)
	rctx := c.Request.Context()
	onlineKey := g.ONLINE_USER + strconv.Itoa(authId)
	rdb.Del(rctx, onlineKey)
	return nil
}

func (s *authService) SendCode(c *gin.Context, email string) error {
	rdb := c.MustGet(g.CTX_RDB).(*g.RedisClient)
	rctx := c.Request.Context()

	code := utils.RandomCode(6)
	// 存入 redis，有效期 15 分钟
	err := rdb.Set(rctx, g.EMAIL_CODE+email, code, 15*time.Minute).Err()
	if err != nil {
		return pkgErrors.NewDefault(pkgErrors.CodeRedisOpError)
	}

	// 发送邮件
	err = utils.SendCodeEmail(email, &utils.EmailData{
		UserName: email,
		Subject:  "注册验证码",
		Code:     code,
	})
	if err != nil {
		return pkgErrors.NewDefault(pkgErrors.CodeSendEmailErr)
	}

	return nil
}

func (s *authService) Register(c *gin.Context, req request.RegisterReq) error {
	req.Email = utils.Format(req.Email)
	db := c.MustGet(g.CTX_DB).(*gorm.DB)

	auth, err := s.repo.GetUserAuthInfoByName(db, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return pkgErrors.NewDefault(pkgErrors.CodeDbOpError)
	}
	if auth != nil {
		return pkgErrors.NewDefault(pkgErrors.CodeEmailExist)
	}

	info := utils.GenEmailVerificationInfo(req.Email, req.Password)
	rdb := c.MustGet(g.CTX_RDB).(*g.RedisClient)
	rctx := c.Request.Context()
	rdb.Set(rctx, info, info, 15*time.Minute)

	EmailData := utils.GetEmailData(req.Email, info)
	err = utils.SendEmail(req.Email, EmailData)
	if err != nil {
		return pkgErrors.NewDefault(pkgErrors.CodeSendEmailErr)
	}
	return nil
}

func (s *authService) VerifyCode(c *gin.Context, code string) error {
	rdb := c.MustGet(g.CTX_RDB).(*g.RedisClient)
	rctx := c.Request.Context()

	val, err := rdb.Get(rctx, code).Result()
	if err != nil || val == "" {
		return pkgErrors.NewDefault(pkgErrors.CodeCodeWrong)
	}
	rdb.Del(rctx, code)

	username, password, err := utils.ParseEmailVerificationInfo(code)
	if err != nil {
		return pkgErrors.NewDefault(pkgErrors.CodeCodeWrong)
	}

	db := c.MustGet(g.CTX_DB).(*gorm.DB)
	_, _, _, err = s.repo.CreateNewUser(db, username, username, password)
	if err != nil {
		return pkgErrors.NewDefault(pkgErrors.CodeDbOpError)
	}

	return nil
}
