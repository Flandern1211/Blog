package request

type UpdateCurrentUserReq struct {
	Nickname string `json:"nickname" binding:"required" validate:"required"`
	Avatar   string `json:"avatar"`
	Intro    string `json:"intro"`
	Website  string `json:"website"`
	Email    string `json:"email"`
}

type UpdateCurrentPasswordReq struct {
	NewPassword string `json:"new_password" binding:"required,min=4,max=20" validate:"required,min=4,max=20"`
	OldPassword string `json:"old_password" binding:"required,min=4,max=20" validate:"required,min=4,max=20"`
}

type UpdateUserReq struct {
	UserAuthId int    `json:"id"`
	Nickname   string `json:"nickname" binding:"required" validate:"required"`
	RoleIds    []int  `json:"role_ids"`
}

type UpdateUserDisableReq struct {
	UserAuthId int  `json:"id"`
	IsDisable  bool `json:"is_disable"`
}

type UserQuery struct {
	PageQuery
	LoginType int8   `form:"login_type"`
	Username  string `form:"username"`
	Nickname  string `form:"nickname"`
}

type ForceOfflineReq struct {
	UserInfoId int `json:"user_info_id"`
}

// 前台用户通过验证码修改密码
type UpdatePasswordByCodeReq struct {
	Email    string `json:"email" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Password string `json:"password" binding:"required,min=4,max=20"`
}
