package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/mark3labs/mcp-go/client"
	mcptransport "github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

func ProbeMCPServer(ctx context.Context, req *gatewaymodels.MCPServerProbeRequest) (*gatewaymodels.MCPServerProbeResult, error) {
	if req == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	transportType := normalizeMCPProbeTransport(req.Transport)
	if transportType == "" {
		return nil, fmt.Errorf("传输方式不能为空")
	}
	if transportType != "sse" && transportType != "streamableHttp" {
		return nil, fmt.Errorf("传输方式仅支持 sse 或 streamableHttp")
	}

	targetURL := strings.TrimSpace(req.URL)
	if targetURL == "" {
		return nil, fmt.Errorf("URL不能为空")
	}
	if _, err := url.ParseRequestURI(targetURL); err != nil {
		return nil, fmt.Errorf("URL格式不正确")
	}

	headers, err := parseMCPProbeHeaders(req.Headers)
	if err != nil {
		return nil, err
	}

	timeout := normalizeMCPProbeTimeout(req.Timeout)
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	mcpClient, err := newProbeClient(transportType, targetURL, headers, timeout)
	if err != nil {
		return nil, err
	}
	defer mcpClient.Close()

	if err := mcpClient.Start(timeoutCtx); err != nil {
		return nil, fmt.Errorf("启动MCP客户端失败: %w", err)
	}

	initResult, err := mcpClient.Initialize(timeoutCtx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    "nova-factory-server-mcp-probe",
				Version: "1.0.0",
			},
			Capabilities: mcp.ClientCapabilities{},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("初始化MCP会话失败: %w", err)
	}

	if err := mcpClient.Ping(timeoutCtx); err != nil {
		return nil, fmt.Errorf("MCP心跳检查失败: %w", err)
	}

	result := &gatewaymodels.MCPServerProbeResult{
		Transport:       transportType,
		URL:             targetURL,
		ProtocolVersion: initResult.ProtocolVersion,
		ServerInfo:      initResult.ServerInfo,
		Capabilities:    initResult.Capabilities,
	}

	if initResult.Capabilities.Tools != nil {
		toolsResult, err := mcpClient.ListTools(timeoutCtx, mcp.ListToolsRequest{})
		if err != nil {
			return nil, fmt.Errorf("获取MCP工具列表失败: %w", err)
		}
		result.Tools = toolsResult.Tools
	}

	if initResult.Capabilities.Prompts != nil {
		promptsResult, err := mcpClient.ListPrompts(timeoutCtx, mcp.ListPromptsRequest{})
		if err != nil {
			result.Warnings = append(result.Warnings, "获取MCP提示词列表失败: "+err.Error())
		} else {
			result.Prompts = promptsResult.Prompts
		}
	}

	if initResult.Capabilities.Resources != nil {
		resourcesResult, err := mcpClient.ListResources(timeoutCtx, mcp.ListResourcesRequest{})
		if err != nil {
			result.Warnings = append(result.Warnings, "获取MCP资源列表失败: "+err.Error())
		} else {
			result.Resources = resourcesResult.Resources
		}

		resourceTemplates, err := mcpClient.ListResourceTemplates(timeoutCtx, mcp.ListResourceTemplatesRequest{})
		if err != nil {
			result.Warnings = append(result.Warnings, "获取MCP资源模板列表失败: "+err.Error())
		} else {
			result.ResourceTemplates = resourceTemplates.ResourceTemplates
		}
	}

	return result, nil
}

func newProbeClient(transportType string, targetURL string, headers map[string]string, timeout time.Duration) (*client.Client, error) {
	httpClient := &http.Client{Timeout: timeout}

	switch transportType {
	case "sse":
		return client.NewSSEMCPClient(
			targetURL,
			client.WithHeaders(headers),
			client.WithHTTPClient(httpClient),
		)
	case "streamableHttp":
		return client.NewStreamableHttpClient(
			targetURL,
			mcptransport.WithHTTPHeaders(headers),
			mcptransport.WithHTTPBasicClient(httpClient),
			mcptransport.WithHTTPTimeout(timeout),
		)
	default:
		return nil, fmt.Errorf("不支持的传输方式: %s", transportType)
	}
}

func normalizeMCPProbeTransport(transport string) string {
	switch strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(transport), "-", ""), "_", "")) {
	case "sse":
		return "sse"
	case "streamablehttp", "httpstream", "httpstreaming":
		return "streamableHttp"
	default:
		return strings.TrimSpace(transport)
	}
}

func normalizeMCPProbeTimeout(timeout int32) time.Duration {
	if timeout <= 0 {
		timeout = 30
	}
	return time.Duration(timeout) * time.Second
}

func parseMCPProbeHeaders(raw json.RawMessage) (map[string]string, error) {
	headers := make(map[string]string)
	content := strings.TrimSpace(string(raw))
	if content == "" || content == "null" {
		return headers, nil
	}

	if strings.HasPrefix(content, "\"") {
		var headerString string
		if err := json.Unmarshal(raw, &headerString); err != nil {
			return nil, fmt.Errorf("请求头格式不正确: %w", err)
		}
		content = strings.TrimSpace(headerString)
		if content == "" {
			return headers, nil
		}
	}

	var headerMap map[string]interface{}
	if err := json.Unmarshal([]byte(content), &headerMap); err != nil {
		return nil, fmt.Errorf("请求头必须是JSON对象或JSON对象字符串")
	}
	for key, value := range headerMap {
		headerKey := strings.TrimSpace(key)
		if headerKey == "" || value == nil {
			continue
		}
		headerVal := strings.TrimSpace(fmt.Sprintf("%v", value))
		if headerVal == "" {
			continue
		}
		headers[headerKey] = headerVal
	}
	return headers, nil
}
