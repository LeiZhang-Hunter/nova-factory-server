package gatewayserviceimpl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	mcpTool "github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
	"net/url"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"
	"nova-factory-server/app/utils/mcp"
	"nova-factory-server/app/utils/store/permissions"
	"strings"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
)

const (
	mcpTransportStdio          = "stdio"
	mcpTransportStreamableHTTP = "streamableHttp"
	mcpTransportSse            = "sse"
)

type MCPServerServiceImpl struct {
	dao gatewaydao.IMCPServerDao
}

func NewMCPServerService(dao gatewaydao.IMCPServerDao) gatewayservice.IMCPServerService {
	return &MCPServerServiceImpl{dao: dao}
}

func (m *MCPServerServiceImpl) Create(c *gin.Context, req *gatewaymodels.MCPServerUpsert) (*gatewaymodels.MCPServer, error) {
	if err := m.prepareUpsert(req, false); err != nil {
		return nil, err
	}
	return m.dao.Create(c, req)
}

func (m *MCPServerServiceImpl) Update(c *gin.Context, req *gatewaymodels.MCPServerUpsert) (*gatewaymodels.MCPServer, error) {
	if err := m.prepareUpsert(req, true); err != nil {
		return nil, err
	}
	current, err := m.dao.GetByID(c, req.ID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, errors.New("MCP服务不存在")
	}
	if mcpServerEnabled(current.Enabled) {
		if isSameMCPServerConfig(current, req) {
			return current, nil
		}
		if !isDisableOnlyUpdate(current, req) {
			return nil, errors.New("MCP服务已启用，请先关闭后再修改")
		}
	}
	return m.dao.Update(c, req)
}

func (m *MCPServerServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的MCP服务")
	}
	for _, id := range ids {
		if id == 0 {
			return errors.New("MCP服务ID不能为空")
		}
		current, err := m.dao.GetByID(c, id)
		if err != nil {
			return err
		}
		if current == nil {
			return errors.New("MCP服务不存在")
		}
		if mcpServerEnabled(current.Enabled) {
			return errors.New("MCP服务已启用，请先关闭后再删除")
		}
	}
	return m.dao.DeleteByIDs(c, ids)
}

func (m *MCPServerServiceImpl) List(c *gin.Context, req *gatewaymodels.MCPServerQuery) (*gatewaymodels.MCPServerListData, error) {
	if req == nil {
		req = new(gatewaymodels.MCPServerQuery)
	}
	req.Transport = normalizeTransport(req.Transport)
	return m.dao.List(c, req)
}

func (m *MCPServerServiceImpl) Probe(ctx context.Context, req *gatewaymodels.MCPServerProbeRequest) (*gatewaymodels.MCPServerProbeResult, error) {
	return mcp.ProbeMCPServer(ctx, req)
}

func (m *MCPServerServiceImpl) prepareUpsert(req *gatewaymodels.MCPServerUpsert, isUpdate bool) error {
	if req == nil {
		return errors.New("参数不能为空")
	}
	if isUpdate && req.ID == 0 {
		return errors.New("id不能为空")
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return errors.New("MCP服务名称不能为空")
	}

	req.Transport = normalizeTransport(req.Transport)
	if req.Transport == "" {
		return errors.New("传输方式不能为空")
	}
	if req.Transport != mcpTransportStdio && req.Transport != mcpTransportStreamableHTTP && req.Transport != mcpTransportSse {
		return errors.New("传输方式仅支持 stdio 或 streamableHttp 或sse")
	}

	req.Description = strings.TrimSpace(req.Description)
	req.Command = strings.TrimSpace(req.Command)
	req.Args = strings.TrimSpace(req.Args)
	req.Env = strings.TrimSpace(req.Env)
	req.URL = strings.TrimSpace(req.URL)
	req.Headers = strings.TrimSpace(req.Headers)
	if req.Timeout <= 0 {
		req.Timeout = 30
	}

	if req.IsCommon == nil {
		req.IsCommon = boolPtr(false)
	}
	if req.Enabled == nil {
		req.Enabled = boolPtr(true)
	}

	switch req.Transport {
	case mcpTransportStdio:
		if req.Command == "" {
			return errors.New("stdio模式启动命令不能为空")
		}
		if err := validateJSONArray(req.Args, "stdio模式参数"); err != nil {
			return err
		}
		if err := validateJSONObject(req.Env, "stdio模式环境变量"); err != nil {
			return err
		}
		req.URL = ""
		req.Headers = ""
	case mcpTransportStreamableHTTP:
		if req.URL == "" {
			return errors.New("streamableHttp模式URL不能为空")
		}
		if _, err := url.ParseRequestURI(req.URL); err != nil {
			return errors.New("streamableHttp模式URL格式不正确")
		}
		if err := validateJSONObject(req.Headers, "streamableHttp请求头"); err != nil {
			return err
		}
		req.Command = ""
		req.Args = ""
		req.Env = ""
	}

	return nil
}

func normalizeTransport(transport string) string {
	switch strings.ToLower(strings.TrimSpace(transport)) {
	case "stdio":
		return mcpTransportStdio
	case "streamablehttp":
		return mcpTransportStreamableHTTP
	default:
		return strings.TrimSpace(transport)
	}
}

func validateJSONArray(content string, fieldName string) error {
	if content == "" {
		return nil
	}
	var data []interface{}
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return errors.New(fieldName + "必须是JSON数组字符串")
	}
	return nil
}

func validateJSONObject(content string, fieldName string) error {
	if content == "" {
		return nil
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return errors.New(fieldName + "必须是JSON对象字符串")
	}
	return nil
}

func boolPtr(v bool) *bool {
	return &v
}

func mcpServerEnabled(enabled *bool) bool {
	return enabled != nil && *enabled
}

func isDisableOnlyUpdate(current *gatewaymodels.MCPServer, req *gatewaymodels.MCPServerUpsert) bool {
	if req.Enabled == nil || *req.Enabled {
		return false
	}
	currentSnapshot := snapshotMCPServer(current)
	reqSnapshot := snapshotMCPServerRequest(req)
	reqSnapshot.Enabled = true
	return currentSnapshot == reqSnapshot
}

func isSameMCPServerConfig(current *gatewaymodels.MCPServer, req *gatewaymodels.MCPServerUpsert) bool {
	return snapshotMCPServer(current) == snapshotMCPServerRequest(req)
}

func snapshotMCPServer(current *gatewaymodels.MCPServer) gatewaymodels.McpServerSnapshot {
	if current == nil {
		return gatewaymodels.McpServerSnapshot{}
	}
	req := &gatewaymodels.MCPServerUpsert{
		ID:          current.ID,
		Name:        current.Name,
		Description: current.Description,
		Transport:   current.Transport,
		Command:     current.Command,
		Args:        current.Args,
		Env:         current.Env,
		URL:         current.URL,
		Headers:     current.Headers,
		Timeout:     current.Timeout,
		IsCommon:    cloneBoolPtr(current.IsCommon),
		Enabled:     cloneBoolPtr(current.Enabled),
	}
	_ = normalizeUpsertForCompare(req)
	return snapshotMCPServerRequest(req)
}

func snapshotMCPServerRequest(req *gatewaymodels.MCPServerUpsert) gatewaymodels.McpServerSnapshot {
	return gatewaymodels.McpServerSnapshot{
		Name:        req.Name,
		Description: req.Description,
		Transport:   req.Transport,
		Command:     req.Command,
		Args:        req.Args,
		Env:         req.Env,
		URL:         req.URL,
		Headers:     req.Headers,
		Timeout:     req.Timeout,
		IsCommon:    mcpServerEnabled(req.IsCommon),
		Enabled:     mcpServerEnabled(req.Enabled),
	}
}

func normalizeUpsertForCompare(req *gatewaymodels.MCPServerUpsert) error {
	if req == nil {
		return nil
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.Transport = normalizeTransport(req.Transport)
	req.Command = strings.TrimSpace(req.Command)
	req.Args = strings.TrimSpace(req.Args)
	req.Env = strings.TrimSpace(req.Env)
	req.URL = strings.TrimSpace(req.URL)
	req.Headers = strings.TrimSpace(req.Headers)
	if req.Timeout <= 0 {
		req.Timeout = 30
	}
	if req.IsCommon == nil {
		req.IsCommon = boolPtr(false)
	}
	if req.Enabled == nil {
		req.Enabled = boolPtr(true)
	}
	switch req.Transport {
	case mcpTransportStdio:
		req.URL = ""
		req.Headers = ""
	case mcpTransportStreamableHTTP:
		req.Command = ""
		req.Args = ""
		req.Env = ""
	}
	return nil
}

func cloneBoolPtr(v *bool) *bool {
	if v == nil {
		return nil
	}
	value := *v
	return &value
}

// ProbePerm 权限探测
func (m *MCPServerServiceImpl) ProbePerm(ctx *gin.Context, mcpServer *gin_mcp.GinMCP, req *gatewaymodels.MCPServerProbeRequest) (*gatewaymodels.MCPServerProbeResult, error) {
	tools := mcpServer.GetTools()
	if len(tools) == 0 {
		return &gatewaymodels.MCPServerProbeResult{
			Tools: make([]mcpTool.Tool, 0),
		}, nil
	}

	// 用户权限
	userPerms := permissions.GetStore().GetPermission(ctx, baizeContext.GetUserId(ctx))
	if len(userPerms) == 0 {
		return &gatewaymodels.MCPServerProbeResult{
			Tools: make([]mcpTool.Tool, 0),
		}, nil
	}
	permissionsData := mcpServer.GetAllPermissions()

	// 构建用户权限映射以便快速查找
	userPermMap := make(map[string]bool, len(userPerms))
	for _, p := range userPerms {
		userPermMap[p] = true
	}

	// 按权限过滤工具列表
	// tool.Name 是 operationId (如 "GET_erp_master_product_list")
	// 通过 m.operations 查找对应的 Method + Path
	// 再通过 registeredPermissions 查找该路由需要的权限
	// 最后检查用户是否拥有该权限
	filtered := make([]mcpTool.Tool, 0)
	for _, tool := range tools {
		operations := mcpServer.GetOperations()
		if operations == nil {
			continue
		}
		op, ok := operations[tool.Name]
		if !ok {
			continue
		}

		// 构造 "METHOD path" 键来查找注册的权限
		permKey := fmt.Sprintf("%s %s", op.Method, op.Path)
		requiredPerm, hasRegistered := permissionsData[permKey]

		if !hasRegistered {
			continue
		}

		// 检查用户是否拥有该权限
		if userPermMap[requiredPerm] {
			filtered = append(filtered, gin_mcp.ToMCPTool(tool))
		} else {
			zap.L().Debug(fmt.Sprintf("[PermissionHook] Tool '%s' (%s) requires permission '%s', user does not have it - filtered out", tool.Name, permKey, requiredPerm))

		}
	}

	return &gatewaymodels.MCPServerProbeResult{
		Tools: filtered,
	}, nil

}
