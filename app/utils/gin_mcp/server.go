package gin_mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"nova-factory-server/app/utils/gin_mcp/pkg/convert"
	"nova-factory-server/app/utils/gin_mcp/pkg/transport"
	"nova-factory-server/app/utils/gin_mcp/pkg/types"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

// isDebugMode returns true if Gin is in debug mode
func isDebugMode() bool {
	return gin.Mode() == gin.DebugMode
}

// GinMCP represents the MCP server configuration for a Gin application
type GinMCP struct {
	engine            *gin.Engine
	name              string
	description       string
	baseURL           string
	tools             []types.Tool
	operations        map[string]types.Operation
	transport         transport.Transport
	config            *Config
	registeredSchemas map[string]types.RegisteredSchemaInfo
	schemasMu         sync.RWMutex
	path              string
	// executeToolFunc holds the function used to execute a tool.
	// It defaults to defaultExecuteTool but can be overridden for testing.
	executeToolFunc func(ctx context.Context, operationID string, parameters map[string]interface{}, headers http.Header) (interface{}, error)
	mcpServer       *server.MCPServer
	handler         http.Handler
}

// Config represents the configuration options for GinMCP
type Config struct {
	Name              string
	Description       string
	BaseURL           string
	IncludeOperations []string
	ExcludeOperations []string
	IncludeTags       []string
	ExcludeTags       []string
	Path              string
}

// New creates a new GinMCP instance
func New(engine *gin.Engine, config *Config) *GinMCP {
	if config == nil {
		config = &Config{
			Name:        "Gin MCP",
			Description: "MCP server for Gin application",
		}
	}

	// Create a new MCP server
	s := server.NewMCPServer(
		"Calculator Demo",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	m := &GinMCP{
		engine:            engine,
		name:              config.Name,
		description:       config.Description,
		baseURL:           config.BaseURL,
		operations:        make(map[string]types.Operation),
		config:            config,
		registeredSchemas: make(map[string]types.RegisteredSchemaInfo),
		path:              config.Path,
		mcpServer:         s,
	}

	m.executeToolFunc = m.defaultExecuteTool // Initialize with the default implementation

	// Add debug logging middleware
	if isDebugMode() {
		engine.Use(func(c *gin.Context) {
			start := time.Now()
			path := c.Request.URL.Path

			log.Printf("[HTTP Request] %s %s (Start)", c.Request.Method, path)
			c.Next()
			log.Printf("[HTTP Request] %s %s completed with status %d in %v",
				c.Request.Method, path, c.Writer.Status(), time.Since(start))
		})
	}

	return m
}

// SetExecuteToolFunc allows overriding the default tool execution function.
// This is useful for implementing dynamic baseURL resolution or custom execution logic.
func (m *GinMCP) SetExecuteToolFunc(fn func(ctx context.Context, operationID string, parameters map[string]interface{}, headers http.Header) (interface{}, error)) {
	m.executeToolFunc = fn
}

// RegisterSchema associates Go struct types with a specific route for automatic schema generation.
// Provide nil if a type (Query or Body) is not applicable for the route.
// Example: mcp.RegisterSchema("POST", "/items", nil, main.Item{})
func (m *GinMCP) RegisterSchema(method string, path string, queryType interface{}, bodyType interface{}) {
	m.schemasMu.Lock()
	defer m.schemasMu.Unlock()

	// Ensure method is uppercase for canonical key
	method = strings.ToUpper(method)
	schemaKey := fmt.Sprintf("%s %s", method, path)

	// Validate types slightly (ensure they are structs or pointers to structs if not nil)
	if queryType != nil {
		queryVal := reflect.ValueOf(queryType)
		if queryVal.Kind() == reflect.Ptr {
			queryVal = queryVal.Elem()
		}
		if queryVal.Kind() != reflect.Struct {
			if isDebugMode() {
				log.Printf("Warning: RegisterSchema queryType for %s is not a struct or pointer to struct, reflection might fail.", schemaKey)
			}
		}
	}
	if bodyType != nil {
		bodyVal := reflect.ValueOf(bodyType)
		if bodyVal.Kind() == reflect.Ptr {
			bodyVal = bodyVal.Elem()
		}
		if bodyVal.Kind() != reflect.Struct {
			if isDebugMode() {
				log.Printf("Warning: RegisterSchema bodyType for %s is not a struct or pointer to struct, reflection might fail.", schemaKey)
			}
		}
	}

	m.registeredSchemas[schemaKey] = types.RegisteredSchemaInfo{
		QueryType: queryType,
		BodyType:  bodyType,
	}
	if isDebugMode() {
		log.Printf("Registered schema types for route: %s", schemaKey)
	}
}

// Mount sets up the MCP routes on the given path
func (m *GinMCP) Mount(path string, operationsPath string) {
	content, err := os.ReadFile(path)
	if err != nil {
		zap.L().Error("read mcp config error", zap.Error(err))
		return
	}

	operationsContent, err := os.ReadFile(operationsPath)
	if err != nil {
		zap.L().Error("read mcp config error", zap.Error(err))
		return
	}

	var tools []mcp.Tool = make([]mcp.Tool, 0)
	err = json.Unmarshal(content, &tools)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(operationsContent, &m.operations)
	if err != nil {
		panic(err)
	}

	if len(tools) == 0 {
		return
	}

	for _, tool := range tools {
		m.mcpServer.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Execute the actual Gin endpoint via internal HTTP call
			//execResult, err := m.executeToolFunc(tool.Name, request.Params) // Use the function field
			arguments, ok := request.Params.Arguments.(map[string]interface{})
			if !ok {
				arguments = make(map[string]interface{})
			}
			result, err := m.executeToolFunc(ctx, tool.Name, arguments, request.Header)
			if err != nil {
				message := fmt.Sprintf("Error executing tool '%s': %v", tool.Name, err)
				return mcp.NewToolResultError(message), err
			}
			text, err := json.Marshal(result)
			if err != nil {
				message := fmt.Sprintf("Error executing tool '%s': %v", tool.Name, err)
				return mcp.NewToolResultError(message), err
			}
			return mcp.NewToolResultText(string(text)), nil
		})
	}

	sseServer := server.NewSSEServer(m.mcpServer, server.WithStaticBasePath("/mcp"))

	//httpWrap := &wrappedHTTP{Handler: sseServer}
	//m.handler = sseServer
	//m.engine.Any("/mcp", func(c *gin.Context) {
	//	sseServer.SSEHandler().ServeHTTP(c.Writer, c.Request)
	//})

	m.engine.Any("/mcp", func(c *gin.Context) {
		sseServer.SSEHandler().ServeHTTP(c.Writer, c.Request)
		//httpWrap.ServeHTTP(c.Writer, c.Request)
	})
	m.engine.Any("/mcp/message", func(c *gin.Context) {

		sseServer.MessageHandler().ServeHTTP(c.Writer, c.Request)
		//httpWrap.ServeHTTP(c.Writer, c.Request)
	})
	return
	//if mountPath == "" {
	//	mountPath = "/mcp"
	//}
	//
	//// 1. Setup tools
	//if err := m.SetupServer(); err != nil {
	//	if isDebugMode() {
	//		log.Printf("Failed to setup server: %v", err)
	//	}
	//	return
	//}
	//
	//// 2. Create transport and register handlers
	//m.transport = transport.NewSSETransport(mountPath)
	//m.transport.RegisterHandler("initialize", m.handleInitialize)
	//m.transport.RegisterHandler("tools/list", m.handleToolsList)
	//m.transport.RegisterHandler("tools/call", m.handleToolCall)
	//
	//// 3. Setup CORS middleware
	//m.engine.Use(func(c *gin.Context) {
	//	if isDebugMode() {
	//		log.Printf("[Middleware] Processing request: Method=%s, Path=%s, RemoteAddr=%s", c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr)
	//	}
	//
	//	if strings.HasPrefix(c.Request.URL.Path, mountPath) {
	//		if isDebugMode() {
	//			log.Printf("[Middleware] Path %s matches mountPath %s. Applying headers.", c.Request.URL.Path, mountPath)
	//		}
	//		c.Header("Access-Control-Allow-Origin", "*")
	//		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	//		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Connection-ID")
	//		c.Header("Access-Control-Expose-Headers", "X-Connection-ID")
	//
	//		if c.Request.Method == "OPTIONS" {
	//			if isDebugMode() {
	//				log.Printf("[Middleware] OPTIONS request for %s. Aborting with 204.", c.Request.URL.Path)
	//			}
	//			c.AbortWithStatus(204)
	//			return
	//		} else if c.Request.Method == "POST" {
	//			if isDebugMode() {
	//				log.Printf("[Middleware] POST request for %s. Proceeding to handler.", c.Request.URL.Path)
	//			}
	//		}
	//	} else {
	//		if isDebugMode() {
	//			log.Printf("[Middleware] Path %s does NOT match mountPath %s. Skipping custom logic.", c.Request.URL.Path, mountPath)
	//		}
	//	}
	//	c.Next() // Ensure processing continues
	//	if isDebugMode() {
	//		log.Printf("[Middleware] Finished processing request: Method=%s, Path=%s, Status=%d", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	//	}
	//})
	//
	//// 4. Setup endpoints
	//if isDebugMode() {
	//	log.Printf("[Server Mount DEBUG] Defining GET %s route", mountPath)
	//}
	//m.engine.GET(mountPath, m.handleMCPConnection)
	//if isDebugMode() {
	//	log.Printf("[Server Mount DEBUG] Defining POST %s route", mountPath)
	//}
	//m.engine.POST(mountPath, func(c *gin.Context) {
	//	m.transport.HandleMessage(c)
	//})
}

// handleMCPConnection handles a new MCP connection request
func (m *GinMCP) handleMCPConnection(c *gin.Context) {
	if isDebugMode() {
		log.Println("[Server DEBUG] handleMCPConnection invoked for GET /mcp")
	}
	// 1. Ensure server is ready
	if len(m.tools) == 0 {
		if err := m.SetupServer(); err != nil {
			errID := fmt.Sprintf("err-%d", time.Now().UnixNano())
			c.JSON(http.StatusInternalServerError, &types.MCPMessage{
				Jsonrpc: "2.0",
				ID:      types.RawMessage([]byte(`"` + errID + `"`)),
				Result: map[string]interface{}{
					"code":    "server_error",
					"message": fmt.Sprintf("Failed to setup server: %v", err),
				},
			})
			return
		}
	}

	// 2. Let transport handle the SSE connection
	m.transport.HandleConnection(c)
}

// handleInitialize handles the initialize request from clients
func (m *GinMCP) handleInitialize(msg *types.MCPMessage) *types.MCPMessage {
	// Parse initialization parameters
	params, ok := msg.Params.(map[string]interface{})
	if !ok {
		return &types.MCPMessage{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Error: map[string]interface{}{
				"code":    -32602,
				"message": "Invalid parameters format",
			},
		}
	}

	// Log initialization request
	if isDebugMode() {
		log.Printf("Received initialize request with params: %+v", params)
	}

	// Return server capabilities with correct structure
	return &types.MCPMessage{
		Jsonrpc: "2.0",
		ID:      msg.ID,
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{
					"enabled": true,
					"config": map[string]interface{}{
						"listChanged": false,
					},
				},
				"prompts": map[string]interface{}{
					"enabled": false,
				},
				"resources": map[string]interface{}{
					"enabled": true,
				},
				"logging": map[string]interface{}{
					"enabled": false,
				},
				"roots": map[string]interface{}{
					"listChanged": false,
				},
			},
			"serverInfo": map[string]interface{}{
				"name":       m.name,
				"version":    "2024-11-05",
				"apiVersion": "2024-11-05",
			},
		},
	}
}

// handleToolsList handles the tools/list request
func (m *GinMCP) handleToolsList(msg *types.MCPMessage) *types.MCPMessage {
	// Ensure server is ready
	if err := m.SetupServer(); err != nil {
		return &types.MCPMessage{
			Jsonrpc: "2.0",
			ID:      msg.ID,
			Error: map[string]interface{}{
				"code":    -32603,
				"message": fmt.Sprintf("Failed to setup server: %v", err),
			},
		}
	}

	// Return tools list with proper format
	return &types.MCPMessage{
		Jsonrpc: "2.0",
		ID:      msg.ID,
		Result: map[string]interface{}{
			"tools": m.tools,
			"metadata": map[string]interface{}{
				"version": "2024-11-05",
				"count":   len(m.tools),
			},
		},
	}
}

// SetupServer initializes the MCP server by discovering routes and converting them to tools
func (m *GinMCP) SetupServer() error {
	if len(m.tools) == 0 {
		// Get all routes from the Gin engine
		//routes := m.engine.Routes()

		// Lock schema map while converting
		m.schemasMu.RLock()
		// Convert routes to tools with registered types
		p := convert.New()
		newTools := make([]types.Tool, 0)
		operations := make(map[string]types.Operation)
		p.RegisteredSchemas = m.registeredSchemas
		err := p.ParseAPI(m.path, "main.go", 100)
		if err != nil {
			return err
		}
		newTools = p.NewTools
		operations = p.Operations
		m.schemasMu.RUnlock()

		// Update tools and operations
		m.tools = newTools
		m.operations = operations

		// Filter tools based on configuration (operation/tag filters)
		m.filterTools()

	}

	return nil
}

func (m *GinMCP) GetTools() []types.Tool {
	return m.tools
}

func (m *GinMCP) GetOperations() map[string]types.Operation {
	return m.operations
}

// haveToolsChanged checks if the tools list has changed
func (m *GinMCP) haveToolsChanged(newTools []types.Tool) bool {
	if len(m.tools) != len(newTools) {
		return true
	}

	// Create maps for easier comparison
	oldToolMap := make(map[string]types.Tool)
	for _, tool := range m.tools {
		oldToolMap[tool.Name] = tool
	}

	// Compare tools
	for _, newTool := range newTools {
		oldTool, exists := oldToolMap[newTool.Name]
		if !exists {
			return true
		}
		// Compare tool definitions (you might want to add more detailed comparison)
		if oldTool.Description != newTool.Description {
			return true
		}
	}

	return false
}

// filterTools filters the tools based on configuration
func (m *GinMCP) filterTools() {
	if len(m.tools) == 0 {
		return
	}

	var filteredTools []types.Tool
	config := m.config // Use the GinMCP config

	// Work with local copies to avoid mutating the caller's config
	includeOps := config.IncludeOperations
	includeTags := config.IncludeTags
	excludeOps := config.ExcludeOperations
	excludeTags := config.ExcludeTags

	// Check for conflicting inclusion filters (prefer operations over tags)
	if len(includeOps) > 0 && len(includeTags) > 0 {
		if isDebugMode() {
			log.Printf("Warning: Both IncludeOperations and IncludeTags are set. Preferring IncludeOperations.")
		}
		includeTags = nil
	}

	// Check for conflicting exclusion filters (prefer operations over tags)
	if len(excludeOps) > 0 && len(excludeTags) > 0 {
		if isDebugMode() {
			log.Printf("Warning: Both ExcludeOperations and ExcludeTags are set. Preferring ExcludeOperations.")
		}
		excludeTags = nil
	}

	// Step 1: Apply inclusion filters (operations take precedence over tags)
	if len(includeOps) > 0 {
		includeMap := make(map[string]bool)
		for _, op := range includeOps {
			includeMap[op] = true
		}
		for _, tool := range m.tools {
			if includeMap[tool.Name] {
				filteredTools = append(filteredTools, tool)
			}
		}
		m.tools = filteredTools
		filteredTools = []types.Tool{}
	} else if len(includeTags) > 0 {
		// Include tools that have at least one matching tag
		includeTagsMap := make(map[string]bool)
		for _, tag := range includeTags {
			includeTagsMap[tag] = true
		}
		for _, tool := range m.tools {
			if hasMatchingTag(tool.Tags, includeTagsMap) {
				filteredTools = append(filteredTools, tool)
			}
		}
		m.tools = filteredTools
		filteredTools = []types.Tool{}
	}

	// Step 2: Apply exclusion filters (operations take precedence over tags)
	// Exclusion always wins - it runs on the result of inclusion filtering
	if len(excludeOps) > 0 {
		excludeMap := make(map[string]bool)
		for _, op := range excludeOps {
			excludeMap[op] = true
		}
		for _, tool := range m.tools {
			if !excludeMap[tool.Name] {
				filteredTools = append(filteredTools, tool)
			}
		}
		m.tools = filteredTools
	} else if len(excludeTags) > 0 {
		// Exclude tools that have at least one matching tag
		excludeTagsMap := make(map[string]bool)
		for _, tag := range excludeTags {
			excludeTagsMap[tag] = true
		}
		for _, tool := range m.tools {
			if !hasMatchingTag(tool.Tags, excludeTagsMap) {
				filteredTools = append(filteredTools, tool)
			}
		}
		m.tools = filteredTools
	}
}

// hasMatchingTag checks if any tag in the toolTags slice exists in the filterTags map
func hasMatchingTag(toolTags []string, filterTags map[string]bool) bool {
	for _, tag := range toolTags {
		if filterTags[tag] {
			return true
		}
	}
	return false
}

// defaultExecuteTool is the default implementation for executing a tool.
// It handles the actual invocation of the underlying Gin handler using the configured baseURL.
func (m *GinMCP) defaultExecuteTool(ctx context.Context, operationID string, parameters map[string]interface{}, headers http.Header) (interface{}, error) {
	if isDebugMode() {
		log.Printf("[Tool Execution] Starting execution of tool '%s' with parameters: %+v", operationID, parameters)
	}

	// Find the operation associated with the tool name (operationID)
	operation, ok := m.operations[operationID]
	if !ok {
		if isDebugMode() {
			log.Printf("Error: Operation details not found for tool '%s'", operationID)
		}
		return nil, fmt.Errorf("operation '%s' not found", operationID)
	}
	if isDebugMode() {
		log.Printf("[Tool Execution] Found operation for tool '%s': Method=%s, Path=%s", operationID, operation.Method, operation.Path)
	}

	// Use the configured baseURL (static)
	baseURL := m.baseURL
	if baseURL == "" {
		// Use relative URL if baseURL is not set
		baseURL = ""
		if isDebugMode() {
			log.Printf("[Tool Execution] Using relative URL for request")
		}
	}

	return m.executeToolLogic(ctx, operation, parameters, baseURL, headers)
}

// executeToolLogic contains the core tool execution logic that can be reused
// with different baseURL resolution strategies
func (m *GinMCP) executeToolLogic(ctx context.Context, operation types.Operation, parameters map[string]interface{}, baseURL string, headers http.Header) (interface{}, error) {

	path := operation.Path
	queryParams := url.Values{}
	pathParams := make(map[string]string)

	// Separate args into path params, query params, and body
	for key, value := range parameters {
		// Check against Gin's format ":key"
		placeholder := ":" + key
		if strings.Contains(path, placeholder) {
			// Store the actual value for substitution later
			pathParams[key] = fmt.Sprintf("%v", value)
			if isDebugMode() {
				log.Printf("[Tool Execution] Found path parameter %s=%v", key, value)
			}
		} else {
			// Assume remaining args are query parameters for GET/DELETE
			if operation.Method == "GET" || operation.Method == "DELETE" {
				queryParams.Add(key, fmt.Sprintf("%v", value))
				if isDebugMode() {
					log.Printf("[Tool Execution] Added query parameter %s=%v", key, value)
				}
			}
		}
	}

	// Substitute path parameters using Gin's format ":key"
	for key, value := range pathParams {
		path = strings.Replace(path, ":"+key, value, -1)
	}

	targetURL := baseURL + path
	if len(queryParams) > 0 {
		targetURL += "?" + queryParams.Encode()
	}

	if isDebugMode() {
		log.Printf("[Tool Execution] Making request: %s %s", operation.Method, targetURL)
	}

	// 3. Create and execute the HTTP request
	var reqBody io.Reader
	if operation.Method == "POST" || operation.Method == "PUT" || operation.Method == "PATCH" {
		// For POST/PUT/PATCH, send all non-path args in the body
		bodyData := make(map[string]interface{})
		for key, value := range parameters {
			// Skip ID field for PUT requests since it's in the path
			if key == "id" && operation.Method == "PUT" {
				continue
			}
			if _, isPath := pathParams[key]; !isPath {
				bodyData[key] = value
				if isDebugMode() {
					log.Printf("[Tool Execution] Added body parameter %s=%v", key, value)
				}
			}
		}
		bodyBytes, err := json.Marshal(bodyData)
		if err != nil {
			if isDebugMode() {
				log.Printf("[Tool Execution] Error marshalling request body: %v", err)
			}
			return nil, err
		}
		reqBody = bytes.NewBuffer(bodyBytes)
		if isDebugMode() {
			log.Printf("[Tool Execution] Request body: %s", string(bodyBytes))
		}
	}

	req, err := http.NewRequest(operation.Method, targetURL, reqBody)
	if err != nil {
		if isDebugMode() {
			log.Printf("[Tool Execution] Error creating request: %v", err)
		}
		return nil, err
	}

	// 再次创建一个可取消的 context，嵌套原始的 context
	authorization := headers.Get("Authorization")
	//req.Header.Set("Authorization")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", authorization)
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if isDebugMode() {
		log.Printf("[Tool Execution] Sending request with headers: %+v", req.Header)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		if isDebugMode() {
			log.Printf("[Tool Execution] Error executing request: %v", err)
		}
		return nil, err
	}
	defer resp.Body.Close()

	// 4. Read and parse the response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if isDebugMode() {
			log.Printf("[Tool Execution] Error reading response body: %v", err)
		}
		return nil, err
	}

	if isDebugMode() {
		log.Printf("[Tool Execution] Response status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if isDebugMode() {
			log.Printf("[Tool Execution] Request failed with status %d", resp.StatusCode)
		}
		// Attempt to return the error body, otherwise just the status
		var errorData interface{}
		if json.Unmarshal(bodyBytes, &errorData) == nil {
			return nil, fmt.Errorf("request failed with status %d: %v", resp.StatusCode, errorData)
		}
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	var resultData interface{}
	if err := json.Unmarshal(bodyBytes, &resultData); err != nil {
		if isDebugMode() {
			log.Printf("[Tool Execution] Error unmarshalling response: %v", err)
		}
		// Return raw body if unmarshalling fails but status was ok
		return string(bodyBytes), nil
	}

	if isDebugMode() {
		log.Printf("[Tool Execution] Successfully completed tool execution")
	}

	return resultData, nil
}

// BaseURLResolver is a function type for resolving baseURL dynamically
type BaseURLResolver func() string

// NewEnvironmentResolver creates a resolver that extracts baseURL from environment variables.
// Useful for Quicknode scenarios where the user endpoint is set via environment.
func NewEnvironmentResolver(envVarName string, fallback string) BaseURLResolver {
	return func() string {
		if value := os.Getenv(envVarName); value != "" {
			return value
		}
		return fallback
	}
}

// NewHeaderResolver creates a resolver that extracts baseURL from HTTP headers.
// This requires access to the current request context.
// Note: This is a template - you'll need to adapt this to your request handling pattern.
func NewHeaderResolver(headerName string, fallback string) BaseURLResolver {
	return func() string {
		// In a real implementation, you would need access to the current request context
		// This could be achieved through:
		// 1. Thread-local storage
		// 2. Context.Context passed through the call chain
		// 3. Middleware that sets a global variable

		// For now, return fallback - see example usage for complete implementation
		return fallback
	}
}

// NewQuicknodeResolver creates a resolver specifically for Quicknode environments.
// It tries multiple common patterns for extracting the user endpoint.
func NewQuicknodeResolver(fallback string) BaseURLResolver {
	return func() string {
		// Try Quicknode-specific environment variables
		if endpoint := os.Getenv("QUICKNODE_USER_ENDPOINT"); endpoint != "" {
			return endpoint
		}
		if endpoint := os.Getenv("USER_ENDPOINT"); endpoint != "" {
			return endpoint
		}
		if host := os.Getenv("HOST"); host != "" {
			if !strings.HasPrefix(host, "http") {
				return "https://" + host
			}
			return host
		}

		return fallback
	}
}

// NewRAGFlowResolver creates a resolver specifically for RAGFlow environments.
// It tries multiple common patterns for extracting the RAGFlow endpoint.
func NewRAGFlowResolver(fallback string) BaseURLResolver {
	return func() string {
		// Try RAGFlow-specific environment variables in order of preference
		if endpoint := os.Getenv("RAGFLOW_ENDPOINT"); endpoint != "" {
			return endpoint
		}
		if workflowURL := os.Getenv("RAGFLOW_WORKFLOW_URL"); workflowURL != "" {
			return workflowURL
		}

		// Try building from base URL and workflow ID
		baseURL := os.Getenv("RAGFLOW_BASE_URL")
		workflowID := os.Getenv("WORKFLOW_ID")
		if baseURL != "" && workflowID != "" {
			baseURL = strings.TrimSuffix(baseURL, "/")
			return baseURL + "/workflow/" + workflowID
		}

		// Try just base URL
		if baseURL != "" {
			return baseURL
		}

		// Try generic HOST variable
		if host := os.Getenv("HOST"); host != "" {
			if !strings.HasPrefix(host, "http") {
				return "https://" + host
			}
			return host
		}

		return fallback
	}
}
