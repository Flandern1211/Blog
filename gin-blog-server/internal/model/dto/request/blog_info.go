package request

type AboutReq struct {
	Content string `json:"content"`
}

type AddOrEditPageReq struct {
	ID    int    `json:"id"`
	Name  string `json:"name" binding:"required" validate:"required"`
	Label string `json:"label" binding:"required" validate:"required"`
	Cover string `json:"cover"`
}
