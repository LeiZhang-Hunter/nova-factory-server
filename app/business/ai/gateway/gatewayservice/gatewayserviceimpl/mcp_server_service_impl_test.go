package gatewayserviceimpl

import (
	"testing"

	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

func TestMCPServerPrepareUpsertForStdio(t *testing.T) {
	service := &MCPServerServiceImpl{}
	req := &gatewaymodels.MCPServerUpsert{
		Name:      "stdio-server",
		Transport: "STDIO",
		Command:   "npx",
		Args:      `["-y","demo"]`,
		Env:       `{"NODE_ENV":"prod"}`,
		URL:       "https://example.com",
		Headers:   `{"Authorization":"Bearer token"}`,
	}

	if err := service.prepareUpsert(req, false); err != nil {
		t.Fatalf("prepareUpsert err: %v", err)
	}
	if req.Transport != mcpTransportStdio {
		t.Fatalf("unexpected transport: %s", req.Transport)
	}
	if req.URL != "" {
		t.Fatalf("expected url to be cleared, got %s", req.URL)
	}
	if req.Headers != "" {
		t.Fatalf("expected headers to be cleared, got %s", req.Headers)
	}
	if req.Enabled == nil || !*req.Enabled {
		t.Fatalf("expected enabled default true")
	}
	if req.IsCommon == nil || *req.IsCommon {
		t.Fatalf("expected isCommon default false")
	}
	if req.Timeout != 30 {
		t.Fatalf("expected timeout 30, got %d", req.Timeout)
	}
}

func TestMCPServerPrepareUpsertForStreamableHTTP(t *testing.T) {
	service := &MCPServerServiceImpl{}
	req := &gatewaymodels.MCPServerUpsert{
		ID:        1,
		Name:      "http-server",
		Transport: "streamablehttp",
		Command:   "npx",
		Args:      `["-y","demo"]`,
		Env:       `{"NODE_ENV":"prod"}`,
		URL:       "https://example.com/mcp",
		Headers:   `{"Authorization":"Bearer token"}`,
		Enabled:   boolPtr(false),
	}

	if err := service.prepareUpsert(req, true); err != nil {
		t.Fatalf("prepareUpsert err: %v", err)
	}
	if req.Transport != mcpTransportStreamableHTTP {
		t.Fatalf("unexpected transport: %s", req.Transport)
	}
	if req.Command != "" || req.Args != "" || req.Env != "" {
		t.Fatalf("expected stdio fields to be cleared")
	}
	if req.Enabled == nil || *req.Enabled {
		t.Fatalf("expected enabled to remain false")
	}
}

func TestMCPServerPrepareUpsertValidateJSON(t *testing.T) {
	service := &MCPServerServiceImpl{}
	req := &gatewaymodels.MCPServerUpsert{
		Name:      "stdio-server",
		Transport: "stdio",
		Command:   "npx",
		Args:      `{"bad":"json"}`,
	}

	if err := service.prepareUpsert(req, false); err == nil {
		t.Fatalf("expected json validation error")
	}
}

func TestMCPServerUpdateRejectEnabledConfigModification(t *testing.T) {
	dao := &mockMCPServerDao{
		current: &gatewaymodels.MCPServer{
			ID:          1,
			Name:        "demo",
			Description: "old",
			Transport:   mcpTransportStdio,
			Command:     "npx",
			Args:        `["-y","demo"]`,
			Env:         `{"NODE_ENV":"prod"}`,
			Timeout:     30,
			IsCommon:    boolPtr(false),
			Enabled:     boolPtr(true),
		},
	}
	service := &MCPServerServiceImpl{dao: dao}
	req := &gatewaymodels.MCPServerUpsert{
		ID:          1,
		Name:        "demo",
		Description: "new",
		Transport:   mcpTransportStdio,
		Command:     "npx",
		Args:        `["-y","demo"]`,
		Env:         `{"NODE_ENV":"prod"}`,
		Timeout:     30,
		IsCommon:    boolPtr(false),
		Enabled:     boolPtr(true),
	}

	_, err := service.Update(&gin.Context{}, req)
	if err == nil || err.Error() != "MCP服务已启用，请先关闭后再修改" {
		t.Fatalf("unexpected err: %v", err)
	}
	if dao.updateCalled {
		t.Fatalf("unexpected update call")
	}
}

func TestMCPServerUpdateAllowDisableOnly(t *testing.T) {
	dao := &mockMCPServerDao{
		current: &gatewaymodels.MCPServer{
			ID:          1,
			Name:        "demo",
			Description: "same",
			Transport:   mcpTransportStdio,
			Command:     "npx",
			Args:        `["-y","demo"]`,
			Env:         `{"NODE_ENV":"prod"}`,
			Timeout:     30,
			IsCommon:    boolPtr(false),
			Enabled:     boolPtr(true),
		},
	}
	service := &MCPServerServiceImpl{dao: dao}
	req := &gatewaymodels.MCPServerUpsert{
		ID:          1,
		Name:        "demo",
		Description: "same",
		Transport:   mcpTransportStdio,
		Command:     "npx",
		Args:        `["-y","demo"]`,
		Env:         `{"NODE_ENV":"prod"}`,
		Timeout:     30,
		IsCommon:    boolPtr(false),
		Enabled:     boolPtr(false),
	}

	data, err := service.Update(&gin.Context{}, req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !dao.updateCalled {
		t.Fatalf("expected update call")
	}
	if data == nil || data.Enabled == nil || *data.Enabled {
		t.Fatalf("expected disabled result")
	}
}

func TestMCPServerDeleteRejectEnabledConfig(t *testing.T) {
	dao := &mockMCPServerDao{
		current: &gatewaymodels.MCPServer{
			ID:      1,
			Enabled: boolPtr(true),
		},
	}
	service := &MCPServerServiceImpl{dao: dao}

	err := service.DeleteByIDs(&gin.Context{}, []int64{1})
	if err == nil || err.Error() != "MCP服务已启用，请先关闭后再删除" {
		t.Fatalf("unexpected err: %v", err)
	}
	if dao.deleteCalled {
		t.Fatalf("unexpected delete call")
	}
}

type mockMCPServerDao struct {
	current      *gatewaymodels.MCPServer
	updateCalled bool
	deleteCalled bool
}

func (m *mockMCPServerDao) Create(c *gin.Context, req *gatewaymodels.MCPServerUpsert) (*gatewaymodels.MCPServer, error) {
	return nil, nil
}

func (m *mockMCPServerDao) Update(c *gin.Context, req *gatewaymodels.MCPServerUpsert) (*gatewaymodels.MCPServer, error) {
	m.updateCalled = true
	result := *m.current
	result.Enabled = cloneBoolPtr(req.Enabled)
	result.Description = req.Description
	return &result, nil
}

func (m *mockMCPServerDao) DeleteByIDs(c *gin.Context, ids []int64) error {
	m.deleteCalled = true
	return nil
}

func (m *mockMCPServerDao) GetByID(c *gin.Context, id int64) (*gatewaymodels.MCPServer, error) {
	return m.current, nil
}

func (m *mockMCPServerDao) List(c *gin.Context, req *gatewaymodels.MCPServerQuery) (*gatewaymodels.MCPServerListData, error) {
	return nil, nil
}
