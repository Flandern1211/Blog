package request

type PageQuery struct {
	Page     int `form:"page"`
	Size     int `form:"size"`
	PageNum  int `form:"page_num"`
	PageSize int `form:"page_size"`
}

// GetPage 返回有效的页码，兼容 page 和 page_num 参数
func (p *PageQuery) GetPage() int {
	if p.Page > 0 {
		return p.Page
	}
	if p.PageNum > 0 {
		return p.PageNum
	}
	return 0
}

// GetSize 返回有效的每页条数，兼容 size 和 page_size 参数
func (p *PageQuery) GetSize() int {
	if p.Size > 0 {
		return p.Size
	}
	if p.PageSize > 0 {
		return p.PageSize
	}
	return 0
}
