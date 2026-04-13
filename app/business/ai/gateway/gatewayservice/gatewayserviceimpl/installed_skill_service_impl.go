package gatewayserviceimpl

import (
	"errors"
	"strings"

	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"

	"github.com/gin-gonic/gin"
)

// InstalledSkillServiceImpl 提供已安装技能的业务实现。
type InstalledSkillServiceImpl struct {
	dao gatewaydao.IInstalledSkillDao
}

// NewInstalledSkillService 创建已安装技能服务。
func NewInstalledSkillService(dao gatewaydao.IInstalledSkillDao) gatewayservice.IInstalledSkillService {
	return &InstalledSkillServiceImpl{dao: dao}
}

// Create 新增已安装技能。
func (i *InstalledSkillServiceImpl) Create(c *gin.Context, req *gatewaymodels.InstalledSkillUpsert) (*gatewaymodels.InstalledSkill, error) {
	if err := i.validateUpsert(c, req, false); err != nil {
		return nil, err
	}
	return i.dao.Create(c, req)
}

// Update 修改已安装技能。
func (i *InstalledSkillServiceImpl) Update(c *gin.Context, req *gatewaymodels.InstalledSkillUpsert) (*gatewaymodels.InstalledSkill, error) {
	if req.ID == 0 {
		return nil, errors.New("id不能为空")
	}
	if err := i.validateUpsert(c, req, true); err != nil {
		return nil, err
	}
	return i.dao.Update(c, req)
}

// DeleteByIDs 删除已安装技能。
func (i *InstalledSkillServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的技能")
	}
	return i.dao.DeleteByIDs(c, ids)
}

// GetByID 查询已安装技能详情。
func (i *InstalledSkillServiceImpl) GetByID(c *gin.Context, id int64) (*gatewaymodels.InstalledSkill, error) {
	if id == 0 {
		return nil, errors.New("id不能为空")
	}
	return i.dao.GetByID(c, id)
}

// List 查询已安装技能列表。
func (i *InstalledSkillServiceImpl) List(c *gin.Context, req *gatewaymodels.InstalledSkillQuery) (*gatewaymodels.InstalledSkillListData, error) {
	if req != nil {
		req.Name = strings.TrimSpace(req.Name)
		req.Slug = strings.TrimSpace(req.Slug)
	}
	return i.dao.List(c, req)
}

// validateUpsert 校验已安装技能新增修改参数。
func (i *InstalledSkillServiceImpl) validateUpsert(c *gin.Context, req *gatewaymodels.InstalledSkillUpsert, isUpdate bool) error {
	if req == nil {
		return errors.New("参数不能为空")
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Slug = strings.TrimSpace(req.Slug)
	req.Version = strings.TrimSpace(req.Version)
	req.Source = strings.TrimSpace(req.Source)
	req.Description = strings.TrimSpace(req.Description)
	if req.Name == "" {
		return errors.New("技能名称不能为空")
	}
	if req.Slug == "" {
		return errors.New("技能标识不能为空")
	}
	if req.Source == "" {
		return errors.New("技能来源不能为空")
	}
	if req.Description == "" {
		return errors.New("技能描述不能为空")
	}
	if req.Enabled == nil {
		req.Enabled = installedSkillBoolPtr(true)
	}
	exists, err := i.dao.GetBySlug(c, req.Slug)
	if err != nil {
		return err
	}
	if exists == nil {
		return nil
	}
	if isUpdate && exists.ID == req.ID {
		return nil
	}
	return errors.New("技能标识已存在")
}

func installedSkillBoolPtr(v bool) *bool {
	return &v
}
