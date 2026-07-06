# nova-factory-server

## 构建命令

```bash
# 初始化工具链
make init

# 生成 Wire 依赖注入代码
make wire

# 生成 Swagger 文档
make swag

# 生成 MCP Handler 注释映射（生产环境二进制部署必需）
make gen-mcp-doc

# 解析 Protobuf
make pb
```

### gen-mcp-doc

`gin_mcp` 模块的代码生成工具，将 Gin handler 函数上的注释（`@summary` / `@operationId` / `@param` / `@return` / `@tags`）编译期嵌入二进制，使生产环境无源码时也能正常提供 MCP 接口的元数据。

在 handler 上添加注释：

```go
// @summary 获取告警规则列表
// @description 返回分页的告警规则
// @param page 页码
// @tags 告警管理
// @operationId listAlerts
func (a *Alert) List(c *gin.Context) { ... }
```

执行生成：

```bash
make gen-mcp-doc
```

或从子目录：

```bash
cd app/utils/gin_mcp/pkg/convert && go generate
```

详细说明见 [tools/handler-doc-gen/README.md](tools/handler-doc-gen/README.md)。
