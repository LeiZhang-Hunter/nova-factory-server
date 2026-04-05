package gatewayserviceimpl

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"

	"github.com/gin-gonic/gin"
)

const (
	mcpTransportStdio          = "stdio"
	mcpTransportStreamableHTTP = "streamableHttp"
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
	return m.dao.Update(c, req)
}

func (m *MCPServerServiceImpl) DeleteByIDs(c *gin.Context, ids []string) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的MCP服务")
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

func (m *MCPServerServiceImpl) prepareUpsert(req *gatewaymodels.MCPServerUpsert, isUpdate bool) error {
	if req == nil {
		return errors.New("参数不能为空")
	}
	if isUpdate && strings.TrimSpace(req.ID) == "" {
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
	if req.Transport != mcpTransportStdio && req.Transport != mcpTransportStreamableHTTP {
		return errors.New("传输方式仅支持 stdio 或 streamableHttp")
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
