package gatewayserviceimpl

import (
	"testing"

	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

func TestAIAgentServiceCreateFillDefaults(t *testing.T) {
	dao := &mockAIAgentDao{}
	service := &AIAgentServiceImpl{dao: dao}
	req := &gatewaymodels.AIAgentUpsert{
		Name:                "demo-agent",
		MCPServerIDs:        `["mcp-1"]`,
		MCPServerEnabledIDs: `["mcp-1"]`,
	}

	_, err := service.Create(&gin.Context{}, req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if req.MCPEnabled == nil || *req.MCPEnabled {
		t.Fatal("expected mcpEnabled default false")
	}
	if req.EnableLLMTemperature == nil || *req.EnableLLMTemperature {
		t.Fatal("expected enableLlmTemperature default false")
	}
}

type mockAIAgentDao struct{}

func (m *mockAIAgentDao) Create(c *gin.Context, req *gatewaymodels.AIAgentUpsert) (*gatewaymodels.AIAgent, error) {
	return &gatewaymodels.AIAgent{Name: req.Name}, nil
}

func (m *mockAIAgentDao) Update(c *gin.Context, req *gatewaymodels.AIAgentUpsert) (*gatewaymodels.AIAgent, error) {
	return &gatewaymodels.AIAgent{ID: req.ID, Name: req.Name}, nil
}

func (m *mockAIAgentDao) DeleteByIDs(c *gin.Context, ids []int64) error {
	return nil
}

func (m *mockAIAgentDao) GetByID(c *gin.Context, id int64) (*gatewaymodels.AIAgent, error) {
	return &gatewaymodels.AIAgent{ID: id, Name: "demo-agent"}, nil
}

func (m *mockAIAgentDao) List(c *gin.Context, req *gatewaymodels.AIAgentQuery) (*gatewaymodels.AIAgentListData, error) {
	return nil, nil
}
