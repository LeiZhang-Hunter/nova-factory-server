package gatewayserviceimpl

import (
	"testing"

	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

func TestInstalledSkillServiceCreateFillDefaultEnabled(t *testing.T) {
	dao := &mockInstalledSkillDao{}
	service := &InstalledSkillServiceImpl{dao: dao}
	req := &gatewaymodels.InstalledSkillUpsert{
		Name:        "demo",
		Slug:        "demo-skill",
		Source:      "local",
		Description: "demo desc",
	}

	_, err := service.Create(&gin.Context{}, req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if req.Enabled == nil || !*req.Enabled {
		t.Fatal("expected enabled default true")
	}
}

func TestInstalledSkillServiceRejectDuplicateSlug(t *testing.T) {
	dao := &mockInstalledSkillDao{
		bySlug: &gatewaymodels.InstalledSkill{ID: 1, Slug: "demo-skill"},
	}
	service := &InstalledSkillServiceImpl{dao: dao}
	req := &gatewaymodels.InstalledSkillUpsert{
		Name:        "demo",
		Slug:        "demo-skill",
		Source:      "local",
		Description: "demo desc",
	}

	_, err := service.Create(&gin.Context{}, req)
	if err == nil || err.Error() != "技能标识已存在" {
		t.Fatalf("unexpected err: %v", err)
	}
}

type mockInstalledSkillDao struct {
	bySlug *gatewaymodels.InstalledSkill
}

func (m *mockInstalledSkillDao) Create(c *gin.Context, req *gatewaymodels.InstalledSkillUpsert) (*gatewaymodels.InstalledSkill, error) {
	return &gatewaymodels.InstalledSkill{
		ID:      1,
		Name:    req.Name,
		Slug:    req.Slug,
		Enabled: req.Enabled,
	}, nil
}

func (m *mockInstalledSkillDao) Update(c *gin.Context, req *gatewaymodels.InstalledSkillUpsert) (*gatewaymodels.InstalledSkill, error) {
	return &gatewaymodels.InstalledSkill{
		ID:      req.ID,
		Name:    req.Name,
		Slug:    req.Slug,
		Enabled: req.Enabled,
	}, nil
}

func (m *mockInstalledSkillDao) DeleteByIDs(c *gin.Context, ids []int64) error {
	return nil
}

func (m *mockInstalledSkillDao) GetByID(c *gin.Context, id int64) (*gatewaymodels.InstalledSkill, error) {
	return nil, nil
}

func (m *mockInstalledSkillDao) GetBySlug(c *gin.Context, slug string) (*gatewaymodels.InstalledSkill, error) {
	return m.bySlug, nil
}

func (m *mockInstalledSkillDao) List(c *gin.Context, req *gatewaymodels.InstalledSkillQuery) (*gatewaymodels.InstalledSkillListData, error) {
	return nil, nil
}
