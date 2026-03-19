# AGENTS.md - Social App Backend

短视频图文社交平台后端服务。

## 技术栈

- Go 1.26 + Gin + GORM + PostgreSQL
- 腾讯云 COS（对象存储）
- JWT 认证

## 项目结构

```
social-app/
├── cmd/server/main.go          # 入口
├── internal/
│   ├── config/                 # 配置管理
│   ├── handler/                # HTTP 处理器
│   ├── middleware/             # 中间件
│   ├── model/                  # 数据模型
│   ├── repository/             # 数据访问层
│   ├── router/                 # 路由注册
│   ├── service/                # 业务逻辑层
│   └── integration/            # 外部服务集成
├── pkg/                        # 公共工具包
│   ├── jwt/
│   ├── logger/
│   ├── response/
│   └── utils/
├── DEVELOPMENT.md              # 开发规范（详细）
└── ARCHITECTURE.md             # 架构设计文档
```

## 开发规范摘要

### 代码风格
- 命名：大驼峰（导出）/ 小驼峰（私有）
- 接口：动词+er 或名词，如 `UserRepository`
- 错误变量：Err 前缀，如 `ErrUserNotFound`
- 导入顺序：标准库 → 第三方 → 内部包

### 错误处理
```go
// 使用 errors.Is/As 传递
// 添加上下文
return nil, fmt.Errorf("failed to get user: %w", err)

// 定义业务错误
var ErrUserNotFound = errors.New("user not found")
```

### 错误码
| 范围 | 模块 |
|------|------|
| 10000-10999 | Auth |
| 20000-20999 | User |
| 30000-30999 | Content |
| 40000-40999 | Comment |
| 50000-50999 | Like |

### API 设计
- RESTful：GET 查询, POST 创建, PUT 更新, DELETE 删除
- 统一响应：`{code, message, data}`
- 分页：`page`(1+), `page_size`(1-50, 默认20)
- 版本：`/api/v1/...`

### 日志规范
- 结构化：`logger.Info("msg", zap.Uint("user_id", id))`
- 禁止记录：密码、Token、身份证号

### 数据库
- 事务：跨表操作必须用事务
- 查询：参数化 `db.Where("id = ?", id)`，禁止拼接
- 删除：软删除 `gorm.DeletedAt`

### 测试
- 文件：`xxx_test.go`
- 函数：`TestXxx`
- 覆盖率：MVP 50%, 正式版 70%

### Git Commit
```
feat(auth): implement JWT refresh token
fix(content): resolve pagination bug
docs: update API documentation
```

## 常用命令

```bash
# 运行服务
go run cmd/server/main.go

# 编译
go build -o bin/server ./cmd/server

# 测试
go test ./...

# 测试覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 代码检查
golangci-lint run ./...
```

## 详细规范

详见 [DEVELOPMENT.md](./DEVELOPMENT.md)
