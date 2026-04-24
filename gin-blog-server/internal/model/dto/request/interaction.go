package request

type MessageQuery struct {
	PageQuery
	Nickname string `form:"nickname"`
	IsReview *bool  `form:"is_review"`
}

type UpdateReviewReq struct {
	Ids      []int `json:"ids"`
	IsReview bool  `json:"is_review"`
}

type CommentQuery struct {
	PageQuery
	Type     int    `form:"type"`
	IsReview *bool  `form:"is_review"`
	Nickname string `form:"nickname"`
}
