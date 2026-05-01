package impl

import (
	"fmt"
	"strings"

	"nova-factory-server/app/business/shop/home/dao"
	"nova-factory-server/app/business/shop/home/models"
	"nova-factory-server/app/business/shop/home/service"

	"github.com/gin-gonic/gin"
)

// ShopHomeModuleItemServiceImpl 提供首页模块明细相关业务能力。
type ShopHomeModuleItemServiceImpl struct {
	dao       dao.IShopHomeModuleItemDao
	moduleDao dao.IShopHomeModuleDao
}

// NewShopHomeModuleItemService 创建首页模块明细服务。
func NewShopHomeModuleItemService(dao dao.IShopHomeModuleItemDao, moduleDao dao.IShopHomeModuleDao) service.IShopHomeModuleItemService {
	return &ShopHomeModuleItemServiceImpl{
		dao:       dao,
		moduleDao: moduleDao,
	}
}

// Set 新增或修改首页模块明细。
func (s *ShopHomeModuleItemServiceImpl) Set(c *gin.Context, req *models.HomeModuleItemSet) (*models.HomeModuleItem, error) {
	if req == nil {
		return nil, fmt.Errorf("参数不能为空")
	}
	req.BusinessType = strings.TrimSpace(req.BusinessType)
	req.ItemName = strings.TrimSpace(req.ItemName)
	req.ItemSubTitle = strings.TrimSpace(req.ItemSubTitle)
	req.ItemImage = strings.TrimSpace(req.ItemImage)
	req.ExtJSON = strings.TrimSpace(req.ExtJSON)

	if req.ModuleID <= 0 {
		return nil, fmt.Errorf("模块ID不能为空")
	}
	module, err := s.moduleDao.GetByID(c, req.ModuleID)
	if err != nil {
		return nil, err
	}
	if module == nil {
		return nil, fmt.Errorf("首页模块不存在")
	}
	if req.BusinessType == "" {
		return nil, fmt.Errorf("业务类型不能为空")
	}
	if req.LinkID <= 0 {
		return nil, fmt.Errorf("关联ID不能为空")
	}
	if req.Sort < 0 {
		return nil, fmt.Errorf("排序不能小于0")
	}
	return s.dao.Set(c, req)
}

// DeleteByIDs 删除首页模块明细。
func (s *ShopHomeModuleItemServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

// GetByID 根据主键获取首页模块明细详情。
func (s *ShopHomeModuleItemServiceImpl) GetByID(c *gin.Context, id int64) (*models.HomeModuleItem, error) {
	return s.dao.GetByID(c, id)
}

// List 查询首页模块明细列表。
func (s *ShopHomeModuleItemServiceImpl) List(c *gin.Context, req *models.HomeModuleItemQuery) (*models.HomeModuleItemListData, error) {
	return s.dao.List(c, req)
}
