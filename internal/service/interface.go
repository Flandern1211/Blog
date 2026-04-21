package service

import (
	"blog/internal/model/dto/request"
	"blog/internal/model/dto/response"
	"blog/internal/model/entity"
)

// UserService 用户服务接口
type UserService interface {
	// Register 用户注册
	Register(req *request.RegisterRequest) error
	// GetUserByID 根据 ID 获取用户
	GetUserByID(id uint) (*entity.User, error)
	// UpdateUser 更新用户信息
	UpdateUser(id uint, req *request.UpdateUserRequest) error
	// ChangePassword 修改密码
	ChangePassword(id uint, req *request.ChangePasswordRequest) error
	// GetUserResponse 获取用户响应
	GetUserResponse(user *entity.User) *response.UserResponse
}

// AuthService 认证服务接口
type AuthService interface {
	// Login 用户登录
	Login(req *request.LoginRequest) (*response.LoginResponse, error)
	// RefreshToken 刷新 Token
	RefreshToken(token string) (string, error)
}
