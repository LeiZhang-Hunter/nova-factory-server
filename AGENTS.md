# AGENTS.md (nova-factory-server)

工厂后端核心服务：Gin + Wire + GORM/sqly，承载 admin/ai/data/iot 四大业务域，数据层包含 MySQL、Redis、IoTDB、ClickHouse、Milvus。

## 构建与生成

```bash
make init                                # 安装 wire / swag
make wire                                # 同步插件并生成 Wire 依赖注入（-tags="ai iot"）
make swag                                # 生成 Swagger 文档
make gen-mcp-doc                         # 将 handler 注释编译进二进制（生产 MCP 元数据）
make pb                                  # 解析 protobuf
go build -tags="ai iot" ./...
go test ./...
```

## 目录结构

```
app/
  main.go                     # 入口
  business/                   # 业务域（每个域内 Controller → Service → DAO）
    admin/                    # 系统管理（basics / monitor / system / tool）
    ai/                       # AI 智能体
    data/                     # 数据平台（术语 / 本体 / 对象 / 血缘 / 搜索）
    iot/                      # IoT 设备监控
    provider_gen.go           # addonsync 生成的插件聚合，勿手改
  cmd/                        # 附加进程（otel、prediction）
  middlewares/  routes/       # 中间件与路由（gin + grpc）
  utils/gin_mcp/              # MCP 服务
  docs/                       # swag 生成
config/                       # config.yaml / config2.yaml
dao/query/                    # gorm gen 生成代码，勿手改
models/                       # entity + dto
sql/                          # mysql / clickhouse 建表与迁移
tools/                        # addonsync、handler-doc-gen 代码生成器
template/                     # go / sql / vue 代码生成模板
ragflow/  uploads/  public/   # 运行时资源
```

## 架构约定

- 固定 Controller → Service → DAO 分层，每层提供 `provider.go` 供 Wire 组装。
- 业务模块注册到 `app/business/provider_gen.go`（由 `make wire` 生成）。
- 接口需要权限时使用权限中间件，权限码按 `模块:资源:动作` 命名（如 `data:term:query`）。
- 删除类操作默认软删除。
- MCP 接口依赖 handler 上的注释（`@summary` / `@description` / `@param` / `@tags` / `@operationId`），改完需重新 `make gen-mcp-doc`。
- 新增依赖后 `go mod tidy`；注入关系变更后重新 `make wire`。

## 业务模块规范（以 app/business/data 为范本）

新增/修改业务模块时，按以下规范执行。

### 分层与目录

```
<module>/
  register.go      # ProviderSet 汇总（Wire 入口）
  routes.go        # GinProviderSet：路由组 + MCP 权限注册
  controller/      # *Controller + 聚合结构（如 DataControllers）+ provider.go
  service/
    i_<domain>_service.go    # I<Domain>Service 接口
    impl/                    # I<Domain>ServiceImpl 实现 + provider.go
  dao/
    i_<domain>_dao.go        # I<Domain>DAO 接口
    impl/                    # I<Domain>DAOImpl 实现 + provider.go
  models/
    entity/                  # GORM 实体 + 常量
    dto/                     # 请求 / 响应结构
  agent/                     # AI 相关能力（按需）
```

依赖方向固定：Controller → Service 接口 → DAO 接口，构造函数注入、返回接口，全部由 Wire 组装；每层实现目录内必须提供 `provider.go`。

### 命名规则

- 文件：`controller/<domain>_controller.go`、`service/i_<domain>_service.go`、`service/impl/i_<domain>_service_impl.go`、`dao/i_<domain>_dao.go`、`dao/impl/i_<domain>_dao_impl.go`、`models/entity/<domain>.go`、`models/dto/<domain>.go`。
- 类型：`XxxController` / `I<Domain>Service` / `I<Domain>ServiceImpl` / `I<Domain>DAO` / `I<Domain>DAOImpl`。
- CRUD 方法统一 `Create / Get / List / Update / Delete`。
- 路由资源统一 kebab-case 复数：`/data/pipeline-rules`、`/data/service-connections`。
- 权限码统一 `模块:资源:动作`：`data:collector:add / query / edit / remove`；需要兜底时用 `HasPermissions([]string{...})`。
- 状态、类型等常量放 entity 包，语义化命名（`RuleStatusEnabled`、`SourceTypeFile`）。

### 调用规则

- 每个控制器实现 `PrivateRoutes(*gin.RouterGroup)` 和 `PrivateMcpRoutes(*gin_mcp.GinMCP)`；路由统一在模块 `routes.go` 挂 `NewSessionAuthMiddlewareBuilder` 鉴权。
- 路由链固定：`middlewares.SetLog("操作名", middlewares.Insert/Update/Delete)`（审计）→ `middlewares.HasPermission("data:xxx:yyy")`（权限）→ handler。
- 入参：POST/PUT 用 `c.ShouldBindJSON(&req)`，失败直接 `baizeContext.Waring(c, err.Error())`；查询参数用 `cast.ToInt(c.DefaultQuery(...))`。
- 分页：`pageNum` / `pageSize` 为查询参数，`pageSize` 钳制到 1~200，`offset = (pageNum-1)*pageSize`；Service 签名固定 `List(c, offset, limit, 筛选参数...) ([]T, int64, error)`。
- 出参：`baizeContext.Success` / `SuccessData` / `SuccessListData` / `Waring`，统一包装 `response.ResponseData`。
- Service 层负责业务校验（唯一性、引用存在性）、状态默认值、敏感字段加密、版本号递增；DAO 只做 SQL。
- DAO 统一 `WithContext(ctx)`；`gorm.ErrRecordNotFound` 转 `(nil, nil)`，由 Service 判定后返回中文错误。
- 列表实现固定：`Count` + `Order("created_at DESC")` + `Offset/Limit`；`Update` 用 `Select(...)` 白名单列。

### Model 封装

- entity 与 dto 严格分离，dto 内请求/响应分离（`CreateXxxRequest` / `UpdateXxxRequest` / `XxxResponse`）。
- 主键：`uuid.NewString()` 生成 string 主键，`gorm:"type:varchar(36);primaryKey"`。
- 列名 snake_case，`TableName()` 显式返回 `data_*` 表名。
- 统一审计字段 `CreatedAt / UpdatedAt / CreatedBy / UpdatedBy`，软删除 `DeletedAt gorm.DeletedAt`。
- 敏感/大字段打 `json:"-"`（配置、凭据、删除标记）。
- 请求体用 binding 校验：`required` / `max` / `oneof`。
- 动态配置用 `json.RawMessage` 透传，结构化约束在 Service 层校验，错误信息带 `pipelines[i].field` 定位。
- 响应必须手工映射 DTO，不直接返回 entity。

### API 注释

- 每个 handler 配完整 Swagger 注释（模板见下），缺失会导致文档与 MCP 元数据不全：

```go
// @Summary 创建采集器
// @Description 创建采集器设备
// @Tags 数据平台-采集器管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.CreateCollectorRequest true "创建采集器请求"
// @Success 200 {object} response.ResponseData{data=dto.CollectorResponse}
// @Router /data/collectors [post]
```

- `@Tags` 统一「数据平台-模块名」，必须带 `@Security BearerAuth`，返回类型统一 `response.ResponseData{data=...}`（列表为 `{data=[]dto.Xxx}`）。
- MCP 权限在 `PrivateMcpRoutes` 里逐路由显式 `router.RegisterPermission(method, path, "data:xxx:yyy")`。

## 注意

- 生成文件（`wire_gen.go`、`dao/query`、`app/docs`、`provider_gen.go`）不要手改。
- 编译默认带 `-tags="ai iot"`。
- 类型转换遵循根 AGENTS.md 的 cast 规则。
