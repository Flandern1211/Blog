package request

type LoginReq struct {
	Username string `json:"username" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"required"`
}

type RegisterReq struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required,min=4,max=20" validate:"required,min=4,max=20"`
}

type SendCodeReq struct {
	Email string `json:"email" binding:"required,email" validate:"required,email"`
}
