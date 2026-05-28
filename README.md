# 信息管理平台 BS 项目

技术栈：

- 前端：Vue 3 + Composition API + axios + Element Plus + Vite
- 后端：Go + Gin + GORM
- 数据库：MySQL 8.0

## 项目结构

```text
backend/
  main.go
  db/mysql.go
  handlers/drug_handler.go
  handlers/specimen_handler.go
  models/drug.go
  models/specimen_application.go
  models/sys_user.go
frontend/
  index.html
  package.json
  vite.config.js
  src/main.js
  src/App.vue
  src/styles.css
sql/
  schema.sql
  migrate_driver_gene_mutation.sql
```

## 数据库初始化

数据库名称为 `medical_info`，用于保存药品、标本留存申请和系统用户信息。

`sql/schema.sql` 会创建以下表：

- `drugs`：药品信息表
- `specimen_applications`：标本留存申请表
- `sys_users`：系统用户表

## 后端运行

启动：

```bash
cd backend
go mod tidy
go run .
```


后端监听端口：`8888`。

## 前端运行

```bash
cd frontend
npm install
npm run dev
```

## 访问地址：

- 首页：`http://localhost:8080/`

