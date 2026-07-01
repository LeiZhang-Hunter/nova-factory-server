package gatewayserviceimpl

import (
	"testing"

	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

type mockAgentConfigKeyDao struct {
	current            *gatewaymodels.AgentConfigKey
	updateAllowID      int64
	updateAllowTools   string
	updateAllowCalled  bool
	updateAllowResult  *gatewaymodels.AgentConfigKey
	getByIDCalledCount int
}

func (m *mockAgentConfigKeyDao) Create(c *gin.Context, req *gatewaymodels.AgentConfigKeyUpsert) (*gatewaymodels.AgentConfigKey, error) {
	return nil, nil
}

func (m *mockAgentConfigKeyDao) Update(c *gin.Context, req *gatewaymodels.AgentConfigKeyUpsert) (*gatewaymodels.AgentConfigKey, error) {
	return nil, nil
}

func (m *mockAgentConfigKeyDao) UpdateAllowMcpTools(c *gin.Context, id int64, tools string) (*gatewaymodels.AgentConfigKey, error) {
	m.updateAllowCalled = true
	m.updateAllowID = id
	m.updateAllowTools = tools
	if m.updateAllowResult != nil {
		return m.updateAllowResult, nil
	}
	return m.current, nil
}

func (m *mockAgentConfigKeyDao) DeleteByIDs(c *gin.Context, ids []int64) error {
	return nil
}

func (m *mockAgentConfigKeyDao) GetByID(c *gin.Context, id int64) (*gatewaymodels.AgentConfigKey, error) {
	m.getByIDCalledCount++
	return m.current, nil
}

func (m *mockAgentConfigKeyDao) GetByKey(c *gin.Context, key string) (*gatewaymodels.AgentConfigKey, error) {
	return nil, nil
}

func (m *mockAgentConfigKeyDao) List(c *gin.Context, req *gatewaymodels.AgentConfigKeyQuery) (*gatewaymodels.AgentConfigKeyListData, error) {
	return nil, nil
}

func TestAgentConfigKeySetAllowMcpToolsNormalizesTools(t *testing.T) {
	dao := &mockAgentConfigKeyDao{
		current:           &gatewaymodels.AgentConfigKey{ID: 1},
		updateAllowResult: &gatewaymodels.AgentConfigKey{ID: 1},
	}
	service := &AgentConfigKeyServiceImpl{dao: dao}

	data, err := service.SetAllowMcpTools(&gin.Context{}, &gatewaymodels.AgentConfigKeyToolUpsert{
		ID:    1,
		Tools: []string{" tool.a ", "", "tool.b", "tool.a"},
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if data == nil || data.ID != 1 {
		t.Fatalf("unexpected result: %#v", data)
	}
	if !dao.updateAllowCalled {
		t.Fatalf("expected UpdateAllowMcpTools to be called")
	}
	if dao.updateAllowID != 1 {
		t.Fatalf("unexpected update id: %d", dao.updateAllowID)
	}
	if dao.updateAllowTools != `["tool.a","tool.b"]` {
		t.Fatalf("unexpected tools payload: %s", dao.updateAllowTools)
	}
}

func TestAgentConfigKeySetAllowMcpToolsClearsTools(t *testing.T) {
	dao := &mockAgentConfigKeyDao{
		current:           &gatewaymodels.AgentConfigKey{ID: 2},
		updateAllowResult: &gatewaymodels.AgentConfigKey{ID: 2},
	}
	service := &AgentConfigKeyServiceImpl{dao: dao}

	_, err := service.SetAllowMcpTools(&gin.Context{}, &gatewaymodels.AgentConfigKeyToolUpsert{
		ID:    2,
		Tools: []string{" ", ""},
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if dao.updateAllowTools != "" {
		t.Fatalf("expected tools payload to be empty, got: %q", dao.updateAllowTools)
	}
}
