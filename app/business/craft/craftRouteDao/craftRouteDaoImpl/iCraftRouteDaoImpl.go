package craftRouteDaoImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type CraftRouteDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewCraftRouteDaoImpl(db *gorm.DB) craftRouteDao.ICraftRouteDao {
	return &CraftRouteDaoImpl{
		db:        db,
		tableName: "sys_craft_route",
	}
}

func (dao *CraftRouteDaoImpl) AddCraftRoute(c *gin.Context, route *craftRouteModels.SysCraftRouteRequest) (*craftRouteModels.SysCraftRoute, error) {
	value := &craftRouteModels.SysCraftRoute{
		RouteID:   snowflake.GenID(),
		RouteCode: route.RouteCode,
		RouteName: route.RouteName,
		RouteDesc: route.RouteDesc,
		Remark:    route.Remark,
		Status:    false,
		DeptID:    baizeContext.GetDeptId(c),
	}
	value.SetCreateBy(baizeContext.GetUserId(c))
	ret := dao.db.Table(dao.tableName).Create(value)
	return value, ret.Error
}

func (dao *CraftRouteDaoImpl) UpdateCraftRoute(c *gin.Context, route *craftRouteModels.SysCraftRouteRequest) (*craftRouteModels.SysCraftRoute, error) {
	if route.RouteID == 0 {
		return nil, errors.New("route.RouteID == 0")
	}
	var info *craftRouteModels.SysCraftRoute
	ret := dao.db.Table(dao.tableName).Where("id=?", route.RouteID).First(&info)
	if ret.Error != nil {
		zap.L().Error("读取工艺流程图错误", zap.Error(ret.Error))
		return nil, errors.New("工艺流程图不存在")
	}
	ret = dao.db.Table(dao.tableName).Where("id=?", route.RouteID).Updates(&info)
	return info, ret.Error
}

func (dao *CraftRouteDaoImpl) RemoveCraftRoute(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选中数据")
	}
	ret := dao.db.Table(dao.tableName).Table(dao.tableName).Where("route_id in (?)", ids).Delete(&craftRouteModels.SysCraftRouteRequest{})
	return ret.Error
}

func (dao *CraftRouteDaoImpl) SelectCraftRoute(c *gin.Context, req *craftRouteModels.SysCraftRouteListReq) (*craftRouteModels.SysCraftRouteListData, error) {
	db := dao.db.Table(dao.tableName).Table(dao.tableName)

	if req != nil && req.RouteCode != "" {
		db = db.Where("route_code = ?", req.RouteCode)
	}

	if req != nil && req.RouteName != "" {
		db = db.Where("route_name LIKE ?", "%"+req.RouteName+"%")
	}
	if req != nil && req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	size := 0
	if req == nil || req.Size <= 0 {
		size = 20
	} else {
		size = int(req.Size)
	}
	offset := 0
	if req == nil || req.Page <= 0 {
		req.Page = 1
	} else {
		offset = int((req.Page - 1) * req.Size)
	}
	db = db.Where("state", commonStatus.NORMAL)
	db = baizeContext.GetGormDataScope(c, db)
	var dto []*craftRouteModels.SysCraftRoute

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &craftRouteModels.SysCraftRouteListData{
			Rows:  make([]*craftRouteModels.SysCraftRoute, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &craftRouteModels.SysCraftRouteListData{
			Rows:  make([]*craftRouteModels.SysCraftRoute, 0),
			Total: 0,
		}, ret.Error
	}
	return &craftRouteModels.SysCraftRouteListData{
		Rows:  dto,
		Total: total,
	}, nil
}

func (dao *CraftRouteDaoImpl) GetById(c *gin.Context, id int64) (*craftRouteModels.SysCraftRoute, error) {
	var data *craftRouteModels.SysCraftRoute
	ret := dao.db.Table(dao.tableName).Where("route_id = ?", id).Where("state = ?", commonStatus.NORMAL).First(&data)
	return data, ret.Error
}

func (dao *CraftRouteDaoImpl) GetByIds(c *gin.Context, ids []int64) ([]*craftRouteModels.SysCraftRoute, error) {
	var data []*craftRouteModels.SysCraftRoute
	ret := dao.db.Table(dao.tableName).Where("route_id in (?)", ids).Where("state = ?", commonStatus.NORMAL).First(&data)
	return data, ret.Error
}
