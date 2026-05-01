package impl

import (
	"fmt"
	"strings"

	"nova-factory-server/app/business/shop/home/dao"
	"nova-factory-server/app/business/shop/home/models"
	"nova-factory-server/app/business/shop/home/service"

	"github.com/gin-gonic/gin"
)

// ShopHomeModuleServiceImpl 提供首页模块相关业务能力。
type ShopHomeModuleServiceImpl struct {
	dao dao.IShopHomeModuleDao
}

// NewShopHomeModuleService 创建首页模块服务。
func NewShopHomeModuleService(dao dao.IShopHomeModuleDao) service.IShopHomeModuleService {
	return &ShopHomeModuleServiceImpl{dao: dao}
}

// Set 新增或修改首页模块。
func (s *ShopHomeModuleServiceImpl) Set(c *gin.Context, req *models.HomeModuleSet) (*models.HomeModule, error) {
	if req == nil {
		return nil, fmt.Errorf("参数不能为空")
	}
	req.PageKey = strings.TrimSpace(req.PageKey)
	req.ModuleType = strings.TrimSpace(req.ModuleType)
	req.ModuleName = strings.TrimSpace(req.ModuleName)
	req.Title = strings.TrimSpace(req.Title)
	req.SubTitle = strings.TrimSpace(req.SubTitle)
	req.MoreLink = strings.TrimSpace(req.MoreLink)
	req.StyleJSON = strings.TrimSpace(req.StyleJSON)
	req.RuleJSON = strings.TrimSpace(req.RuleJSON)
	req.ExtJSON = strings.TrimSpace(req.ExtJSON)

	if req.PageKey == "" {
		req.PageKey = "index"
	}
	if req.ModuleType == "" {
		return nil, fmt.Errorf("模块类型不能为空")
	}
	if req.ModuleName == "" {
		return nil, fmt.Errorf("模块名称不能为空")
	}
	if req.Sort < 0 {
		return nil, fmt.Errorf("排序不能小于0")
	}
	if req.LimitNum < 0 {
		return nil, fmt.Errorf("展示数量不能小于0")
	}
	if req.StartTime > 0 && req.EndTime > 0 && req.StartTime > req.EndTime {
		return nil, fmt.Errorf("开始时间不能大于结束时间")
	}
	return s.dao.Set(c, req)
}

// DeleteByIDs 删除首页模块。
func (s *ShopHomeModuleServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

// GetByID 根据主键获取首页模块详情。
func (s *ShopHomeModuleServiceImpl) GetByID(c *gin.Context, id int64) (*models.HomeModule, error) {
	return s.dao.GetByID(c, id)
}

// List 查询首页模块列表。
func (s *ShopHomeModuleServiceImpl) List(c *gin.Context, req *models.HomeModuleQuery) (*models.HomeModuleListData, error) {
	return s.dao.List(c, req)
}
