# 信息管理平台 BS 项目

## 项目概况

- 项目类型：信息管理平台 BS 项目
- 前端：Vue 3 + Composition API + axios + Element Plus + Vite
- 后端：Go + Gin + GORM
- 数据库：MySQL 8.0
- 数据库名称：`medical_info`

## 启动脚本

项目提供了一键开发启动脚本，会自动释放前后端端口、启动服务、写入运行日志，并打开登录页。

```powershell
powershell.exe -ExecutionPolicy Bypass -File scripts\start-dev.ps1
```

常用参数：

- `-Help`：查看脚本帮助。
- `-Stop`：停止占用当前配置端口的前后端进程。
- `-NoBrowser`：启动服务但不自动打开浏览器。
- `-BackendOnly`：只重启后端服务。
- `-FrontendOnly`：只重启前端服务。
- `-BackendPort 8888`：指定后端端口，默认 `8888`。
- `-FrontendPort 8080`：指定前端端口，默认 `8080`。

示例：

```powershell
powershell.exe -ExecutionPolicy Bypass -File scripts\start-dev.ps1 -NoBrowser
powershell.exe -ExecutionPolicy Bypass -File scripts\start-dev.ps1 -Stop
powershell.exe -ExecutionPolicy Bypass -File scripts\start-dev.ps1 -BackendOnly
powershell.exe -ExecutionPolicy Bypass -File scripts\start-dev.ps1 -FrontendOnly
```

运行日志输出到 `runtime_logs/` 目录。

## 手动启动

后端启动：

```bash
cd backend
go mod tidy
go run .
```

后端默认监听端口：`8888`。

前端启动：

```bash
cd frontend
npm install
npm.cmd run dev
```

前端默认监听端口：`8080`。

访问地址：

- 登录页：`http://localhost:8080/login`
- 首页入口：`http://localhost:8080/`
- 后端 API：`http://localhost:8888/api`

## 验证命令

后端测试：

```bash
cd backend
go test ./...
```

前端构建：

```bash
cd frontend
npm.cmd run build
```

## 当前功能

- 用户登录、用户新增、用户查询、用户删除。
- 药品新增、药品查询。
- 标本留存申请新增、查询、批量导入预览、批量导入。
- 文件上传、文件列表查询、文件下载。
- 追溯码识别、确认、查询、删除。
- 基于 RBAC 的接口权限控制。
- 前端登录成功后保存用户信息，并进入药品管理页面。

## 接口说明

后端 API 统一返回：

- `code`
- `message`
- `errorMessage`
- `data`

成功响应的 `code` 固定为 `0`。

主要接口：

- `POST /api/users/login`
- `POST /api/users/add`
- `GET /api/users/get`
- `POST /api/users/delete`
- `POST /api/drugs/add`
- `GET /api/drugs/get`
- `POST /api/specimens/add`
- `POST /api/specimens/import/preview`
- `POST /api/specimens/import`
- `GET /api/specimens/get`
- `POST /api/files/upload`
- `GET /api/files/get`
- `GET /api/files/download/:id`
- `POST /api/trace_codes/recognize`
- `POST /api/trace_codes/confirm`
- `GET /api/trace_codes/get`
- `POST /api/trace_codes/delete`

兼容接口：

- `POST /api/fileUploadAndDownload/upload`
- `GET /api/fileUploadAndDownload/get`
- `GET /api/fileUploadAndDownload/download/:id`
