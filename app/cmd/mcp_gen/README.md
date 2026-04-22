# mcp_gen

`mcp_gen` 用于扫描项目中的 Gin 路由和 `gin_mcp.RegisterSchema()` 注册信息，生成 MCP 所需的两份文件：

- `mcp.json`：MCP tools 定义
- `operations.json`：tool 与实际 HTTP 路由的映射关系

这两个文件通常会被项目启动时的 `McpServer.Mount()` 使用。

## 工作逻辑

`mcp_gen` 的执行流程如下：

1. 读取配置文件和命令行参数
2. 创建一个临时 Gin engine
3. 初始化 `gin_mcp.GinMCP`
4. 调用 `SetupServer()` 扫描指定目录下的路由与 schema
5. 生成 tools 数据并写入 `mcp.json`
6. 生成 operations 数据并写入 `operations.json`

入口代码见 [main.go](file:///home/zhanglei/project/zhanglei/nova-factory-server/app/cmd/mcp_gen/main.go)。

## 配置来源

支持两种方式：

- 配置文件
- 命令行参数覆盖

优先级如下：

1. 命令行参数
2. `config.yaml` 中的 `mcpGen`
3. 回退到已有配置：
   - `mcp.path`
   - `mcp.operationsPath`
   - `host`
4. 默认值

## 配置项

`mcpGen` 支持以下字段：

```yaml
mcpGen:
  name: Product API
  description: API for managing products.
  baseURL: http://localhost:8080
  scanPath: ../app
  toolsOutputPath: /home/zhanglei/project/zhanglei/nova-factory-server/config/mcp.json
  operationsOutputPath: /home/zhanglei/project/zhanglei/nova-factory-server/config/operations.json
```

字段说明：

- `name`：MCP server 名称
- `description`：MCP server 描述
- `baseURL`：生成 operation 时使用的基础地址
- `scanPath`：GinMCP 扫描源码的目录
- `toolsOutputPath`：生成的 `mcp.json` 输出路径
- `operationsOutputPath`：生成的 `operations.json` 输出路径

## 已有配置回退规则

如果 `mcpGen` 没有配置完整，程序会自动补齐：

- `toolsOutputPath` 为空时，使用 `mcp.path`
- `operationsOutputPath` 为空时，使用 `mcp.operationsPath`
- `baseURL` 为空时，使用全局 `host`

当前项目里已有：

```yaml
mcp:
  path: /home/zhanglei/project/zhanglei/nova-factory-server/config/mcp.json
  operationsPath: /home/zhanglei/project/zhanglei/nova-factory-server/config/operations.json
```

## 默认值

当配置文件和命令行都没有提供时，默认值如下：

```text
name=Product API
description=API for managing products.
baseURL=http://localhost:8080
scanPath=../app
toolsOutputPath=./config/mcp.json
operationsOutputPath=./config/operations.json
```

## 命令行参数

支持以下参数：

```bash
--config
--name
--description
--base-url
--scan-path
--tools-output
--operations-output
```

参数含义：

- `--config`：配置文件路径，默认 `./config/config.yaml`
- `--name`：覆盖 MCP server 名称
- `--description`：覆盖 MCP server 描述
- `--base-url`：覆盖基础地址
- `--scan-path`：覆盖扫描目录
- `--tools-output`：覆盖 `mcp.json` 输出路径
- `--operations-output`：覆盖 `operations.json` 输出路径

## 使用示例

使用默认配置文件：

```bash
go run ./app/cmd/mcp_gen --config ./config/config.yaml
```

使用配置文件并覆盖部分参数：

```bash
go run ./app/cmd/mcp_gen \
  --config ./config/config.yaml \
  --base-url http://127.0.0.1:8080 \
  --tools-output ./config/mcp.json \
  --operations-output ./config/operations.json
```

直接执行已编译二进制：

```bash
./app/cmd/mcp_gen/mcp_gen --config ./config/config.yaml
```

## 输出结果

执行成功后会生成：

- `mcp.json`
- `operations.json`

程序在写文件前会自动创建目标目录，因此输出目录不需要手工提前创建。

## 与业务代码的关系

要让某个接口被更好地转换成 MCP tool，通常需要在控制器里显式注册 schema，例如：

```go
router.RegisterSchema("POST", "/erp/order-audit/import", nil, ordermodels.OrderAuditImportReq{})
```

含义是：

- 路径是 `POST /erp/order-audit/import`
- 没有 query 参数
- body 使用 `OrderAuditImportReq`

这样生成出来的 MCP 参数描述会更完整，尤其是结构体字段上补了 `jsonschema:"description=..."` 之后。

## 常见问题

### 1. 为什么生成出来的工具参数不完整？

通常是因为接口没有调用 `RegisterSchema()`，或者请求结构体缺少 `jsonschema:"description=..."` 标签。

### 2. 为什么输出文件没有写到预期位置？

优先检查：

- 是否传了命令行覆盖参数
- `mcpGen` 配置是否生效
- 是否回退到了 `mcp.path` / `mcp.operationsPath`

### 3. `scanPath` 应该填什么？

应填写给 `GinMCP` 做源码扫描的目录。当前工具默认值是：

```text
../app
```

如果你从项目根目录执行，且目录结构发生变化，需要按实际路径调整。
