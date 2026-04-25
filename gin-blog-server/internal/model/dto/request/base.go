package request

type PageQuery struct {
	Page int `form:"page"`
	Size int `form:"size"`
}
