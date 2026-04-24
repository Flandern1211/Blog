重构项目后端·实现的整体结构，使其符合MVC架构，下面的为主要模块

**示例文件树**

```bash
├─api
│  └─v1
│      ├─admin
│      ├─auth
│      ├─operationLogs
│      ├─permissionManage
│      │  ├─menu
│      │  ├─permission
│      │  └─role
│      ├─queue
│ 
├─app
├─middleware
├─model
│  ├─dto
│  │  ├─request
│  │  └─response
│  └─entity
├─repository
└─service
```

### 1. model层

- **文件:**
- **内容**：只保留数据对象的定义，不直接进行和数据库的交互
- **示例：只有对UserAuth的定义**

```bash
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

```

### 2.  DTO 与输入验证

- **文件**: gin-blog-server/internal/model/dto
- 内容：只包含对Respnse和Request的定义
- **规范**: 使用结构体定义 Request/Response，并利用 validate 标签进行验证。
- **示例**:

    ```go
    type CreateRequest struct {
        Name  string `json:"name" validate:"required,min=3"`
        Email string `json:"email" validate:"required,email"`
    }
    ```


### 3. Controller 层

- **例子**: gin-blog-server/internal/api/v1/auth/controller.go
- **职责**: 仅负责 **Bind**（参数绑定）、**Validate**（校验）和 **Render**（返回响应）。
- 必须通过 **Interface** 调用 Service 层，以便于 Mock。

### 4. Service 层 (业务逻辑)

- **文件**: gin-blog-server/internal/service
- **职责**: 处理核心业务逻辑、权限判断和事务管理。
- 如果涉及数据库操作，调用 repository 层。

### 5. Repository 层 (数据持久化)

- **文件**:gin-blog-server/internal/respository
- **职责**: 封装 SQL 或 GORM 操作。严禁在 Service 层直接写原生 SQL。