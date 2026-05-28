# 药品信息管理 BS 项目

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
frontend/
  index.html
  package.json
  vite.config.js
  src/main.js
  src/App.vue
  src/styles.css
sql/
  schema.sql
```

## 数据库初始化

在 MySQL 8.0 中执行：

```sql
CREATE DATABASE IF NOT EXISTS medical_info
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE medical_info;

CREATE TABLE IF NOT EXISTS drugs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL COMMENT '药品名称',
  manufacturer VARCHAR(150) DEFAULT '' COMMENT '生产厂家',
  approval_number VARCHAR(100) DEFAULT '' COMMENT '批准文号',
  specification VARCHAR(100) DEFAULT '' COMMENT '规格',
  price DECIMAL(10, 2) NOT NULL COMMENT '价格',
  stock INT NOT NULL COMMENT '库存数量',
  created_at DATETIME(3) DEFAULT NULL,
  updated_at DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (id),
  INDEX idx_drugs_name (name),
  CONSTRAINT chk_drugs_price_positive CHECK (price > 0),
  CONSTRAINT chk_drugs_stock_positive CHECK (stock > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS specimen_applications (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(50) NOT NULL COMMENT '姓名',
  gender VARCHAR(10) NOT NULL COMMENT '性别',
  age INT NOT NULL COMMENT '年龄',
  id_number VARCHAR(40) NOT NULL COMMENT 'ID号',
  sample_type VARCHAR(20) NOT NULL COMMENT '送检标本类型',
  pathology_type VARCHAR(50) NOT NULL COMMENT '病理类型',
  pdl1_expression INT NOT NULL DEFAULT 0 COMMENT 'PD-L1表达百分比',
  driver_gene_mutation VARCHAR(255) DEFAULT '' COMMENT '驱动基因突变',
  stage VARCHAR(10) NOT NULL COMMENT '分期',
  last_treatment VARCHAR(255) DEFAULT '' COMMENT '末次治疗',
  follow_up_treatment VARCHAR(255) DEFAULT '' COMMENT '后续治疗方案',
  doctor VARCHAR(50) NOT NULL COMMENT '送检医师',
  inspection_date VARCHAR(10) NOT NULL COMMENT '送检日期',
  created_at DATETIME(3) DEFAULT NULL,
  updated_at DATETIME(3) DEFAULT NULL,
  PRIMARY KEY (id),
  INDEX idx_specimen_name (name),
  INDEX idx_specimen_id_number (id_number),
  CONSTRAINT chk_specimen_age_positive CHECK (age > 0),
  CONSTRAINT chk_specimen_pdl1_range CHECK (pdl1_expression BETWEEN 0 AND 100),
  CONSTRAINT chk_specimen_gender CHECK (gender IN ('男', '女')),
  CONSTRAINT chk_specimen_sample_type CHECK (sample_type IN ('组织', '血浆')),
  CONSTRAINT chk_specimen_stage CHECK (stage IN ('I', 'II', 'III', 'IV'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

也可以直接执行文件：

```bash
mysql -uroot -p < sql/schema.sql
```

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

本机访问：`http://localhost:8080`

公网或局域网访问：`http://你的本机公网IP:8080`

Vite 已配置监听 `0.0.0.0:8080`，并将 `/api` 代理到 `http://localhost:8888`。后端监听 `0.0.0.0:8888`，CORS 已允许公网来源访问。

如果其他电脑仍无法访问，请检查：

- 本机防火墙是否放行 TCP `8080` 端口
- 云服务器安全组是否放行 TCP `8080` 端口
- 家用宽带/路由器是否已做公网端口转发到本机 `8080`
- 对方访问的是否是真实公网 IP，而不是 `127.0.0.1`、内网 IP 或运营商 CGNAT 地址

## API 接口

### 新增药品

`POST /api/drugs/add`

请求：

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

成功响应：

```json
{
  "message": "保存成功",
  "data": {
    "id": 1,
    "name": "阿莫西林胶囊",
    "manufacturer": "某某制药有限公司",
    "approvalNumber": "国药准字H12345678",
    "specification": "0.25g*24粒",
    "price": 18.5,
    "stock": 100,
    "createdAt": "2026-05-27T14:00:00+08:00",
    "updatedAt": "2026-05-27T14:00:00+08:00"
  }
}
```

### 查询药品列表

`GET /api/drugs/get`

按药品名称搜索：

```text
GET /api/drugs/get?name=阿莫西林
```

成功响应：

```json
{
  "message": "查询成功",
  "data": [
    {
      "id": 1,
      "name": "阿莫西林胶囊",
      "manufacturer": "某某制药有限公司",
      "approvalNumber": "国药准字H12345678",
      "specification": "0.25g*24粒",
      "price": 18.5,
      "stock": 100,
      "createdAt": "2026-05-27T14:00:00+08:00",
      "updatedAt": "2026-05-27T14:00:00+08:00"
    }
  ]
}
```

### 新增标本留存申请单

`POST /api/specimens/add`

请求：

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
