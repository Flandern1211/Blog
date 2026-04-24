package request

type AddOrEditRoleReq struct {
	ID          int    `json:"id"`
	Name        string `json:"name" binding:"required" validate:"required"`
	Label       string `json:"label" binding:"required" validate:"required"`
	IsDisable   bool   `json:"is_disable"`
	ResourceIds []int  `json:"resource_ids"`
	MenuIds     []int  `json:"menu_ids"`
}

type AddOrEditMenuReq struct {
	ID           int    `json:"id"`
	ParentId     int    `json:"parent_id"`
	Name         string `json:"name" binding:"required" validate:"required"`
	Path         string `json:"path" binding:"required" validate:"required"`
	Component    string `json:"component" binding:"required" validate:"required"`
	Icon         string `json:"icon"`
	OrderNum     int8   `json:"order_num"`
	Redirect     string `json:"redirect"`
	Catalogue    bool   `json:"is_catalogue"`
	Hidden       bool   `json:"is_hidden"`
	KeepAlive    bool   `json:"keep_alive"`
	External     bool   `json:"is_external"`
	ExternalLink string `json:"external_link"`
}

type AddOrEditResourceReq struct {
	ID       int    `json:"id"`
	Url      string `json:"url"`
	Method   string `json:"request_method"`
	Name     string `json:"name" binding:"required" validate:"required"`
	ParentId int    `json:"parent_id"`
}

type EditAnonymousReq struct {
	ID        int  `json:"id" binding:"required" validate:"required"`
	Anonymous bool `json:"is_anonymous"`
}
