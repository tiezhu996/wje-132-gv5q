# SafetyPlatform（建筑施工安全管理平台）

面向工地安全管理员和监理人员的安全管理平台，支持安全隐患排查、事故记录与跟踪、安全培训管理和人员资质证书管理。

## 快速启动（Docker Compose 一键部署）

```bash
cp .env.example .env
docker compose up -d --build
```

启动完成后访问：

- 前端：http://localhost:18702
- 后端 API：http://localhost:19202
- 后端健康检查：http://localhost:19202/healthz
- MySQL：localhost:3306

预置账号（密码见 database/init.sql 与 backend/internal/service/seed.go）：

| 手机号 | 密码 | 角色 |
| --- | --- | --- |
| 13800000001 | Admin@123 | 管理员 |
| 13800000002 | User@123 | 安全管理员 |
| 13800000003 | User@123 | 监理 |
| 13800000004 | User@123 | 工人 |

## 本地开发

后端：

```bash
cd backend && go mod tidy && go run ./cmd/server
```

构建：`cd backend && go build ./...`

前端：

```bash
cd frontend && npm install && npm run dev
```

前端开发服务器通过 Vite 代理将 `/api` 转发到 `http://localhost:19202`。

## 技术栈

| 层 | 技术 |
| --- | --- |
| 前端 | React 18 + TypeScript + Ant Design 5 + ECharts + Zustand + Vite |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL 8.0 |
| 认证 | JWT（github.com/golang-jwt/jwt/v5）+ RBAC |
| 其他依赖 | gin-contrib/cors、golang.org/x/crypto/bcrypt、go-playground/validator/v10 |

## 项目目录结构

```
wje-132/
├── docker-compose.yml
├── .env.example
├── database/init.sql              # MySQL 初始化脚本（建表 + 种子数据）
├── backend/
│   ├── cmd/server/main.go
│   └── internal/
│       ├── config/
│       ├── model/                 # user/safety_incident/safety_inspection/inspection_item/safety_training/worker_certification/audit_log
│       ├── repository/            # 按实体分文件
│       ├── service/               # 业务逻辑 + dashboard + 种子数据 + 单元测试
│       ├── handler/               # 按实体分文件（含 upload_handler、audit_log_handler）
│       ├── router/                # router.go + 按实体分文件
│       ├── middleware/            # auth/rbac/audit_log/error_handler/rate_limiter/upload/cors/request_logger
│       ├── dto/
│       ├── constants/             # 枚举、错误码、日志模板、文案
│       └── util/                  # jwt/logger/formatters/app_error/file_upload
└── frontend/
    └── src/
        ├── api/                   # user/incident/inspection/training/certification/dashboard/auditLog/upload
        ├── stores/                # authStore/userStore/incidentStore/inspectionStore/trainingStore
        ├── types/
        ├── components/common/     # StatusBadge/RiskLevelTag/UserAvatar/EmptyState/AvatarUploader/RoleGuard/ErrorBoundary
        ├── hooks/                 # useIncident/usePagination/useFileUpload/useAuth
        ├── pages/                 # Dashboard/IncidentManage/InspectionManage/TrainingManage/CertReview/Profile/AuditLogs/Login
        ├── router/                # index.tsx + guards.tsx
        ├── utils/                 # getSeverityColor/dateFormat/request
        └── constants/             # incident/user/errorCodes
```

## 环境变量

| 变量 | 说明 | 默认值 |
| --- | --- | --- |
| COMPOSE_PROJECT_NAME | Docker Compose 项目名/容器前缀 | safety-platform |
| DB_NAME | 数据库名 | safety_db |
| DB_USER | 数据库用户 | safety_user |
| DB_PASSWORD | 数据库密码 | safety_pwd |
| DB_ROOT_PASSWORD | 数据库 root 密码 | safety_root |
| JWT_SECRET | JWT 签名密钥 | change_me_to_a_long_random_string |
| JWT_EXPIRE_HOURS | JWT 过期小时数 | 72 |
| APP_CORS_ORIGINS | 允许的跨域来源（逗号分隔） | http://localhost:18702 |
| FRONTEND_PORT | 前端端口 | 18702 |
| BACKEND_PORT | 后端端口 | 19202 |
| DB_PORT | 数据库端口 | 3306 |

## Docker 部署说明

- 端口映射：前端 18702:80、后端 19202:8080、MySQL 3306:3306。
- 数据卷：`db-data` 持久化 MySQL 数据；`upload-data` 持久化上传图片。
- 服务依赖：backend `depends_on` db（service_healthy），frontend `depends_on` backend（service_healthy）。
- 前端 Nginx 将 `/api/` 反代到 `http://backend:8080/`，支持 SPA 路由 `try_files`。
- 常见问题：
  - 端口冲突：修改 `.env` 中 `FRONTEND_PORT/BACKEND_PORT/DB_PORT` 后重新 `docker compose up -d`。
  - 数据库重置：`docker compose down -v` 后重新启动。

## 枚举出现位置清单

### SeverityLevel（near_miss/minor/moderate/major/fatal）
- 后端：`backend/internal/constants/incident.go`、`backend/internal/model/safety_incident.go`、`backend/internal/service/safety_incident_service.go`、`backend/internal/util/formatters.go`、`backend/internal/constants/log_templates.go`、`backend/internal/constants/error_codes.go`、`backend/internal/dto/dto_incident.go`、`database/init.sql`
- 前端：`frontend/src/constants/incident.ts`、`frontend/src/utils/getSeverityColor.ts`、`frontend/src/components/common/RiskLevelTag.tsx`、`frontend/src/pages/IncidentManage.tsx`、`frontend/src/pages/Dashboard.tsx`、`frontend/src/pages/CertReview.tsx`

### IncidentStatus（reported/investigating/resolved/closed）
- 后端：`backend/internal/constants/incident.go`、`backend/internal/model/safety_incident.go`、`backend/internal/service/safety_incident_service.go`、`backend/internal/util/formatters.go`、`backend/internal/constants/log_templates.go`、`backend/internal/constants/error_codes.go`、`database/init.sql`
- 前端：`frontend/src/constants/incident.ts`、`frontend/src/components/common/StatusBadge.tsx`、`frontend/src/pages/IncidentManage.tsx`、`frontend/src/pages/Dashboard.tsx`

### UserRole（admin/safety_manager/inspector/worker）
- 后端：`backend/internal/constants/user.go`、`backend/internal/model/user.go`、`backend/internal/middleware/rbac.go`、`backend/internal/router/*.go`、`backend/internal/util/formatters.go`、`database/init.sql`
- 前端：`frontend/src/constants/user.ts`、`frontend/src/stores/authStore.ts`、`frontend/src/components/common/RoleGuard.tsx`、`frontend/src/router/guards.tsx`、`frontend/src/pages/Login.tsx`

## API 接口清单

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /healthz | 服务健康检查 |
| GET | /api/healthz | Nginx 反代健康检查 |
| GET | /api/v1/healthz | API 版本健康检查 |
| POST | /api/v1/auth/register | 用户注册 |
| POST | /api/v1/auth/login | 用户登录，返回 JWT |
| GET | /api/v1/users/me | 当前登录用户信息 |
| PUT | /api/v1/users/me | 修改个人资料 |
| GET | /api/v1/users | 用户列表（仅管理员） |
| GET | /api/v1/dashboard/stats | 安全概览统计 |
| GET | /api/v1/incidents | 事件分页列表 |
| POST | /api/v1/incidents | 上报安全事件 |
| GET | /api/v1/incidents/:id | 事件详情 |
| POST | /api/v1/incidents/:id/assign | 指派调查 |
| POST | /api/v1/incidents/:id/rectify | 提交整改 |
| POST | /api/v1/incidents/:id/close | 关闭事件 |
| GET | /api/v1/inspections | 检查计划列表 |
| POST | /api/v1/inspections | 创建检查计划 |
| GET | /api/v1/inspections/:id | 检查详情与检查项 |
| POST | /api/v1/inspections/:id/execute | 执行检查 |
| GET | /api/v1/inspections/:id/report | 检查报告 |
| GET | /api/v1/inspection-items/by-inspection/:id | 按检查查询检查项 |
| GET | /api/v1/trainings | 培训列表 |
| POST | /api/v1/trainings | 创建培训 |
| GET | /api/v1/trainings/:id | 培训详情 |
| POST | /api/v1/trainings/:id/record | 记录培训成绩 |
| GET | /api/v1/certifications | 资质分页列表 |
| GET | /api/v1/certifications/by-user | 当前用户资质 |
| POST | /api/v1/certifications | 提交资质 |
| POST | /api/v1/certifications/:id/review | 审核资质 |
| GET | /api/v1/audit-logs | 审计日志（仅管理员） |
| POST | /api/v1/upload/image | 图片上传 |

## 主要功能

- 安全概览：近 30 天事件趋势折线图、风险等级分布饼图、待整改列表、本月培训完成率。
- 事件管理：上报事件、指派调查、提交整改、关闭事件，按严重等级/状态/时间筛选。
- 检查管理：创建检查计划、逐项执行检查（合格/不合格）、得分与检查报告。
- 培训管理：创建培训、记录签到与通过率。
- 资质审核：提交资质、审核、过期预警。
- 审计日志：写操作自动记录（管理员查看）。
- 角色权限：JWT + RBAC（admin/safety_manager/inspector/worker）。

## License

MIT License
