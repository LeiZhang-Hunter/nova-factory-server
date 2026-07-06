# handler-doc-gen

`gin_mcp` 模块的代码生成工具。从 Gin handler 函数的注释中提取 MCP Tool 元信息（`@summary` / `@operationId` / `@param` / `@return` / `@tags`），生成编译期嵌入的 `handlerDocMap`，使得生产环境二进制部署（无 `.go` 源码）时也能正常提供 MCP 接口的摘要、标签等元数据。

## 使用方式

在 handler 函数上添加注释：

```go
// @summary 获取告警规则列表
// @description 返回分页的告警规则
// @param page 页码
// @tags 告警管理
// @operationId listAlerts
func (a *Alert) List(c *gin.Context) { ... }
```

### 生成

```bash
cd app/utils/gin_mcp/pkg/convert
go generate
```

或从项目根目录：

```bash
go run ./tools/handler-doc-gen/
```

生成文件位于 `app/utils/gin_mcp/pkg/convert/handler_docs_gen.go`。

### 工作原理

- **构建时**：`go generate` 遍历 `app/` 下所有 `.go` 文件，提取注释生成 `handlerDocMap`
- **运行时**：`parseHandlerComments` 优先从 `handlerDocMap` 查表，命中后不再读磁盘源文件

### 支持的注释标记

| 标记 | 说明 |
|------|------|
| `@summary` | 接口摘要 |
| `@description` | 详细描述 |
| `@param name desc` | 参数名及说明 |
| `@return` | 返回值说明 |
| `@tags` | 标签（逗号或空格分隔） |
| `@operationId` | 自定义 Operation ID |

标记名不区分大小写，值保留原始大小写。
