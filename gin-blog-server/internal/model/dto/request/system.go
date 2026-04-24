package request

type FriendLinkQuery struct {
	PageQuery
	Keyword string `form:"keyword"`
}

type AddOrEditLinkReq struct {
	ID      int    `json:"id"`
	Name    string `json:"name" binding:"required" validate:"required"`
	Avatar  string `json:"avatar"`
	Address string `json:"address" binding:"required" validate:"required"`
	Intro   string `json:"intro"`
}

type OperationLogQuery struct {
	PageQuery
	Keyword string `form:"keyword"`
}
