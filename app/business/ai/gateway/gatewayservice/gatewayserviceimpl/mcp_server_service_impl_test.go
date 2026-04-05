package gatewayserviceimpl

import (
	"testing"

	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
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
		ID:        "1",
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
