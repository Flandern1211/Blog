package entity

import "time"

type UserAuth struct {
	Model
	Username      string     `gorm:"unique;type:varchar(50)" json:"username"`
	Password      string     `gorm:"type:varchar(100)" json:"-"`
	LoginType     int        `gorm:"type:tinyint(1);comment:登录类型" json:"login_type"`
	IpAddress     string     `gorm:"type:varchar(20);comment:登录IP地址" json:"ip_address"`
	IpSource      string     `gorm:"type:varchar(50);comment:IP来源" json:"ip_source"`
	LastLoginTime *time.Time `json:"last_login_time"`
	IsDisable     bool       `json:"is_disable"`
	IsSuper       bool       `json:"is_super"` // 超级管理员只能后台设置

	UserInfoId int       `json:"user_info_id"`
	UserInfo   *UserInfo `json:"info"`
	Roles      []*Role   `json:"roles" gorm:"many2many:user_auth_role"`
}

type UserInfo struct {
	Model
	Email     string `gorm:"type:varchar(50)" json:"email"`
	Nickname  string `gorm:"type:varchar(50)" json:"nickname"`
	Avatar    string `gorm:"type:varchar(255)" json:"avatar"`
	Intro     string `gorm:"type:varchar(255)" json:"intro"`
	Website   string `gorm:"type:varchar(255)" json:"website"`
	IsDisable bool   `json:"is_disable"`
}

type Role struct {
	Model
	Name      string `gorm:"unique" json:"name"`
	Label     string `gorm:"unique" json:"label"`
	IsDisable bool   `json:"is_disable"`

	Resources []Resource `json:"resources" gorm:"many2many:role_resource"`
	Menus     []Menu     `json:"menus" gorm:"many2many:role_menu"`
	Users     []UserAuth `json:"users" gorm:"many2many:user_auth_role"`
}

type Resource struct {
	Model
	Name      string `gorm:"unique;type:varchar(50)" json:"name"`
	ParentId  int    `json:"parent_id"`
	Url       string `gorm:"type:varchar(255)" json:"url"`
	Method    string `gorm:"type:varchar(10)" json:"request_method"`
	Anonymous bool   `json:"is_anonymous"`

	Roles []*Role `json:"roles" gorm:"many2many:role_resource"`
}

type Menu struct {
	Model
	ParentId     int    `json:"parent_id"`
	Name         string `gorm:"uniqueIndex:idx_name_and_path;type:varchar(20)" json:"name"`
	Path         string `gorm:"uniqueIndex:idx_name_and_path;type:varchar(50)" json:"path"`
	Component    string `gorm:"type:varchar(50)" json:"component"`
	Icon         string `gorm:"type:varchar(50)" json:"icon"`
	OrderNum     int8   `json:"order_num"`
	Redirect     string `gorm:"type:varchar(50)" json:"redirect"`
	Catalogue    bool   `json:"is_catalogue"`
	Hidden       bool   `json:"is_hidden"`
	KeepAlive    bool   `json:"keep_alive"`
	External     bool   `json:"is_external"`
	ExternalLink string `gorm:"type:varchar(255)" json:"external_link"`

	Roles []*Role `json:"roles" gorm:"many2many:role_menu"`
}

type RoleResource struct {
	RoleId     int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_resource"`
	ResourceId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_resource"`
}

type UserAuthRole struct {
	UserAuthId int `gorm:"primaryKey;uniqueIndex:idx_user_auth_role"`
	RoleId     int `gorm:"primaryKey;uniqueIndex:idx_user_auth_role"`
}

type RoleMenu struct {
	RoleId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_menu"`
	MenuId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_menu"`
}
