package impl

import (
	"errors"
	"strings"

	"nova-factory-server/app/business/wms/admin/warehouse/dao"
	"nova-factory-server/app/business/wms/admin/warehouse/models"
	"nova-factory-server/app/business/wms/admin/warehouse/service"

	"github.com/gin-gonic/gin"
)

// WarehouseAreaServiceImpl 提供 WMS 库区相关业务能力。
type WarehouseAreaServiceImpl struct {
	dao dao.IWarehouseAreaDao
}

// NewWarehouseAreaService 创建 WMS 库区服务。
func NewWarehouseAreaService(dao dao.IWarehouseAreaDao) service.IWarehouseAreaService {
	return &WarehouseAreaServiceImpl{dao: dao}
}

// Set 新增或修改 WMS 库区。
func (s *WarehouseAreaServiceImpl) Set(c *gin.Context, req *models.WarehouseAreaSet) (*models.WarehouseArea, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	req.AreaName = strings.TrimSpace(req.AreaName)
	if req.WarehouseID <= 0 {
		return nil, errors.New("仓库ID不能为空")
	}
	if req.AreaName == "" {
		return nil, errors.New("库区名称不能为空")
	}
	exists, err := s.dao.GetByName(c, req.WarehouseID, req.ParentID, req.AreaName)
	if err != nil {
		return nil, err
	}
	if exists != nil && exists.ID != req.ID {
		return nil, errors.New("同一父级下的库区名称已存在")
	}
	return s.dao.Set(c, req)
}

// DeleteByIDs 删除 WMS 库区。
func (s *WarehouseAreaServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

// GetByID 查询 WMS 库区详情。
func (s *WarehouseAreaServiceImpl) GetByID(c *gin.Context, id int64) (*models.WarehouseArea, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

// List 查询 WMS 库区列表。
func (s *WarehouseAreaServiceImpl) List(c *gin.Context, req *models.WarehouseAreaQuery) (*models.WarehouseAreaListData, error) {
	if req != nil {
		req.AreaName = strings.TrimSpace(req.AreaName)
	}
	return s.dao.List(c, req)
}
