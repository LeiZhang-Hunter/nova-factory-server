package gin_mcp

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"nova-factory-server/app/constant/agent"
	key2 "nova-factory-server/app/utils/store/key"
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

	key, ok := h[agent.AuthorizationApiKey]
	if !ok {
		result.Tools = make([]mcp.Tool, 0)
		return
	}

	if len(key) == 0 {
		result.Tools = make([]mcp.Tool, 0)
		return
	}

	// 读取用户id
	w := httptest.NewRecorder()
	ginCtx, _ := gin.CreateTestContext(w)
	ginCtx.Request = &http.Request{Header: h}
	settingTools, err := key2.GetStore().GetTool(ginCtx, key[0])
	if len(settingTools) == 0 {
		result.Tools = make([]mcp.Tool, 0)
		return
	}

	if err != nil {
		zap.L().Error("after hook error", zap.Error(err))
		return
	}

	filtered := make([]mcp.Tool, 0)

	// 获取路由到权限的注册映射
	var permToolsMap = make(map[string]uint8)
	for _, t := range settingTools {
		permToolsMap[t] = 0
	}

	for _, v := range result.Tools {
		_, ok = permToolsMap[v.Name]
		if !ok {
			continue
		}
		filtered = append(filtered, v)
	}

	result.Tools = filtered
}
