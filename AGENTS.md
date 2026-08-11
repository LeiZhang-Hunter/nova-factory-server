# AGENTS.md

## 项目概览

nova-factory-server（白泽工厂管理系统）是一个基于 Go + Gin 构建的企业级平台，集成了后台管理、工业 AI 智能体、物联网设备管理和 ERP 等业务模块。

| 项 | 值 |
|---|-----|
| 语言 / 框架 | Go 1.26 · Gin · GORM · Wire · Zap |
| 启动入口 | `app/main.go` |
| 配置 | `--config` 命令行指定 YAML 路径，结构体见 `app/setting/settings.go` |
| 默认凭据 | `admin` / `admin123` |
| 许可证 | MIT |

## 构建与模块开关

项目使用 **Go build tags** 按需编译模块。当前常驻模块为 `ai` 和 `iot`，其他模块（shop、erp 等）作为外部 addon 链接接入。

```bash
# 安装 wire / swag 工具链
make init

# 生成 Wire 依赖注入（默认 tag: ai iot）
make wire

# 链接外部 addon 并生成全量 Wire
make wire-addons                        # 全部 addon
make wire-addons ADDONS="shop erp"      # 指定 addon
make wire-addons ADDONS_DIR=/path/to/nova-factory-addons-be

# 编译
go build -tags="ai iot" ./app/

# Swagger 文档生成
make swag                               # 等价 cd app/ && swag init

# MCP handler 元数据嵌入
make gen-mcp-doc

# gRPC Protobuf
make pb
```

Swagger 文档地址: `http://localhost:8080/swagger/doc.json`

## 架构

### 整体分层

```
┌─────────────────────────────────────────────────┐
│  Controller   HTTP handler · 参数校验 · Swagger  │
├─────────────────────────────────────────────────┤
│  Service      业务逻辑接口 + impl 实现            │
├─────────────────────────────────────────────────┤
│  DAO          数据访问接口 + impl 实现 (GORM)     │
├─────────────────────────────────────────────────┤
│  Model        entity / request / response / query│
└─────────────────────────────────────────────────┘
```

### 目录结构

```
app/
├── main.go                 # 入口，Wire 注入，启动 HTTP/gRPC
├── routes/
│   ├── routes.go           # Gin Engine + MCP Server 初始化
│   └── grpc_routes.go      # gRPC Server + Wire 注册
├── setting/
│   └── settings.go         # AppConfig 结构体 + Viper 加载
├── middlewares/
│   ├── session_auth.go     # Session 鉴权中间件（支持 HTTP/WS/多域）
│   ├── permission.go       # 细粒度权限校验中间件
│   └── logger.go           # 请求日志中间件
├── datasource/             # 数据源连接层
│   ├── mysql/              # GORM MySQL 连接（dbs、gorm2 双实例）
│   ├── redis/              # Redis 缓存 + Session 存储
│   ├── clickhouse/         # ClickHouse 时序存储
│   ├── iotdb/              # IoTDB/Prometheus 指标存储
│   ├── milvus/             # Milvus 向量数据库（AI 检索）
│   └── objectFile/         # 文件存储（S3 / 本地）
├── baize/                  # 通用工具（base_entity、分页、时间等）
├── utils/                  # 通用工具库
│   ├── baizeContext/       # 从 gin.Context 提取用户/角色/权限
│   ├── gin_mcp/            # MCP 协议 HTTP 适配 + handler 注解解析
│   ├── llm/                # LLM 客户端（模型注册、embedding）
│   ├── observer/           # 事件观察者（集成平台同步）
│   ├── grpc/               # gRPC 工具 + agent 配置热加载
│   └── ...
├── constant/               # 模块常量（按领域分包）:
│   ├── agent/ aiagent/ gateway/  # AI 相关
│   ├── device/ metric/ iotdb/    # IoT 相关
│   ├── order/ shop/              # 商城/ERP
│   ├── craft/                    # 工艺路由
│   ├── sessionStatus/ dataScopeAspect/ userStatus/
│   └── ...
└── business/               # 业务模块
    ├── admin/              # 后台管理
    ├── ai/                 # AI 智能体
    ├── iot/                # 物联网
    └── ...
```

### 核心业务模块

| 模块 | 路径 | 子模块 |
|------|------|--------|
| **admin** | `business/admin/` | system（用户/角色/菜单/部门/字典）、monitor（定时任务/操作日志）、tool（代码生成器） |
| **ai** | `business/ai/` | agent/dataset（RAGFlow 数据集）、gateway（LLM 网关/MCP/Agent 编排）、core（LLM 客户端抽象） |
| **iot** | `business/iot/` | asset（设备/摄像头/建筑/物料）、metric（IoTDB/ClickHouse 指标）、alert（告警规则/模板）、dashboard、craft（工艺路线）、daemonize（设备守护进程） |

### 每个业务子模块的内部结构

```
business/<module>/<sub>/
├── controller/          # HTTP handler，每个文件一个 Controller 结构体
├── service/             # 业务接口定义 (i_xxx.go)
├── service/impl/        # 业务实现，单文件单结构体
├── dao/                 # 数据接口定义 (i_xxx.go)
├── dao/impl/            # 数据实现，单文件单结构体
└── models/              # entity / request / response / query / api
```

### 鉴权体系

```
请求 → Logger 中间件 → SessionAuth 中间件(Session/Token)
    → Permission 中间件(RBAC 权限校验) → Controller
```

- **SessionAuth**: 支持 HTTP / WebSocket 双传输协议，多认证域（admin / shop），必选/可选认证模式
- **Permission**: 基于 `baizeContext.GetPermission(c)` 做细粒度接口权限校验
- **baizeContext**: 从 `gin.Context` 中提取当前登录用户的 ID、角色、部门、权限、数据权限范围等

## 代码约定

### DAO 规范

- **单文件单结构体**: 每个 DAO 实现独立成文件
- **显式表名**: 所有查询必须 `d.db.Table(d.tableName)`，禁止依赖 GORM 自动推断
- **构造函数注入**: `func NewIXxxDaoImpl(db *gorm.DB) dao.IXxxDao`
- **命名**: 接口 `i_xxx_dao.go`，实现 `i_xxx_dao_impl.go`

```go
type IShopOrderDaoImpl struct {
    db        *gorm.DB
    tableName string
}

func NewIShopOrderDaoImpl(db *gorm.DB) dao.IShopOrderDao {
    return &IShopOrderDaoImpl{
        db:        db,
        tableName: "shop_order",
    }
}

func (d *IShopOrderDaoImpl) FindByID(c *gin.Context, id int64) (*model.Order, error) {
    var order model.Order
    err := d.db.WithContext(c).Table(d.tableName).Where("id = ?", id).First(&order).Error
    return &order, err
}
```

### Swagger 注解

每个 handler 函数上方加注释，然后执行 `make swag` 生成：

```go
// 获取告警规则列表
// @Summary 分页查询告警规则
// @Tags 告警管理
// @Param object body request.AlertQuery true "查询条件"
// @Success 200 {object} response.ResponseData "成功"
// @Router /alert/list [post]
func (a *Alert) List(c *gin.Context) { ... }
```

### MCP 元数据

在 handler 上加 `@summary` / `@operationId` / `@param` / `@tags` 注释，执行 `make gen-mcp-doc` 即可编译期嵌入，使生产无源码环境也能正常暴露 MCP 接口。

### 添加新模块

1. 在 `app/business/` 下创建模块目录，实现 Controller → Service → DAO
2. 编写 `register.go` 提供 Wire `ProviderSet`
3. 编写 `routes.go` 用 `//go:build <tag>` 约束 + `GinProviderSet` 注册路由
4. 在 `app/wire.go` 中引入模块的 `GinProviderSet`
5. `make wire` 生成依赖注入代码

### 技能 (Skills)

项目内置两个专用技能，生成代码时应优先调用：

- **go-swagger-helper**: 补充或维护 Swagger 注解
- **go-service-gen**: 根据 Service 接口生成 ServiceImpl + DAO 接口 + DAOImpl 完整四层代码
