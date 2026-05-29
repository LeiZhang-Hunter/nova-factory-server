package impl

import (
	"errors"
	"strings"

	mastermodels "nova-factory-server/app/business/erp/master/mastermodels"
	masterservice "nova-factory-server/app/business/erp/master/masterservice"
	"nova-factory-server/app/business/wms/admin/warehouse/dao"
	"nova-factory-server/app/business/wms/admin/warehouse/models"
	wmsservice "nova-factory-server/app/business/wms/admin/warehouse/service"

	"github.com/gin-gonic/gin"
)

// WarehouseAreaServiceImpl 提供 WMS 库区相关业务能力。
type WarehouseAreaServiceImpl struct {
	dao              dao.IWarehouseAreaDao
	warehouseService masterservice.IWarehouseService
}

// NewWarehouseAreaService 创建 WMS 库区服务。
func NewWarehouseAreaService(dao dao.IWarehouseAreaDao, warehouseService masterservice.IWarehouseService) wmsservice.IWarehouseAreaService {
	return &WarehouseAreaServiceImpl{dao: dao, warehouseService: warehouseService}
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
	data, err := s.dao.List(c, req)
	if err != nil {
		return nil, err
	}
	if err := s.fillWarehouseAreaNames(c, data.Rows); err != nil {
		return nil, err
	}
	return data, nil
}

// fillWarehouseAreaNames 填充库区
func (s *WarehouseAreaServiceImpl) fillWarehouseAreaNames(c *gin.Context, rows []*models.WarehouseArea) error {
	if len(rows) == 0 {
		return nil
	}
	warehouseMap, err := s.getWarehouseMap(c, collectWarehouseIDs(rows))
	if err != nil {
		return err
	}
	parentAreaMap, err := s.getParentAreaMap(c, collectParentAreaIDs(rows))
	if err != nil {
		return err
	}
	for _, row := range rows {
		if row == nil {
			continue
		}
		if warehouse, ok := warehouseMap[row.WarehouseID]; ok && warehouse != nil {
			row.WarehouseName = warehouse.Name
		}
		if parentArea, ok := parentAreaMap[row.ParentID]; ok && parentArea != nil {
			row.ParentAreaName = parentArea.AreaName
		}
	}
	return nil
}

func (s *WarehouseAreaServiceImpl) getWarehouseMap(c *gin.Context, ids []int64) (map[int64]*mastermodels.Warehouse, error) {
	result := make(map[int64]*mastermodels.Warehouse, len(ids))
	if len(ids) == 0 || s.warehouseService == nil {
		return result, nil
	}
	rows, err := s.warehouseService.GetByIDs(c, ids)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if row == nil {
			continue
		}
		result[row.ID] = row
	}
	return result, nil
}

func (s *WarehouseAreaServiceImpl) getParentAreaMap(c *gin.Context, ids []int64) (map[int64]*models.WarehouseArea, error) {
	result := make(map[int64]*models.WarehouseArea, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	rows, err := s.dao.GetByIDs(c, ids)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if row == nil {
			continue
		}
		result[row.ID] = row
	}
	return result, nil
}

func collectWarehouseIDs(rows []*models.WarehouseArea) []int64 {
	seen := make(map[int64]struct{})
	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.WarehouseID <= 0 {
			continue
		}
		if _, ok := seen[row.WarehouseID]; ok {
			continue
		}
		seen[row.WarehouseID] = struct{}{}
		ids = append(ids, row.WarehouseID)
	}
	return ids
}

func collectParentAreaIDs(rows []*models.WarehouseArea) []int64 {
	seen := make(map[int64]struct{})
	ids := make([]int64, 0, len(rows))
	for _, row := range rows {
		if row == nil || row.ParentID <= 0 {
			continue
		}
		if _, ok := seen[row.ParentID]; ok {
			continue
		}
		seen[row.ParentID] = struct{}{}
		ids = append(ids, row.ParentID)
	}
	return ids
}
