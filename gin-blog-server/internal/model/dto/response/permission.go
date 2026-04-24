package response

import (
	"gin-blog/internal/model/entity"
	"time"
)

type RoleVO struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Label       string    `json:"label"`
	IsDisable   bool      `json:"is_disable"`
	CreatedAt   time.Time `json:"created_at"`
	ResourceIds []int     `json:"resource_ids" gorm:"-"`
	MenuIds     []int     `json:"menu_ids" gorm:"-"`
}

type MenuTreeVO struct {
	entity.Menu
	Children []MenuTreeVO `json:"children"`
}

type ResourceTreeVO struct {
	ID        int              `json:"id"`
	CreatedAt time.Time        `json:"created_at"`
	Name      string           `json:"name"`
	Url       string           `json:"url"`
	Method    string           `json:"request_method"`
	Anonymous bool             `json:"is_anonymous"`
	Children  []ResourceTreeVO `json:"children"`
}

type TreeOptionVO struct {
	ID       int            `json:"key"`
	Label    string         `json:"label"`
	Children []TreeOptionVO `json:"children"`
}

type OptionVO struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}
