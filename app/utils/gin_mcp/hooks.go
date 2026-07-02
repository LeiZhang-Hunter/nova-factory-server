package gin_mcp

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/mcp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	key2 "nova-factory-server/app/utils/store/key"
	"nova-factory-server/app/utils/store/permissions"
)

type contextKey int

const (
	// This const is used as key for context value lookup
	requestHeader contextKey = iota
)

func (m *GinMCP) afterHook(ctx context.Context, id any, message *mcp.ListToolsRequest, result *mcp.ListToolsResult) {
	h := message.Header
	if h == nil {
		result.Tools = make([]mcp.Tool, 0)
		return
	}

	key, ok := h["Authorization"]
	if !ok {
		result.Tools = make([]mcp.Tool, 0)
		return
	}

	if len(key) == 0 {
		result.Tools = make([]mcp.Tool, 0)
		return
	}

	// 读取用户id
	userId := key2.GetStore().GetUserId(key[0])
	if userId == 0 {
		result.Tools = make([]mcp.Tool, 0)
		return
	}

	// 获取用户权限列表
	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = &http.Request{Header: h}

	userPerms := permissions.GetStore().GetPermission(ginCtx, userId)

	// 如果没有权限，返回空列表
	if len(userPerms) == 0 {
		result.Tools = make([]mcp.Tool, 0)
		return
	}

	// 构建用户权限映射以便快速查找
	userPermMap := make(map[string]bool, len(userPerms))
	for _, p := range userPerms {
		userPermMap[p] = true
	}

	// 获取路由到权限的注册映射
	m.permissionsMu.RLock()
	registeredPerms := make(map[string]string, len(m.registeredPermissions))
	for k, v := range m.registeredPermissions {
		registeredPerms[k] = v
	}
	m.permissionsMu.RUnlock()

	// 按权限过滤工具列表
	// tool.Name 是 operationId (如 "GET_erp_master_product_list")
	// 通过 m.operations 查找对应的 Method + Path
	// 再通过 registeredPermissions 查找该路由需要的权限
	// 最后检查用户是否拥有该权限
	filtered := make([]mcp.Tool, 0)
	for _, tool := range result.Tools {
		op, ok := m.operations[tool.Name]
		if !ok {
			// 没有 operation 记录的工具，放宽处理：允许通过
			if isDebugMode() {
				log.Printf("[PermissionHook] Tool '%s' has no operation mapping, allowing by default", tool.Name)
			}
			filtered = append(filtered, tool)
			continue
		}

		// 构造 "METHOD path" 键来查找注册的权限
		permKey := fmt.Sprintf("%s %s", op.Method, op.Path)
		requiredPerm, hasRegistered := registeredPerms[permKey]

		if !hasRegistered {
			// 该路由未注册权限要求，视为公开接口，允许通过
			if isDebugMode() {
				log.Printf("[PermissionHook] Tool '%s' (%s) has no registered permission, allowing by default", tool.Name, permKey)
			}
			continue
		}

		// 检查用户是否拥有该权限
		if userPermMap[requiredPerm] {
			filtered = append(filtered, tool)
		} else {
			if isDebugMode() {
				log.Printf("[PermissionHook] Tool '%s' (%s) requires permission '%s', user does not have it - filtered out", tool.Name, permKey, requiredPerm)
			}
		}
	}

	result.Tools = filtered
}
