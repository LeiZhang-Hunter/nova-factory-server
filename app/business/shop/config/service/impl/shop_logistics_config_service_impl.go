package impl

import (
	"errors"
	"strings"

	"nova-factory-server/app/business/shop/config/dao"
	"nova-factory-server/app/business/shop/config/models"
	"nova-factory-server/app/business/shop/config/service"

	"github.com/gin-gonic/gin"
)

// ShopLogisticsConfigServiceImpl 物流配置服务实现
type ShopLogisticsConfigServiceImpl struct {
	dao dao.IShopLogisticsConfigDao
}

// NewShopLogisticsConfigService 创建物流配置服务
func NewShopLogisticsConfigService(dao dao.IShopLogisticsConfigDao) service.IShopLogisticsConfigService {
	return &ShopLogisticsConfigServiceImpl{dao: dao}
}

func (s *ShopLogisticsConfigServiceImpl) Create(c *gin.Context, req *models.ShopLogisticsConfigSet) (*models.ShopLogisticsConfig, error) {
	if err := s.validateUnique(c, req); err != nil {
		return nil, err
	}
	return s.dao.Create(c, req)
}

func (s *ShopLogisticsConfigServiceImpl) Update(c *gin.Context, req *models.ShopLogisticsConfigSet) (*models.ShopLogisticsConfig, error) {
	if err := s.validateUnique(c, req); err != nil {
		return nil, err
	}
	return s.dao.Update(c, req)
}

func (s *ShopLogisticsConfigServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ShopLogisticsConfigServiceImpl) GetByID(c *gin.Context, id int64) (*models.ShopLogisticsConfig, error) {
	return s.dao.GetByID(c, id)
}

func (s *ShopLogisticsConfigServiceImpl) GetByType(c *gin.Context, typ string) (*models.ShopLogisticsConfig, error) {
	return s.dao.GetByType(c, typ)
}

func (s *ShopLogisticsConfigServiceImpl) List(c *gin.Context, req *models.ShopLogisticsConfigQuery) (*models.ShopLogisticsConfigListData, error) {
	return s.dao.List(c, req)
}

func (s *ShopLogisticsConfigServiceImpl) validateUnique(c *gin.Context, req *models.ShopLogisticsConfigSet) error {
	req.Type = strings.TrimSpace(req.Type)
	if req.Type == "" {
		return errors.New("配置类型不能为空")
	}
	exist, err := s.dao.GetByType(c, req.Type)
	if err != nil {
		return err
	}
	if exist != nil && exist.ID != req.ID {
		return errors.New("配置类型已存在")
	}
	return nil
}
