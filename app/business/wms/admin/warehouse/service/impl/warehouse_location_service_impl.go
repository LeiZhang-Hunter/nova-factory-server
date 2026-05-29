package impl

import (
	"errors"
	"strings"

	"nova-factory-server/app/business/wms/admin/warehouse/dao"
	"nova-factory-server/app/business/wms/admin/warehouse/models"
	"nova-factory-server/app/business/wms/admin/warehouse/service"

	"github.com/gin-gonic/gin"
)

// WarehouseLocationServiceImpl 提供 WMS 库位相关业务能力。
type WarehouseLocationServiceImpl struct {
	dao dao.IWarehouseLocationDao
}

// NewWarehouseLocationService 创建 WMS 库位服务。
func NewWarehouseLocationService(dao dao.IWarehouseLocationDao) service.IWarehouseLocationService {
	return &WarehouseLocationServiceImpl{dao: dao}
}

// Set 新增或修改 WMS 库位。
func (s *WarehouseLocationServiceImpl) Set(c *gin.Context, req *models.WarehouseLocationSet) (*models.WarehouseLocation, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	trimWarehouseLocationSet(req)
	if req.WarehouseID <= 0 {
		return nil, errors.New("仓库ID不能为空")
	}
	if req.WarehouseName == "" {
		return nil, errors.New("仓库名称不能为空")
	}
	if req.WarehouseAreaID <= 0 {
		return nil, errors.New("库区ID不能为空")
	}
	if req.WarehouseAreaName == "" {
		return nil, errors.New("库区名称不能为空")
	}
	if req.LocationName == "" {
		return nil, errors.New("库位名称不能为空")
	}
	if err := validateWarehouseLocationNumbers(req); err != nil {
		return nil, err
	}
	exists, err := s.dao.GetByName(c, req.WarehouseAreaID, req.LocationName)
	if err != nil {
		return nil, err
	}
	if exists != nil && exists.ID != req.ID {
		return nil, errors.New("同一库区下的库位名称已存在")
	}
	if req.TagNumber != "" {
		tagItem, err := s.dao.GetByTagNumber(c, req.TagNumber)
		if err != nil {
			return nil, err
		}
		if tagItem != nil && tagItem.ID != req.ID {
			return nil, errors.New("标签号已存在")
		}
	}
	return s.dao.Set(c, req)
}

// DeleteByIDs 删除 WMS 库位。
func (s *WarehouseLocationServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

// GetByID 查询 WMS 库位详情。
func (s *WarehouseLocationServiceImpl) GetByID(c *gin.Context, id int64) (*models.WarehouseLocation, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

// List 查询 WMS 库位列表。
func (s *WarehouseLocationServiceImpl) List(c *gin.Context, req *models.WarehouseLocationQuery) (*models.WarehouseLocationListData, error) {
	if req != nil {
		req.LocationName = strings.TrimSpace(req.LocationName)
		req.RoadwayNumber = strings.TrimSpace(req.RoadwayNumber)
		req.ShelfNumber = strings.TrimSpace(req.ShelfNumber)
		req.LayerNumber = strings.TrimSpace(req.LayerNumber)
		req.TagNumber = strings.TrimSpace(req.TagNumber)
	}
	return s.dao.List(c, req)
}

func trimWarehouseLocationSet(req *models.WarehouseLocationSet) {
	req.WarehouseName = strings.TrimSpace(req.WarehouseName)
	req.WarehouseAreaName = strings.TrimSpace(req.WarehouseAreaName)
	req.LocationName = strings.TrimSpace(req.LocationName)
	req.RoadwayNumber = strings.TrimSpace(req.RoadwayNumber)
	req.ShelfNumber = strings.TrimSpace(req.ShelfNumber)
	req.LayerNumber = strings.TrimSpace(req.LayerNumber)
	req.TagNumber = strings.TrimSpace(req.TagNumber)
}

func validateWarehouseLocationNumbers(req *models.WarehouseLocationSet) error {
	if req.LocationLength < 0 {
		return errors.New("库位长度不能小于0")
	}
	if req.LocationWidth < 0 {
		return errors.New("库位宽度不能小于0")
	}
	if req.LocationHeight < 0 {
		return errors.New("库位高度不能小于0")
	}
	if req.LocationVolume < 0 {
		return errors.New("库位容积不能小于0")
	}
	if req.LocationLoad < 0 {
		return errors.New("库位承重不能小于0")
	}
	return nil
}
