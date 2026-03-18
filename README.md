# Social App Backend

短视频图文社交平台后端服务

## 技术栈

- Go 1.21+
- Gin Web Framework
- GORM ORM
- PostgreSQL
- JWT Authentication
- Zap Logger

## 项目结构

```
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   ├── router/
│   ├── service/
│   └── integration/
│       ├── cos/
│       └── moderation/
├── pkg/
│   ├── jwt/
│   ├── logger/
│   ├── response/
│   └── utils/
├── go.mod
└── README.md
```

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/yourorg/social-app.git
cd social-app
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，填入实际配置
```

### 3. 安装依赖

```bash
go mod tidy
```

### 4. 运行服务

```bash
go run cmd/server/main.go
```

### 5. 健康检查

```bash
curl http://localhost:8080/health
```

## API 文档

参见 [ARCHITECTURE.md](./ARCHITECTURE.md)

## 开发进度

- [x] 项目骨架
- [x] 路由注册
- [x] 中间件
- [x] 数据模型
- [ ] 业务逻辑实现
- [ ] 单元测试
- [ ] 集成测试
