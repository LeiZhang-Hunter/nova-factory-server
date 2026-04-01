package craftroutedaoimpl

import (
	"encoding/json"
	"nova-factory-server/app/business/iot/craft/craftroutedao"
	"nova-factory-server/app/business/iot/craft/craftroutemodels"
	"nova-factory-server/app/business/iot/craft/craftroutemodels/api/v1"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ISysCraftRouteConfigDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewISysCraftRouteConfigDaoImpl(db *gorm.DB) craftroutedao.ISysCraftRouteConfigDao {
	return &ISysCraftRouteConfigDaoImpl{
		db:        db,
		tableName: "sys_craft_route_config",
	}
}

func (i *ISysCraftRouteConfigDaoImpl) Save(c *gin.Context, routeId uint64, topo *craftroutemodels.ProcessTopo, configContent []byte) (*craftroutemodels.SysCraftRouteConfig, error) {
	var info *craftroutemodels.SysCraftRouteConfig
	content, err := json.Marshal(topo)
	if err != nil {
		return nil, err
	}
	ret := i.db.Table(i.tableName).Where("route_id = ?", routeId).Where("state = ?", commonStatus.NORMAL).First(&info)
	if ret.Error != nil {
		zap.L().Error("get info error", zap.Error(ret.Error))
		info = nil
	}

	// 读取工序列表，工序内容

	if info == nil {
		var config craftroutemodels.SysCraftRouteConfig
		config.RouteConfigID = snowflake.GenID()
		config.RouteID = topo.Route.RouteID
		config.Context = string(content)
		config.Config = string(configContent)
		config.SetCreateBy(baizeContext.GetUserId(c))
		ret = i.db.Table(i.tableName).Create(&config)
		if ret.Error != nil {
			zap.L().Error("save info error", zap.Error(ret.Error))
		}

		config.Topo = topo
		return &config, ret.Error
	} else {
		info.Context = string(content)
		info.Config = string(configContent)
		info.SetUpdateBy(baizeContext.GetUserId(c))
		ret = i.db.Table(i.tableName).Where("route_id = ?", routeId).Updates(&info)
		if ret.Error != nil {
			zap.L().Error("save info error", zap.Error(ret.Error))
		}
	}
	info.Topo = topo
	return info, nil
}

func (i *ISysCraftRouteConfigDaoImpl) GetById(routeId uint64) (*craftroutemodels.SysCraftRouteConfig, error) {
	var info *craftroutemodels.SysCraftRouteConfig
	ret := i.db.Table(i.tableName).Where("route_id = ?", routeId).Where("state = ?", commonStatus.NORMAL).First(&info)
	if ret.Error != nil {
		zap.L().Error("get info error", zap.Error(ret.Error))
		return nil, ret.Error
	}
	var topo craftroutemodels.ProcessTopo
	err := json.Unmarshal([]byte(info.Context), &topo)
	if err != nil {
		zap.L().Error("json unmarshal error", zap.Error(err))
		return nil, err
	}
	info.Topo = &topo
	return info, nil
}

func (i *ISysCraftRouteConfigDaoImpl) GetConfigByIds(routeIds []int64) ([]*v1.Router, error) {
	var list []*craftroutemodels.SysCraftRouteConfig
	ret := i.db.Table(i.tableName).Where("route_id in (?)", routeIds).Where("state = ?", commonStatus.NORMAL).Find(&list)
	if ret.Error != nil {
		zap.L().Error("get info error", zap.Error(ret.Error))
		return nil, ret.Error
	}

	routers := make([]*v1.Router, 0)
	for _, info := range list {
		var router v1.Router
		err := json.Unmarshal([]byte(info.Config), &router)
		if err != nil {
			zap.L().Error("json unmarshal error", zap.Error(err))
			return nil, err
		}
		routers = append(routers, &router)
	}

	return routers, nil
}
