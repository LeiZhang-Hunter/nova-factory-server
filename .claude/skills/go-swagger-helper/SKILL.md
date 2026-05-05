---
name: "go-swagger-helper"
description: "维护 Go 项目的 Swagger 注解与生成流程。Invoke when user asks to add/update Swagger comments, routes, docs generation, or swag init troubleshooting."
---

# Go Swagger Helper

该 Skill 用于维护当前仓库的 Golang Swagger 能力，基于项目现有的 `swaggo/swag`、`gin-swagger` 与 `swag init` 工作流。

## 何时调用

- 用户要求为 Go 接口补充或修改 Swagger 注解
- 用户要求排查 `swag init` 生成失败
- 用户要求补充 `@Summary`、`@Param`、`@Success`、`@Router` 等注解
- 用户要求分析 Swagger 文档入口、访问地址或生成命令
- 用户要求让新接口出现在 Swagger 文档中

## 项目约定

- Swagger 路由接入位于 `app/routes/routes.go`
- Swagger 文档包通过 `nova-factory-server/app/docs` 引入
- 项目使用 `github.com/swaggo/files` 和 `github.com/swaggo/gin-swagger`
- 生成命令优先使用 `Makefile` 中的 `make swag`
- 对应命令实际执行为：

```bash
cd app/ && swag init
```

- README 中说明的 Swagger 地址为：

```text
http://localhost:8080/swagger/doc.json
```

## 工作步骤

### 1. 识别 Swagger 入口与注解风格

- 先检查控制器函数上方是否已有 Swagger 注解
- 保持现有中文注释风格
- 保持如下常见结构：

```go
// XXX 接口说明
// @Summary 接口摘要
// @Description 接口描述
// @Tags 模块名称
// @Param object body somepackage.Request true "请求参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "成功"
// @Router /path [post]
```

### 2. 修改或补全注解

- 优先复用现有请求/响应模型
- `@Param` 要与 `ShouldBindJSON`、`ShouldBindQuery`、路径参数方式一致
- `@Router` 要与真实路由一致
- `@Tags` 要沿用相邻接口的分组名称
- 不要凭空创造不存在的响应类型

### 3. 生成或验证文档

- 优先执行：

```bash
make swag
```

- 若只需直接执行生成命令，可使用：

```bash
cd app/ && swag init
```

- 若仓库缺少 `swag` 命令，可参考 `Makefile` 中的安装方式：

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 4. 排错准则

- 若 `swag init` 报模型解析错误，先检查结构体字段类型、匿名字段、循环引用
- 若接口未出现在文档中，先检查控制器注解与真实路由是否匹配
- 若 `@Param` 无法解析，检查导出的请求结构体是否存在且包路径正确
- 若 `docs` 包未更新，确认命令执行目录是否为 `app/`

## 输出要求

- 修改 Swagger 时，优先只改必要的注解与相关模型
- 完成后说明变更位置、生成命令、验证结果
- 引用代码时使用可点击文件链接

## 示例请求

- “给这个 gin handler 补齐 swagger 注解”
- “帮我修复 swag init 报错”
- “这个接口为什么没出现在 swagger 里”
- “帮我给新增接口接入 swagger 文档”
