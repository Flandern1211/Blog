package request

type PageQuery struct {
	Page int `form:"page" binding:"required,min=1"`
	Size int `form:"size" binding:"required,min=1,max=100"`
}
