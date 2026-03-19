# 后端开发规范

> 参考：Google Go Style Guide、Uber Go Style Guide、Microsoft REST API Guidelines

---

## 一、代码风格

### 1.1 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 包名 | 小写单词，不用下划线 | `handler`, `repository` |
| 文件名 | 小写+下划线 | `user_handler.go` |
| 接口 | 动词+er 或名词 | `UserRepository`, `Reader` |
| 结构体 | 大驼峰 | `UserService` |
| 函数/方法 | 大驼峰（导出）/小驼峰（私有） | `GetUser`, `validateToken` |
| 常量 | 大驼峰 | `StatusPublished` |
| 错误变量 | Err 前缀 | `ErrUserNotFound` |

```go
// Good
type UserRepository interface {
    GetByID(ctx context.Context, id uint) (*User, error)
}

// Bad
type USER_REPOSITORY interface {}
type user_repository interface {}
```

### 1.2 注释规范

```go
// CreateUser creates a new user with the given credentials.
// It returns ErrUsernameExists if the username is already taken.
//
// Example:
//   user, err := svc.CreateUser(ctx, "john", "password123")
func (s *UserService) CreateUser(ctx context.Context, username, password string) (*User, error) {
    // implementation
}
```

- 导出函数必须有注释
- 注释以函数名开头
- 复杂逻辑添加行内注释
- 避免冗余注释

### 1.3 文件组织

```go
// 导入顺序
import (
    // 标准库
    "context"
    "fmt"
    
    // 第三方库
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    
    // 内部包
    "social-app/internal/model"
    "social-app/pkg/logger"
)
```

---

## 二、错误处理

### 2.1 错误定义

```go
// internal/errors/errors.go
package errors

import "errors"

var (
    ErrUserNotFound      = errors.New("user not found")
    ErrUsernameExists    = errors.New("username already exists")
    ErrInvalidPassword   = errors.New("invalid password")
    ErrTokenExpired      = errors.New("token expired")
    ErrContentNotFound   = errors.New("content not found")
)

// AppError 业务错误，包含错误码
type AppError struct {
    Code    int
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return e.Message + ": " + e.Err.Error()
    }
    return e.Message
}

func NewAppError(code int, message string, err error) *AppError {
    return &AppError{Code: code, Message: message, Err: err}
}
```

### 2.2 错误传递

```go
// Good: 添加上下文
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
    var user model.User
    if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to get user by id %d: %w", id, err)
    }
    return &user, nil
}

// Bad: 吞掉错误
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
    var user model.User
    r.db.First(&user, id) // 错误被忽略
    return &user, nil
}
```

### 2.3 错误码规范

| 范围 | 模块 |
|------|------|
| 10000-10999 | Auth 认证 |
| 20000-20999 | User 用户 |
| 30000-30999 | Content 内容 |
| 40000-40999 | Comment 评论 |
| 50000-50999 | Like 点赞 |
| 90000-90999 | System 系统 |

---

## 三、日志规范

### 3.1 日志级别

| 级别 | 使用场景 |
|------|----------|
| DEBUG | 开发调试信息 |
| INFO | 正常业务流程 |
| WARN | 异常但可恢复 |
| ERROR | 需要关注的错误 |
| FATAL | 服务无法继续 |

### 3.2 结构化日志

```go
// Good: 结构化字段
logger.Info("user login",
    zap.Uint("user_id", userID),
    zap.String("username", username),
    zap.String("ip", c.ClientIP()),
)

// Bad: 字符串拼接
logger.Info(fmt.Sprintf("user %d login from %s", userID, c.ClientIP()))
```

### 3.3 敏感信息

```go
// 禁止记录
- 密码
- Token
- 身份证号
- 银行卡号

// 脱敏处理
logger.Info("user registered",
    zap.String("phone", maskPhone(phone)), // 138****8888
)
```

---

## 四、API 设计规范

### 4.1 RESTful 规范

| 操作 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 创建 | POST | /api/v1/contents | 创建资源 |
| 查询列表 | GET | /api/v1/contents | 分页查询 |
| 查询详情 | GET | /api/v1/contents/:id | 单个资源 |
| 更新 | PUT | /api/v1/contents/:id | 全量更新 |
| 部分更新 | PATCH | /api/v1/contents/:id | 部分更新 |
| 删除 | DELETE | /api/v1/contents/:id | 删除资源 |

### 4.2 统一响应格式

```go
// pkg/response/response.go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

type PagedData struct {
    List     interface{} `json:"list"`
    Total    int64       `json:"total"`
    Page     int         `json:"page"`
    PageSize int         `json:"page_size"`
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:    0,
        Message: "success",
        Data:    data,
    })
}

func Error(c *gin.Context, code int, message string) {
    c.JSON(http.StatusOK, Response{
        Code:    code,
        Message: message,
    })
}
```

### 4.3 分页规范

```go
// Query Parameters
// page: 页码，从 1 开始，默认 1
// page_size: 每页数量，默认 20，最大 50

type Pagination struct {
    Page     int `form:"page" binding:"min=1"`
    PageSize int `form:"page_size" binding:"min=1,max=50"`
}

func (p *Pagination) Normalize() {
    if p.Page < 1 {
        p.Page = 1
    }
    if p.PageSize < 1 || p.PageSize > 50 {
        p.PageSize = 20
    }
}

func (p *Pagination) Offset() int {
    return (p.Page - 1) * p.PageSize
}
```

### 4.4 版本控制

- URL 路径版本：`/api/v1/...`
- Breaking Change 时升级大版本：`/api/v2/...`
- 同一版本内保持向后兼容

---

## 五、数据库规范

### 5.1 事务使用

```go
// 需要事务的操作
func (s *ContentService) Publish(ctx context.Context, req *CreateContentReq) (*Content, error) {
    return s.db.Transaction(func(tx *gorm.DB) (interface{}, error) {
        // 1. 创建内容
        content := &model.Content{...}
        if err := tx.Create(content).Error; err != nil {
            return nil, err
        }
        
        // 2. 更新用户计数
        if err := tx.Model(&model.User{}).
            Where("id = ?", content.UserID).
            Update("content_count", gorm.Expr("content_count + 1")).Error; err != nil {
            return nil, err
        }
        
        return content, nil
    })
}
```

### 5.2 查询规范

```go
// Good: 使用 Where 链式调用
db.Where("status = ?", status).
   Where("created_at > ?", startTime).
   Order("created_at DESC").
   Limit(limit).
   Offset(offset).
   Find(&contents)

// Bad: 字符串拼接（SQL 注入风险）
db.Raw(fmt.Sprintf("SELECT * FROM contents WHERE status = '%s'", status))
```

### 5.3 软删除

```go
// 使用 gorm.DeletedAt
type Content struct {
    ID        uint           `gorm:"primaryKey"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
    // ...
}

// 查询时自动排除已删除
db.Find(&contents) // WHERE deleted_at IS NULL

// 包含已删除
db.Unscoped().Find(&contents)

// 永久删除
db.Unscoped().Delete(&content)
```

---

## 六、测试规范

### 6.1 测试文件命名

```
user_service.go      -> user_service_test.go
user_repository.go   -> user_repository_test.go
```

### 6.2 测试函数命名

```go
func TestUserService_Create(t *testing.T) {}
func TestUserService_Create_DuplicateUsername(t *testing.T) {}
func TestUserService_GetByID_NotFound(t *testing.T) {}
```

### 6.3 表格驱动测试

```go
func TestValidatePassword(t *testing.T) {
    tests := []struct {
        name     string
        password string
        wantErr  bool
    }{
        {"valid", "Password123!", false},
        {"too short", "abc", true},
        {"no number", "password!", true},
        {"no special", "Password123", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidatePassword(tt.password)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 6.4 覆盖率要求

| 阶段 | 目标覆盖率 |
|------|------------|
| MVP | 50% |
| 正式版 | 70% |
| 核心模块 | 80% |

```bash
# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 七、Git 规范

### 7.1 分支策略

```
main        -> 生产环境
develop     -> 开发环境
feature/*   -> 功能分支
hotfix/*    -> 紧急修复
```

### 7.2 Commit Message

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type:**
| 类型 | 说明 |
|------|------|
| feat | 新功能 |
| fix | Bug 修复 |
| docs | 文档更新 |
| style | 代码格式 |
| refactor | 重构 |
| test | 测试 |
| chore | 构建/工具 |

**示例:**
```
feat(auth): implement JWT refresh token

- Add refresh token storage with SHA256 hash
- Implement token rotation on refresh
- Revoke all tokens on logout

Closes #123
```

### 7.3 PR 规范

```markdown
## 变更说明
- 实现了什么功能
- 修复了什么问题

## 测试
- [ ] 单元测试已通过
- [ ] 本地测试已通过

## Checklist
- [ ] 代码符合开发规范
- [ ] 已添加必要注释
- [ ] 已更新相关文档
```

---

## 八、安全规范

### 8.1 输入校验

```go
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50,alphanum"`
    Password string `json:"password" binding:"required,min=8,max=72"`
    Nickname string `json:"nickname" binding:"max=100"`
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "参数错误: "+err.Error())
        return
    }
    // ...
}
```

### 8.2 敏感数据

```go
// 密码存储
hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// 密码校验
err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

// JSON 响应隐藏
type User struct {
    Password string `json:"-"` // 不序列化
}
```

### 8.3 SQL 注入防护

```go
// Good: 参数化查询
db.Where("username = ?", username).First(&user)

// Bad: 字符串拼接
db.Raw("SELECT * FROM users WHERE username = '" + username + "'")
```

---

## 九、代码审查清单

- [ ] 命名清晰，符合规范
- [ ] 错误处理完整，无吞错误
- [ ] 日志级别适当，无敏感信息
- [ ] 事务使用正确
- [ ] 无 SQL 注入风险
- [ ] 输入校验完整
- [ ] 测试覆盖关键逻辑
- [ ] 注释清晰，无冗余
