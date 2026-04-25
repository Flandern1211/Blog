# Swagger API 文档使用指南

## 概述

本项目使用 [swaggo/swag](https://github.com/swaggo/swag) 自动生成 OpenAPI 2.0（Swagger）规范的 API 文档。文档生成后可通过浏览器访问 Swagger UI 交互式页面，直接查看和测试所有 API 端点。

---

## 目录结构

```
gin-blog-server/docs/
├── docs.go        # Go 代码内嵌的 Swagger 规范（程序运行时直接使用）
├── swagger.json   # OpenAPI 2.0 JSON 格式
└── swagger.yaml   # OpenAPI 2.0 YAML 格式
```

三个文件内容等价，`docs.go` 将规范编译进二进制文件，运行时无需加载外部文件。

---

## 访问 Swagger UI

服务启动后，浏览器访问：

```
http://localhost:{port}/swagger/index.html
```

页面展示所有 API 端点列表，支持在线测试（填写参数 -> 点击 Execute -> 查看响应）。

---

## 使用流程

### 1. 在 Handler 代码中添加注释

在 Controller 的方法上按 swaggo 格式添加注释。示例：

```go
// Login 用户登录
// @Tags     UserAuth
// @Summary  登录
// @Description 用户登录
// @Accept   json
// @Produce  json
// @Param    form  body  request.LoginReq  true  "登录"
// @Success  0     {object}  handler.Response{data=model.LoginVO}
// @Router   /login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
    // ...
}
```

#### 常用注释说明

| 注释 | 说明 |
|------|------|
| `@Tags` | API 分类标签，相同标签归为一组 |
| `@Summary` | 接口简要说明 |
| `@Description` | 接口详细描述 |
| `@Accept` | 请求 Content-Type，如 `json`、`multipart/form-data` |
| `@Produce` | 响应 Content-Type |
| `@Param` | 参数：`参数名 参数类型 参数类型 是否必填 描述` |
| `@Success` | 成功响应：`HTTP状态码 响应类型 响应结构 描述` |
| `@Failure` | 失败响应 |
| `@Router` | 路由路径和 HTTP 方法：`路径 [方法]` |
| `@Security` | 认证方式，如 `ApiKeyAuth` |

#### 参数类型对照

| 写法 | 说明 |
|------|------|
| `body` | JSON 请求体 |
| `query` | URL 查询参数 |
| `path` | URL 路径参数 |
| `formData` | multipart/form-data 表单字段 |
| `header` | 请求头参数 |

#### 响应结构写法

本项目统一响应格式为 `handler.Response` 泛型结构：

```go
type Response struct {
    Code    int         `json:"code"`
    Data    interface{} `json:"data"`
    Message string      `json:"message"`
}
```

注释中响应结构的写法：

```go
// 返回 string 类型数据
@Success 0  {object}  handler.Response{data=string}

// 返回 model.LoginVO 对象
@Success 0  {object}  handler.Response{data=model.LoginVO}

// 返回分页结果
@Success 0  {object}  handler.Response{data=handler.PageResult-model.TagVO}

// 返回数组
@Success 0  {object}  handler.Response{data=[]model.CategoryVO}
```

### 2. 安装 swag CLI

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

确认安装成功：

```bash
swag --version
```

### 3. 生成文档

在 `gin-blog-server/` 目录下执行：

```bash
swag init -g internal/manager.go -o docs
```

- `-g` 指定入口文件（包含 `// @title`、`// @version` 等全局注释的文件）
- `-o` 指定输出目录（默认 `docs/`）

执行后自动更新 `docs/` 下的三个文件。

### 4. 启动服务验证

```go
// internal/manager.go 中已注册 Swagger 路由
docs.SwaggerInfo.BasePath = "/api"
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

启动服务，访问 `/swagger/index.html` 确认文档正确显示。

---

## 在 Gin 中的完整集成

### 依赖

```go
import (
    "gin-blog/docs"              // 自动生成的 docs 包
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)
```

### 注册

```go
// internal/manager.go
func RegisterHandlers(r *gin.Engine) {
    docs.SwaggerInfo.BasePath = "/api"          // 设置基础路径
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))  // 挂载 UI
    // ... 注册业务路由
}
```

---

## 常见问题

### 生成的 docs.go 提示 "DO NOT EDIT"

`docs.go` 是由 `swag init` 自动生成的，**不要手动编辑**。如需更新文档，修改 Controller 的注释后重新执行 `swag init`。

### API 在 Swagger UI 上不显示

可能的原因：
1. 没有添加对应的 `@Tags` 和 `@Router` 注释
2. `swag init` 后未重启服务
3. `@Router` 中的路径与 `BasePath` 拼接后不匹配

### 响应数据与实际返回不一致

检查 `@Success` 中声明的结构体是否与 Controller 中实际返回的类型一致。本项目统一使用 `handler.Response{}` 包装，`data` 字段需要正确指定类型。

### 认证接口报 401

在 Swagger UI 页面上方点击 **Authorize** 按钮，输入认证 Token，后续请求会自动附带。

---

## 最佳实践

1. **新增 API 后立即更新文档** — 项目红线规定"新增 API 端点必须同时更新 @docs 中的文档"
2. **保持 Controller 注释与代码一致** — 参数名、类型、是否必填等信息随代码变更同步更新
3. **运行 `swag init` 后重启服务** — 否则浏览器看到的还是旧文档
