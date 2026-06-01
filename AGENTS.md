# AGENTS.md

本文件用于约定后续在本项目中协作的 AI Agent / Codex 开发规范。

## 项目概况

- 项目类型：信息管理平台 BS 项目
- 前端：Vue 3 + Composition API + axios + Element Plus + Vite
- 后端：Go + Gin + GORM
- 数据库：MySQL 8.0
- 数据库名称：`medical_info`

## 重要约定

- 用户已明确说明：后续修改默认不要修改 `README.md`，除非用户再次明确要求。
- 前端开发端口：`8080`
- 后端服务端口：`8888`
- 后端 API 返回格式需要包含：
  - `code`
  - `message`
  - `errorMessage`
  - `data`
- 成功响应的 `code` 固定为 `0`。

## 常用命令

后端验证：

```bash
cd backend
go test ./...
```

前端构建验证：

```bash
cd frontend
npm.cmd run build
```

前端开发启动：

```bash
cd frontend
npm.cmd run dev
```

## 代码修改规则

- 后端新增接口时，优先复用 `backend/handlers/response.go` 中的统一响应方法。
- 后端模型放在 `backend/models`。
- 后端接口逻辑放在 `backend/handlers`。
- 路由统一在 `backend/main.go` 注册。
- Go 代码中双引号内的字符串命名使用英文下划线 `_` 分隔，不要把多个词直接连在一起。
- 前端主要页面逻辑目前集中在 `frontend/src/App.vue`。
- 前端全局样式在 `frontend/src/styles.css`。
- 不要提交或依赖 `frontend/dist`、`frontend/node_modules`、`.gocache`、`.gomodcache` 等生成目录。

## 当前功能

- 用户创建：`POST /api/users/add`
- 用户登录：`POST /api/users/login`
- 药品新增：`POST /api/drugs/add`
- 药品查询：`GET /api/drugs/get`
- 标本留存申请新增：`POST /api/specimens/add`
- 标本留存申请查询：`GET /api/specimens/get`

## 登录说明

- 首页 `/` 是登录页。
- 登录成功后进入 `/drugs`。
- 登录接口返回用户信息，前端保存到 `localStorage`。
- 密码使用 bcrypt 哈希校验，数据库内不能存明文密码。

## UI 风格

- 管理页参考 Gin-Vue-Admin 风格：
  - 深色左侧菜单栏
  - 白色顶部导航栏
  - 浅灰内容背景
  - 白色业务卡片
- 用户信息展示在右上角，头像 + 用户名 + 下拉箭头，退出登录放入下拉菜单。
