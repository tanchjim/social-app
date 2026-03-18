# 短视频图文社交平台完整架构设计文档

| 项目 | 内容 |
|------|------|
| 文档版本 | v1.0 |
| 最后更新 | 2026-03-17 |
| 适用范围 | 当前版本交付范围 |
| 维护者 | 开发团队 |

## 一、文档概述

本文档定义短视频图文社交平台当前版本的系统架构、技术选型、后端分层、核心业务流程、API 规范、数据库结构、部署基线、开发排期及风险控制要求。

### 1.1 功能范围

当前版本包含以下功能：

- 用户注册、登录、Token 刷新
- 用户资料查询与更新
- 图文内容与短视频内容发布
- 信息流列表与内容详情
- 评论发布与评论列表
- 点赞与取消点赞
- COS 直传
- 内容审核

当前版本不包含以下功能：

- 关注系统
- 私信系统
- 推荐算法
- 后台管理系统

## 二、总体架构

### 2.1 架构说明

系统采用 `Taro App + Go 后端 + 云服务` 的单体应用架构，内部按模块分层。后端以单服务单进程部署，对外统一提供 `Gin HTTP Server`。

| 层级 | 组件 | 职责 |
|------|------|------|
| 客户端 | Taro 3 + React | 页面渲染、交互处理、媒体上传、接口调用 |
| 接入层 | Gin HTTP Server | HTTP 路由分发、参数校验、认证鉴权、统一响应 |
| 业务层 | Service | 业务规则执行、审核调用、数据组装 |
| 数据访问层 | Repository | PostgreSQL 读写 |
| 数据层 | PostgreSQL、COS | 结构化数据、媒体文件存储 |
| 外部服务 | 内容审核服务 | 图文与视频内容审核 |

### 2.2 技术栈

| 类别 | 技术 |
|------|------|
| 前端 | Taro 3、React |
| 后端 | Go、Gin、GORM |
| 数据库 | PostgreSQL |
| 对象存储 | 腾讯云 COS |
| 日志 | Zap |
| 指标监控 | Prometheus |
| 鉴权 | JWT |

### 2.3 核心技术决策

| 主题 | 决策 |
|------|------|
| 内容审核 | 采用同步审核；发布请求在审核完成后返回结果 |
| 视频处理 | 客户端压缩后上传；视频文件大小上限为 100 MB，分辨率统一为 720p |
| 缓存策略 | MVP 阶段不使用缓存，所有查询直接访问 PostgreSQL |
| 日志与监控 | 后端输出 Zap 结构化日志，暴露 Prometheus 基础指标 |
| 媒体上传 | 客户端通过上传签名直传 COS，后端仅签发上传凭证与保存元数据 |
| 点赞一致性 | PostgreSQL 为数据真源；点赞关系写入事务表，内容计数同步更新 |

## 三、后端架构设计

### 3.1 分层约束

- `Handler` 负责参数解析、鉴权校验和响应封装，不直接访问数据库。
- `Service` 负责业务编排、审核调用和事务处理。
- `Repository` 负责 PostgreSQL 访问，不包含业务逻辑。
- 模块之间通过 `Service` 暴露的接口协作，不允许跨模块直接访问其他模块的 `Repository`。

### 3.2 目录结构

```text
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── router/
│   │   └── router.go
│   ├── handler/
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── content.go
│   │   ├── comment.go
│   │   ├── like.go
│   │   ├── upload.go
│   │   └── health.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   ├── model/
│   │   ├── user.go
│   │   ├── content.go
│   │   ├── comment.go
│   │   └── like.go
│   ├── repository/
│   │   ├── user.go
│   │   ├── content.go
│   │   ├── comment.go
│   │   └── like.go
│   ├── service/
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── content.go
│   │   ├── comment.go
│   │   ├── like.go
│   │   └── upload.go
│   └── integration/
│       ├── cos/
│       └── moderation/
├── pkg/
│   ├── jwt/
│   ├── logger/
│   ├── response/
│   └── utils/
├── go.mod
└── go.sum
```

### 3.3 模块职责

| 模块 | 职责 |
|------|------|
| Auth | 注册、登录、Token 刷新、JWT 签发 |
| User | 用户资料查询与更新、我的内容列表查询 |
| Content | 内容发布、信息流查询、详情查询、删除、审核结果查询 |
| Comment | 评论发布、评论列表查询 |
| Like | 点赞、取消点赞、点赞状态与计数维护 |
| Upload | COS 上传签名生成与上传约束校验 |

### 3.4 鉴权模型

```go
type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}
```

鉴权规则如下：

- 采用 JWT 认证，`Access Token` 无状态，`Refresh Token` 通过数据库管理。
- 访问受保护接口时，客户端在请求头中传递 `Authorization: Bearer <access_token>`。
- `Access Token` 用于接口访问，有效期 `24` 小时。
- `Refresh Token` 用于刷新登录态，有效期 `7` 天，存储于 `refresh_tokens` 表。
- 刷新接口返回新的 `Access Token` 与 `Refresh Token`，同时撤销旧 `Refresh Token`。
- `Refresh Token` 仅可使用一次，刷新后旧 `Token` 立即失效。
- 用户登出或修改密码时，撤销该用户所有 `Refresh Token`。

#### Token 安全设计

| 场景 | 处理方式 |
|------|----------|
| 正常刷新 | 生成新Token对，撤销旧Refresh Token |
| 并发刷新 | 使用数据库唯一约束，只有第一个请求成功 |
| 用户登出 | 撤销该用户所有Refresh Token |
| 修改密码 | 撤销该用户所有Refresh Token，强制重新登录 |
| Token泄露 | 用户可通过登出撤销所有Token |

#### Refresh Token 存储设计

`Refresh Token` 不存明文，仅存储 `SHA256` 哈希值：

```go
type RefreshToken struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"index;not null"`
    TokenHash string    `gorm:"uniqueIndex;size:64;not null"`  // SHA256(refresh_token)
    Revoked   bool      `gorm:"default:false"`
    CreatedAt time.Time
    ExpiresAt time.Time `gorm:"index"`
}
```

刷新流程：

1. 验证 `Refresh Token` 签名和过期时间
2. 计算 `Token` 哈希，查询 `refresh_tokens` 表
3. 检查 `Revoked` 状态，已撤销则返回错误
4. 生成新 `Token` 对
5. 撤销旧 `Token`，存入新 `Token`
6. 返回新 `Token` 对

### 3.5 Repository 接口规范

```go
type ContentRepository interface {
    Create(ctx context.Context, content *model.Content) error
    GetByID(ctx context.Context, id uint) (*model.Content, error)
    List(ctx context.Context, status string, offset, limit int) ([]*model.Content, error)
    ListByUser(ctx context.Context, userID uint, status string, offset, limit int) ([]*model.Content, error)
    Update(ctx context.Context, content *model.Content) error
    SoftDelete(ctx context.Context, id uint, deletedAt time.Time) error
}
```

### 3.6 Service 职责划分

| 方法 | 职责 |
|------|------|
| `CreateContent` | 参数校验、草稿入库、审核调用、状态更新 |
| `GetFeed` | 信息流查询、作者信息组装 |
| `ListMyContents` | 按状态筛选当前用户内容列表 |
| `GetReviewResult` | 读取审核结果与拒绝原因 |
| `LikeContent` | 幂等控制、点赞关系写入、计数更新 |
| `CreateComment` | 内容存在性校验、评论入库、评论计数更新 |

## 四、核心业务流程

### 4.1 用户注册与登录流程

1. 客户端调用注册接口提交用户名、密码、昵称。
2. 后端校验参数并写入用户表，密码使用 `bcrypt` 加密后存储。
3. 客户端调用登录接口。
4. 后端校验账号密码，签发 `Access Token` 与 `Refresh Token`。

### 4.2 媒体上传与内容发布流程

1. 客户端调用上传签名接口获取 COS 上传凭证。
2. 客户端将图片或视频文件直传至 COS。
3. 客户端调用内容发布接口提交内容元数据与媒体地址。
4. 后端校验媒体地址、内容类型、标题和描述，并写入 `contents` 表，初始状态为 `draft`。
5. 后端将内容状态更新为 `reviewing`，同步调用内容审核服务。
6. 审核通过后更新内容状态为 `published`。
7. 审核不通过时更新内容状态为 `rejected`，写入 `reject_reason`，并返回错误码 `20003`。

### 4.3 信息流读取流程

1. 客户端调用信息流接口并传入分页参数。
2. 后端直接读取 PostgreSQL，并按发布时间倒序分页查询。
3. 信息流仅返回状态为 `published` 的内容。
4. 返回内容列表、作者信息、点赞状态和分页信息。

### 4.4 点赞流程

1. 客户端调用点赞接口并传入动作类型。
2. 后端根据 `content_id + user_id` 唯一约束判断当前状态。
3. 点赞动作执行插入并递增 `like_count`。
4. 取消点赞动作执行删除并递减 `like_count`。
5. 点赞结果写库成功后返回最新计数。

### 4.5 内容状态流转

- `draft`：内容草稿，记录已创建，尚未进入审核完成态。
- `reviewing`：内容已提交审核，等待审核服务返回结果。
- `published`：审核通过，可在信息流和详情页中对外展示。
- `rejected`：审核拒绝，仅作者可查看审核结果与拒绝原因。
- `deleted`：作者主动删除后的软删除状态，不再对外展示。
- 正常发布链路为 `draft -> reviewing -> published`。
- 审核拒绝链路为 `draft -> reviewing -> rejected`。
- 拒绝后重新提交审核时，状态可从 `rejected -> reviewing`。
- 删除操作采用软删除，状态可从 `draft`、`rejected` 或 `published` 流转为 `deleted`，并记录 `deleted_at`。

## 五、API 设计

### 5.1 通用规范

#### 基础路径

所有业务接口统一使用前缀：

```text
/api/v1
```

#### 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

错误响应格式：

```json
{
  "code": 30002,
  "message": "server internal error",
  "data": null
}
```

#### 分页规则

- `page` 默认值为 `1`。
- `page_size` 默认值为 `20`，最大值为 `50`。

#### 时间格式

所有时间字段统一使用 `RFC 3339` 格式，例如 `2026-03-17T10:00:00Z`。

### 5.2 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 10001 | 用户名已存在 |
| 10002 | 密码错误 |
| 10003 | 用户不存在 |
| 10004 | Token 无效 |
| 10005 | Token 已过期 |
| 10006 | Refresh Token 已撤销 |
| 10007 | Refresh Token 不存在 |
| 20001 | 内容不存在 |
| 20002 | 无权操作 |
| 20003 | 内容审核未通过 |
| 30001 | 参数错误 |
| 30002 | 服务器内部错误 |

### 5.3 接口总览

| 模块 | 方法 | 路径 | 说明 |
|------|------|------|------|
| Auth | POST | `/api/v1/auth/register` | 用户注册 |
| Auth | POST | `/api/v1/auth/login` | 用户登录 |
| Auth | POST | `/api/v1/auth/refresh` | 刷新 Token |
| Auth | POST | `/api/v1/auth/logout` | 用户登出 |
| Upload | POST | `/api/v1/uploads/sign` | 获取 COS 上传签名 |
| User | GET | `/api/v1/users/:id` | 获取用户信息 |
| User | PUT | `/api/v1/users/:id` | 更新用户信息 |
| User | GET | `/api/v1/users/me/contents` | 获取我的内容列表 |
| Content | POST | `/api/v1/contents` | 发布内容 |
| Content | GET | `/api/v1/contents` | 获取信息流 |
| Content | GET | `/api/v1/contents/:id` | 获取内容详情 |
| Content | GET | `/api/v1/contents/:id/review-result` | 获取审核结果 |
| Content | DELETE | `/api/v1/contents/:id` | 删除内容 |
| Comment | GET | `/api/v1/contents/:id/comments` | 获取评论列表 |
| Comment | POST | `/api/v1/contents/:id/comments` | 发表评论 |
| Like | POST | `/api/v1/contents/:id/like` | 点赞或取消点赞 |
| Health | GET | `/health` | 健康检查 |

### 5.4 Auth 模块

#### POST `/api/v1/auth/register`

Request:

```json
{
  "username": "testuser",
  "password": "password123",
  "nickname": "测试用户"
}
```

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1,
    "username": "testuser",
    "nickname": "测试用户"
  }
}
```

#### POST `/api/v1/auth/login`

Request:

```json
{
  "username": "testuser",
  "password": "password123"
}
```

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1,
    "username": "testuser",
    "nickname": "测试用户",
    "avatar": "https://cos.xxx.com/avatar/1.jpg",
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

#### POST `/api/v1/auth/refresh`

Request:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

#### POST `/api/v1/auth/logout`

说明：撤销当前用户所有 `Refresh Token`，用户需要重新登录。客户端应在调用后清除本地存储的 `Token`。

Request:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

### 5.5 Upload 模块

#### POST `/api/v1/uploads/sign`

Request:

```json
{
  "file_name": "video-001.mp4",
  "content_type": "video/mp4",
  "file_size": 52428800
}
```

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "upload_url": "https://cos.xxx.com/upload/video-001.mp4",
    "object_key": "videos/2026/03/17/video-001.mp4",
    "expired_at": "2026-03-17T10:30:00Z"
  }
}
```

### 5.6 User 模块

#### GET `/api/v1/users/:id`

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1,
    "username": "testuser",
    "nickname": "测试用户",
    "avatar": "https://cos.xxx.com/avatar/1.jpg",
    "bio": "这是我的简介",
    "content_count": 10,
    "like_count": 100,
    "created_at": "2026-03-17T10:00:00Z"
  }
}
```

#### PUT `/api/v1/users/:id`

Request:

```json
{
  "nickname": "新昵称",
  "bio": "新的简介",
  "avatar": "https://cos.xxx.com/avatar/new.jpg"
}
```

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1,
    "nickname": "新昵称",
    "bio": "新的简介",
    "avatar": "https://cos.xxx.com/avatar/new.jpg"
  }
}
```

#### GET `/api/v1/users/me/contents`

Query Parameters:

- `status`：可选，支持 `draft`、`reviewing`、`published`、`rejected`、`deleted`
- `page`
- `page_size`

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "content_id": 1,
        "type": "video",
        "title": "我的第一个视频",
        "status": "rejected",
        "cover_url": "https://cos.xxx.com/covers/xxx.jpg",
        "like_count": 0,
        "comment_count": 0,
        "reject_reason": "视频封面包含违规元素",
        "created_at": "2026-03-17T10:00:00Z",
        "updated_at": "2026-03-17T10:05:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 20
  }
}
```

### 5.7 Content 模块

#### POST `/api/v1/contents`

Request:

```json
{
  "type": "video",
  "title": "我的第一个视频",
  "description": "视频描述",
  "media_url": "https://cos.xxx.com/videos/xxx.mp4",
  "cover_url": "https://cos.xxx.com/covers/xxx.jpg",
  "tags": ["生活", "日常"]
}
```

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "content_id": 1,
    "status": "published",
    "type": "video",
    "title": "我的第一个视频",
    "media_url": "https://cos.xxx.com/videos/xxx.mp4",
    "cover_url": "https://cos.xxx.com/covers/xxx.jpg",
    "created_at": "2026-03-17T10:00:00Z"
  }
}
```

#### GET `/api/v1/contents`

Query Parameters:

- `page`
- `page_size`

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "content_id": 1,
        "type": "video",
        "title": "视频标题",
        "description": "视频描述",
        "cover_url": "https://cos.xxx.com/covers/xxx.jpg",
        "media_url": "https://cos.xxx.com/videos/xxx.mp4",
        "author": {
          "user_id": 1,
          "nickname": "作者昵称",
          "avatar": "https://cos.xxx.com/avatar/1.jpg"
        },
        "like_count": 100,
        "comment_count": 20,
        "is_liked": false,
        "created_at": "2026-03-17T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

#### GET `/api/v1/contents/:id`

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "content_id": 1,
    "type": "video",
    "title": "视频标题",
    "description": "视频描述",
    "status": "published",
    "cover_url": "https://cos.xxx.com/covers/xxx.jpg",
    "media_url": "https://cos.xxx.com/videos/xxx.mp4",
    "author": {
      "user_id": 1,
      "nickname": "作者昵称",
      "avatar": "https://cos.xxx.com/avatar/1.jpg"
    },
    "like_count": 100,
    "comment_count": 20,
    "is_liked": false,
    "tags": ["生活", "日常"],
    "created_at": "2026-03-17T10:00:00Z"
  }
}
```

#### GET `/api/v1/contents/:id/review-result`

说明：该接口仅内容作者可访问，用于查看审核状态与拒绝原因。

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "content_id": 1,
    "status": "rejected",
    "reject_reason": "视频封面包含违规元素",
    "reviewed_at": "2026-03-17T10:05:00Z"
  }
}
```

#### DELETE `/api/v1/contents/:id`

说明：删除操作为软删除，后端将内容状态更新为 `deleted` 并写入 `deleted_at`。

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": null
}
```

### 5.8 Comment 模块

#### GET `/api/v1/contents/:id/comments`

Query Parameters:

- `page`
- `page_size`

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "comment_id": 1,
        "content": "评论内容",
        "author": {
          "user_id": 2,
          "nickname": "评论者昵称",
          "avatar": "https://cos.xxx.com/avatar/2.jpg"
        },
        "created_at": "2026-03-17T10:00:00Z"
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

#### POST `/api/v1/contents/:id/comments`

Request:

```json
{
  "content": "这是一条评论"
}
```

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "comment_id": 1,
    "content": "这是一条评论",
    "author": {
      "user_id": 1,
      "nickname": "我的昵称",
      "avatar": "https://cos.xxx.com/avatar/1.jpg"
    },
    "created_at": "2026-03-17T10:00:00Z"
  }
}
```

### 5.9 Like 模块

#### POST `/api/v1/contents/:id/like`

Request:

```json
{
  "action": "like"
}
```

Response:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "is_liked": true,
    "like_count": 101
  }
}
```

## 六、数据库设计

### 6.1 数据模型

| 表名 | 说明 |
|------|------|
| `users` | 用户基础信息 |
| `contents` | 内容主表 |
| `comments` | 评论表 |
| `likes` | 点赞关系表 |
| `refresh_tokens` | Refresh Token 存储（支持撤销） |

### 6.2 PostgreSQL 初始化脚本

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(100),
    avatar VARCHAR(500),
    bio TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);

CREATE TABLE contents (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    type VARCHAR(20) NOT NULL CHECK (type IN ('image', 'video')),
    title VARCHAR(200),
    description TEXT,
    media_url VARCHAR(500) NOT NULL,
    cover_url VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'reviewing', 'published', 'rejected', 'deleted')),
    reject_reason TEXT,
    like_count INTEGER NOT NULL DEFAULT 0,
    comment_count INTEGER NOT NULL DEFAULT 0,
    tags TEXT[] NOT NULL DEFAULT '{}',
    deleted_at TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_contents_user_id ON contents(user_id);
CREATE INDEX idx_contents_status ON contents(status);
CREATE INDEX idx_contents_created_at ON contents(created_at DESC);

CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    content_id BIGINT NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_comments_content_id ON comments(content_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);

CREATE TABLE likes (
    id BIGSERIAL PRIMARY KEY,
    content_id BIGINT NOT NULL REFERENCES contents(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (content_id, user_id)
);

CREATE INDEX idx_likes_content_id ON likes(content_id);
CREATE INDEX idx_likes_user_id ON likes(user_id);

CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens(expires_at);
```

### 6.3 初始化测试数据

```sql
INSERT INTO users (username, password, nickname) VALUES
('test1', '$2a$10$...', '测试用户1'),
('test2', '$2a$10$...', '测试用户2');
```

说明：`password` 字段存储 `bcrypt` 密文。

## 七、配置与部署

### 7.1 环境变量

```env
# Server
SERVER_PORT=8080
GIN_MODE=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=social_app

# JWT
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRE_HOURS=24
JWT_REFRESH_EXPIRE_DAYS=7

# Tencent COS
COS_SECRET_ID=your_cos_secret_id
COS_SECRET_KEY=your_cos_secret_key
COS_BUCKET=your_bucket_name
COS_REGION=ap-guangzhou

# Content Moderation
MODERATION_API_KEY=your_moderation_api_key
MODERATION_API_URL=https://api.example.com/moderate
```

### 7.2 部署基线

后端部署要求：

- 采用单体部署，单服务单进程运行，`Gin HTTP Server` 为唯一对外入口
- PostgreSQL 已创建并完成初始化
- 环境变量已完整配置
- Go 依赖已安装
- 服务启动命令为 `go run cmd/server/main.go`
- 健康检查接口 `GET /health` 可访问

前端部署要求：

- API 基础地址已配置
- 依赖已安装
- 小程序或 App 构建产物已生成
- COS 上传配置已完成

COS 配置要求：

- 存储桶已创建
- 跨域规则已配置
- 上传目录权限已配置

说明：本节为开发、测试环境的最小部署基线；生产环境的灾备与高可用要求以第八章为准。

## 八、灾备与高可用设计

### 8.1 设计目标

灾备方案按“核心读链路优先、数据安全优先、发布链路可降级”设计，目标如下：

| 目标项 | 目标值 | 说明 |
|--------|--------|------|
| 核心读链路可用性 | `>= 99.9%` | 包含登录、信息流、详情、评论列表 |
| 发布链路可用性 | `>= 99.5%` | 审核服务异常时允许降级，不要求强一致实时发布 |
| 单机故障影响范围 | 单实例级别 | 任一应用实例故障不应导致整站不可用 |
| 数据恢复能力 | 支持误删恢复与时间点恢复 | 依赖 PostgreSQL 备份 + WAL 连续归档 |

### 8.2 单点故障分析

当前文档中的 MVP 架构存在明显单点，生产环境需要按下表消除：

| 组件 | 当前单点风险 | 故障影响 | 生产方案 |
|------|--------------|----------|----------|
| `Gin HTTP Server` 单实例 | 进程退出、机器宕机即整站不可用 | 登录、发布、查询全部失败 | 至少 `2` 个无状态实例，跨可用区部署，前置 `SLB/Ingress` 健康检查 |
| PostgreSQL 单实例 | 数据库不可连即核心业务不可用 | 登录、发内容、评论、点赞、信息流全部失败 | 托管 PostgreSQL 主备架构，同城双可用区，自动故障切换 |
| 审核服务同步调用 | 外部审核不可用会阻塞发布 | 新内容无法发布，已有内容浏览不受影响 | 审核超时熔断，写库后转异步补审，状态保留为 `reviewing` |
| COS 单地域桶 | 上传或媒体访问依赖单地域存储 | 新内容无法上传；历史内容可能无法播放或封面加载失败 | 开启对象版本控制，关键桶启用跨地域复制，客户端支持占位图和重试 |
| 负载均衡 / DNS 单入口 | 接入层故障会导致整站不可达 | 所有请求失败 | 采用云负载均衡高可用实例；跨地域场景配合 DNS 故障切换 |
| 配置与密钥仅单份保存 | 误删后无法快速重建实例 | 无法扩容、无法恢复服务 | 配置中心或密钥管理服务托管，保留加密备份 |

### 8.3 生产推荐拓扑

推荐采用“同城双可用区高可用 + 异地灾备”的分层部署：

- 接入层：`SLB/Ingress` 部署在主地域，对后端实例做存活与就绪检查。
- 应用层：`2~3` 个 Go API 实例，跨两个可用区部署，实例保持无状态，支持滚动发布。
- 数据层：PostgreSQL 主库在 `AZ-A`，同步或准同步备库在 `AZ-B`，异地保留一个异步灾备实例或冷备恢复实例。
- 对象存储：COS 开启版本控制；核心媒体桶按成本选择跨地域复制到灾备地域。
- 外部审核：调用链增加超时、重试、熔断与降级开关；审核结果落库，避免只存在外部系统。
- 监控告警：Prometheus 采集应用、数据库、审核调用与备份任务指标；关键告警同步到值班通道。

### 8.4 数据库备份策略

PostgreSQL 采用“全量备份 + WAL 增量归档 + 周期逻辑备份”的组合方案：

| 备份类型 | 方式 | 频率 | 保留周期 | 用途 |
|----------|------|------|----------|------|
| 全量物理备份 | 数据库快照或 `pg_basebackup` | 每日 `02:00` | `14` 天 | 整库恢复、构建新实例 |
| 增量备份 | WAL 连续归档 | 实时归档，最多 `5` 分钟一批 | `14` 天 | 时间点恢复（PITR） |
| 逻辑备份 | `pg_dump` 导出结构和关键表 | 每周一次 | `8` 周 | 误删单表、跨环境校验、审计留档 |
| 月度归档 | 每月首日保留一份全量 | 每月一次 | `6` 个月 | 长周期回溯、重大事故兜底 |

备份要求：

- 备份文件必须存放到独立存储，不与生产数据库同机或同磁盘。
- 备份文件开启加密、校验和完整性校验，防止恢复时发现文件损坏。
- 每月至少执行一次恢复演练，验证“全量备份 + WAL”可以成功恢复。
- 备份任务失败、WAL 归档延迟、磁盘空间不足必须有告警。

### 8.5 数据恢复方案

#### 8.5.1 误删数据恢复

| 场景 | 恢复方式 | 目标时长 |
|------|----------|----------|
| 内容误删 | 由于 `contents` 已采用软删除，优先通过将状态从 `deleted` 恢复为原状态处理 | `5~15` 分钟 |
| 评论 / 点赞误删 | 从最近时间点恢复到临时库，导出目标时间点前的数据，再回放到生产库 | `30~60` 分钟 |
| 用户资料误改或误删 | 通过 PITR 恢复到临时库，对比后按主键回写 | `30~60` 分钟 |

恢复流程：

1. 先冻结相关写操作，防止二次覆盖。
2. 确认误操作发生时间，选择最近的恢复时间点。
3. 在临时实例执行 PITR，不直接覆盖生产库。
4. 对恢复出的数据做主键、时间范围和业务状态校验。
5. 通过脚本定向回写生产库，并记录审计日志。

#### 8.5.2 数据损坏恢复

| 场景 | 处置策略 |
|------|----------|
| 主库机器或磁盘损坏 | 立即切换到同城备库，应用连接串自动指向新主库 |
| 逻辑错误已复制到备库 | 使用最近全量备份 + WAL 恢复到故障前时间点，拉起新实例后切流 |
| 单表或索引损坏 | 优先在临时实例验证修复；无法在线修复时切换到备库，再离线修复原主库 |

### 8.6 服务降级策略

审核服务不可用时，目标是“保住读链路，限制发布链路，不放过未审核内容”。

| 故障场景 | 降级策略 | 用户侧表现 |
|----------|----------|------------|
| 审核服务超时或返回 `5xx` | 内容元数据和媒体地址正常入库，状态保持 `reviewing`，由后台重试任务补审 | 用户收到“内容已提交，审核排队中” |
| 审核服务持续不可用超过阈值（如 `10` 分钟） | 打开降级开关，暂停同步审核；仅允许保存草稿或提交待审，不进入 `published` | 发布成功率下降，但浏览、评论、点赞继续可用 |
| 审核队列积压严重 | 优先限制视频发布，仅保留图文待审入口，避免大文件占满审核资源 | 用户可继续保存内容，但发布时间延后 |

实施要求：

- 审核调用超时建议控制在 `2~3` 秒，失败最多重试 `1` 次。
- 发布接口在降级模式下返回明确状态，不返回“已发布”假象。
- 运营侧需有“待审积压”看板，便于人工补审或清理异常任务。

### 8.7 多机房 / 多地域部署考虑

现阶段不建议直接做双地域双写，原因是评论、点赞、内容状态流转都依赖单写数据库事务，跨地域双写会显著增加冲突处理和运维复杂度。推荐分阶段建设：

| 阶段 | 部署方式 | 适用阶段 |
|------|----------|----------|
| Phase 1 | 同城双可用区：应用多实例 + PostgreSQL 主备 | 上线初期、MVP 到小规模生产 |
| Phase 2 | 异地灾备：异步数据库副本 + COS 跨地域复制 + DNS 切换 | 日活增长后，需要抗地域级故障 |
| Phase 3 | 热备增强：异地应用实例常驻，定期演练切流 | 对可用性要求进一步提升时 |

多地域切换原则：

- 正常情况下只允许单地域写入，避免双写冲突。
- 异地机房默认不承接生产流量，处于热备或温备状态。
- 发生地域级故障时，通过 DNS 或网关切换到灾备地域，数据库提升为主库后开放写流量。
- 切换完成后需校验用户、内容、评论、点赞四类核心表数据一致性。

### 8.8 灾难恢复目标（RTO / RPO）

| 故障级别 | 典型场景 | RTO | RPO |
|----------|----------|-----|-----|
| 实例级故障 | 单个 API 实例宕机 | `< 5` 分钟 | `0` |
| 可用区级数据库故障 | 主库宕机、存储损坏 | `< 15` 分钟 | `0~5` 分钟 |
| 审核服务故障 | 第三方审核不可用 | `< 3` 分钟进入降级 | `0`，不丢内容元数据 |
| 地域级故障 | 主地域整体不可用 | `< 60` 分钟 | `<= 15` 分钟 |
| 人为误删 / 误操作 | 删除内容、评论或表数据 | `< 2` 小时 | `<= 5` 分钟 |

说明：

- 若同城主备采用同步复制，数据库主备切换场景可做到 `RPO = 0`。
- 若异地灾备采用异步复制，地域级故障的 `RPO` 取决于复制延迟和 WAL 归档频率，建议按 `15` 分钟内设计。

### 8.9 MVP 阶段最小灾备方案

MVP 阶段不追求多地域全量建设，但以下能力不能省略：

| 必选项 | 最低要求 |
|--------|----------|
| 应用高可用 | `2` 个 API 实例，部署在不同可用区或不同宿主机，前置负载均衡 |
| 数据库高可用 | 使用云 PostgreSQL 主备版，开启自动故障切换 |
| 数据备份 | 开启每日全量备份、WAL 连续归档、`7~14` 天 PITR |
| 审核降级 | 审核失败时内容保留为 `reviewing`，已有内容浏览不受影响 |
| 运维值守 | 建立数据库连接异常、实例存活、备份失败三类基础告警 |
| 演练机制 | 上线前至少做一次主备切换演练和一次误删恢复演练 |

可延后项：

- 异地数据库热备
- COS 跨地域复制
- 跨地域 DNS 自动切换

对当前项目而言，MVP 的最低可行结论是：不能继续保持“单应用实例 + 单数据库实例”的真实生产部署；至少要做到“应用双实例 + 数据库主备 + 备份可恢复 + 审核可降级”。

## 九、开发排期与协作机制

### 9.1 两周开发排期

#### Week 1：基础架构、用户模块、内容模块

| 天数 | 后端任务 | 前端任务 |
|------|----------|----------|
| Day 1 | 项目骨架搭建、数据库初始化、配置管理 | Taro 项目初始化、路由配置 |
| Day 2 | JWT 认证中间件、用户模型与仓储实现 | 登录与注册页面 UI |
| Day 3 | Auth Service 与 Handler、注册登录接口 | 登录与注册联调 |
| Day 4 | User Service 与 Handler、用户信息接口 | 个人中心页面 |
| Day 5 | Content 模型与仓储实现 | 信息流页面 UI |
| Day 6 | Content Service 与 Handler、发布与列表接口 | 发布内容页面 |
| Day 7 | 内容详情接口、我的内容列表接口 | 内容详情页 |

#### Week 2：互动模块、COS 集成、联调测试

| 天数 | 后端任务 | 前端任务 |
|------|----------|----------|
| Day 1 | Comment 模型、仓储、服务实现 | 评论列表 UI |
| Day 2 | Comment Handler 与评论接口 | 评论功能联调 |
| Day 3 | Like 模型、仓储、服务实现 | 点赞状态与动画 |
| Day 4 | COS 签名接口、上传接口联调 | 图片与视频上传组件 |
| Day 5 | 内容审核服务接入与状态流转实现 | 上传进度与失败提示 |
| Day 6 | 审核结果接口、软删除逻辑 | 审核结果页与错误处理 |
| Day 7 | 联调测试与缺陷修复 | 联调测试与缺陷修复 |

### 9.2 协作机制

- 后端优先输出 API 文档与错误码定义。
- 前端基于 Mock 数据并行开发，联调阶段切换真实接口。
- 每日固定安排 1 小时联调。
- 需求、缺陷和联调问题统一在飞书群同步。

## 十、风险控制

| 风险点 | 影响 | 控制措施 |
|--------|------|----------|
| COS 配额超限 | 用户无法上传 | 按用户维度限制上传频率与单日流量，并配置监控告警 |
| 内容审核超时 | 发布链路阻塞 | 审核请求超时后进入降级策略，内容保留在 `reviewing` 并通过后台任务补审，不影响已发布内容浏览 |
| 点赞并发冲突 | 点赞计数不准确 | 使用唯一约束与数据库事务控制点赞关系和计数一致性 |
| Token 泄露 | 账号安全风险 | Access Token 短时有效，Refresh Token 设置合理有效期并要求客户端安全存储 |
| 数据库连接池耗尽 | 服务不可用 | 限制连接池大小并监控活跃连接数与慢查询 |
