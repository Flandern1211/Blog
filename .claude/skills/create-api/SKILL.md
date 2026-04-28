# 名称
create-api-go

# 描述
按照 Go 后端工程规范创建新的 API 端点 (Controller/Service/Repository)
disable-model-invocation: true
allowed-tools: Bash(go test ./...) Bash(go build ./...) Bash(golangci-lint run)
按照我们的规范创建新的 Go API 端点: $ARGUMENTS

## 必须遵循的步骤（不可跳过）

### 1. 检查现有模式

- 观察 internal/handler/ 或 internal/controller/ 中的代码，理解参数绑定和错误处理的统一封装（如 Response 结构体）。
- 确保新端点的错误处理（如 ErrUserNotFound）在 internal/errors/ 中统一定义。

### 2. 定义 DTO 与输入验证

- **文件**: internal/model/dto/{resource}.go
- **规范**: 使用结构体定义 Request/Response，并利用 validate 标签进行验证。
- **示例**:

    ```go
    type CreateRequest struct {
        Name  string `json:"name" validate:"required,min=3"`
        Email string `json:"email" validate:"required,email"`
    }
    ```


### 3. Controller 层 

- **文件**: internal/handler/{resource}_handler.go
- **职责**: 仅负责 **Bind**（参数绑定）、**Validate**（校验）和 **Render**（返回响应）。
- 必须通过 **Interface** 调用 Service 层，以便于 Mock。

### 4. Service 层 (业务逻辑)

- **文件**: internal/service/{resource}_service.go
- **职责**: 处理核心业务逻辑、权限判断和事务管理。
- 如果涉及数据库操作，调用 repository 层。

### 5. Repository 层 (数据持久化)

- **文件**: internal/repository/{resource}_repo.go
- **职责**: 封装 SQL 或 GORM 操作。严禁在 Service 层直接写原生 SQL。

### 6. 编写单元测试 (必须执行)

- **文件**: internal/service/{resource}_service_test.go
- **模式**: 必须使用 **表格驱动测试 (Table-Driven Tests)**。
- **Mock**: 使用 gomock 或 testify/mock 模拟数据库层，确保测试不依赖真实数据库。
- **覆盖率**: 至少包含 Success, InvalidInput, DataConflict 三种场景。

### 7. 更新接口文档

- 如果使用 Gin-Swagger，请在 Handler 函数上方补充 // @Summary, // @Router 等注释。
- 或手动更新项目根目录的 docs/openapi.yaml。

### 8. 自动化验证

- go vet ./... — 静态检查。
- golangci-lint run — 运行代码风格和潜在缺陷检查。
- go test -v -run Test{Resource} ./internal/service/... — 运行相关业务逻辑测试。