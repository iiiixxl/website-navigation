# Site.Navigation.Api

基于 `monitor.api` 结构搭建的 Go Gin 后端脚手架，用于 Site.Navigation 相关接口扩展。

## 技术栈

- Go 1.22+
- Gin
- GORM + MySQL
- JWT 鉴权（支持 AllowAnonymous）
- 本地 `config/config.json`（标准库解析，无 Consul / Viper）
- Redis（可选）

## 目录结构

```
Site.Navigation.Api/
├── main.go
├── go.mod
├── config/
│   ├── config.go              # 读取本地 config.json
│   ├── config.example.json    # 配置示例
│   ├── config.json            # 本地配置（已 gitignore）
│   ├── database.go
│   └── redis.go
├── controller/
│   ├── health_controller.go
│   ├── auth_controller.go
│   ├── user_controller.go
│   └── prompt_controller.go
├── service/
│   ├── user_service.go
│   └── prompt_service.go
├── model/
│   ├── user.go
│   └── prompt.go
├── middleware/
│   ├── auth.go
│   ├── allow_anonymous.go
│   └── route_whitelist.go
├── router/
│   ├── router.go
│   └── route_helper.go
├── utils/
│   ├── jwt.go
│   └── redis_helper.go
└── scripts/
    ├── 02_t_user.sql
    └── 03_t_prompt.sql
```

## 快速开始

### 1. 安装依赖

```bash
cd Site.Navigation.Api
go mod tidy
```

### 2. 配置

复制并修改配置：

```bash
copy config\config.example.json config\config.json
```

重点修改：

- `database.dsn`：MySQL 连接串
- `jwt.secret`：JWT 密钥
- `redis.enabled`：本地可先设为 `false`
- `server.port`：默认 `18080`

### 3. 建表

用户表、AI 模板表分别执行：

```bash
# 先手动选中库 site.navigation，再执行
# scripts/02_t_user.sql
# scripts/03_t_prompt.sql
```

用户表也可单独参考 `scripts/02_t_user.sql`。AI 模板含初始数据见 `scripts/03_t_prompt.sql`。

### 4. 启动

```bash
go run main.go
```

健康检查：`GET http://localhost:18080/health`

## 配置说明

只读 `config/config.json`，把连接串等写死在文件里即可，例如：

```json
{
  "database": {
    "dsn": "user:pwd@tcp(host:3306)/db?charset=utf8mb4&parseTime=True&loc=Local"
  },
  "jwt": { "secret": "your-secret", "expire_hours": 24 },
  "redis": { "enabled": false },
  "server": { "port": "18080" }
}
```

## 主要接口

### 用户 / 认证

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/health` | 否 | 健康检查 |
| POST | `/api/v1/auth/login` | 否 | 登录拿 token |
| POST | `/api/v1/users/CreateUser` | 否 | 创建用户 |
| GET | `/api/v1/users/GetUserList` | 是 | 用户列表 |
| GET | `/api/v1/users/GetUserById?id=1` | 是 | 用户详情 |
| POST | `/api/v1/users/UpdateUser` | 是 | 更新用户 |
| POST | `/api/v1/users/DeleteUser` | 是 | 软删除用户 |

### AI 提问模板

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| GET | `/api/v1/prompts/GetTree` | 否 | 分类+模板树（对齐前端） |
| GET | `/api/v1/prompts/GetTree?keyword=审查` | 否 | 按名称/正文搜索 |
| GET | `/api/v1/prompts/GetItemById?id=1` | 否 | 单个模板 |
| POST | `/api/v1/prompts/CreateCategory` | 是 | 创建分类 |
| POST | `/api/v1/prompts/UpdateCategory` | 是 | 更新分类 |
| POST | `/api/v1/prompts/DeleteCategory` | 是 | 删除分类（模板一并软删） |
| POST | `/api/v1/prompts/CreateItem` | 是 | 创建模板 |
| POST | `/api/v1/prompts/UpdateItem` | 是 | 更新模板正文 |
| POST | `/api/v1/prompts/DeleteItem` | 是 | 删除模板 |

`GetTree` 返回结构示例（与静态页 `PROMPT_DATA` 同形）：

```json
{
  "categories": [
    {
      "id": 1,
      "title": "常用模板",
      "sort_order": 1,
      "items": [
        {
          "id": 1,
          "category_id": 1,
          "name": "开发",
          "content": "请根据需求实现代码...",
          "sort_order": 1
        }
      ]
    }
  ]
}
```

正文使用 `LONGTEXT` 存储，**换行原样保留**，不会合并。

鉴权头：`Authorization: Bearer <token>`

## 与 monitor.api 的对应关系

- 分层：`controller / service / model / middleware / router / config / utils`
- 路由：`RegisterRoute` + `AllowAnonymous`
- JWT：中间件全局挂载，白名单跳过
- 配置：仅本地 JSON，不做 Consul
