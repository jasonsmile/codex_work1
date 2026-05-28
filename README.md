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

在 MySQL 8.0 中执行：

```bash
mysql -uroot -proot < sql/schema.sql
```

`sql/schema.sql` 会创建以下表：

- `drugs`：药品信息表
- `specimen_applications`：标本留存申请表
- `sys_users`：系统用户表

用户表建表语句：

```sql
CREATE TABLE IF NOT EXISTS sys_users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  uuid VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户UUID',
  username VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户登录名',
  password VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户登录密码',
  header_img VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT 'https://api.dicebear.com/10.x/bottts/png' COMMENT '用户头像',
  authority_id BIGINT UNSIGNED DEFAULT '888' COMMENT '用户角色ID',
  phone VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户手机号',
  email VARCHAR(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户邮箱',
  enable BIGINT DEFAULT '1' COMMENT '用户是否被冻结 1正常 2冻结',
  created_at DATETIME(3) DEFAULT NULL,
  updated_at DATETIME(3) DEFAULT NULL,
  deleted_at DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_sys_users_uuid (uuid),
  KEY idx_sys_users_username (username)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

后端启动时也会通过 GORM `AutoMigrate` 自动检查并创建 `drugs`、`specimen_applications`、`sys_users` 三张表。

## 后端运行

默认连接：

- 用户名：`root`
- 密码：`root`
- 地址：`127.0.0.1:3306`
- 数据库：`medical_info`

启动：

```bash
cd backend
go mod tidy
go run .
```

自定义数据库连接：

```bash
MYSQL_DSN="root:你的密码@tcp(127.0.0.1:3306)/medical_info?charset=utf8mb4&parseTime=True&loc=Local" go run .
```

或使用拆分环境变量：

```bash
MYSQL_USER=root
MYSQL_PASSWORD=你的密码
MYSQL_HOST=127.0.0.1
MYSQL_PORT=3306
MYSQL_DATABASE=medical_info
go run .
```

后端监听端口：`8888`。

## 前端运行

```bash
cd frontend
npm install
npm run dev
```

访问地址：

- 首页：`http://localhost:8080/`
- 药品信息管理：`http://localhost:8080/drugs`
- 标本留存信息：`http://localhost:8080/specimens`
- 关于我们：`http://localhost:8080/about`

公网或局域网访问：`http://你的本机公网IP:8080`

Vite 已配置监听 `0.0.0.0:8080`，并将 `/api` 代理到 `http://localhost:8888`。后端监听 `0.0.0.0:8888`，CORS 已允许外部来源访问。

如果其他电脑仍无法访问，请检查：

- 本机防火墙是否放行 TCP `8080` 端口
- 云服务器安全组是否放行 TCP `8080` 端口
- 路由器是否已把公网 `8080` 端口转发到本机
- 对方访问的是否是真实公网 IP，而不是 `127.0.0.1`、内网 IP 或运营商 CGNAT 地址

## API 接口

### 新增药品

`POST /api/drugs/add`

```json
{
  "name": "阿莫西林胶囊",
  "manufacturer": "某某制药有限公司",
  "approvalNumber": "国药准字H12345678",
  "specification": "0.25g*24粒",
  "price": 18.5,
  "stock": 100
}
```

### 查询药品列表

`GET /api/drugs/get`

按药品名称搜索：

```text
GET /api/drugs/get?name=阿莫西林
```

### 新增标本留存申请单

`POST /api/specimens/add`

```json
{
  "name": "张三",
  "gender": "男",
  "age": 58,
  "idNumber": "ID20260528001",
  "sampleType": "组织",
  "pathologyType": "腺癌",
  "pdl1Expression": 60,
  "driverGeneMutation": "EGFR 19del",
  "stage": "III",
  "lastTreatment": "一线化疗",
  "followUpTreatment": "靶向治疗",
  "doctor": "李医生",
  "inspectionDate": "2026-05-28"
}
```

### 查询标本留存申请单列表

`GET /api/specimens/get`
