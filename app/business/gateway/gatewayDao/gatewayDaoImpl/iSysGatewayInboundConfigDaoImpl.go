package gatewayDaoImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/gateway/gatewayDao"
	"nova-factory-server/app/business/gateway/gatewayModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type ISysGatewayInboundConfigDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewISysGatewayInboundConfigDaoImpl(db *gorm.DB) gatewayDao.ISysGatewayInboundConfigDao {
	return &ISysGatewayInboundConfigDaoImpl{
		db:        db,
		tableName: "sys_gateway_inbound_config",
	}
}

func (i *ISysGatewayInboundConfigDaoImpl) Add(c *gin.Context, config *gatewayModels.SysSetGatewayInboundConfig) (*gatewayModels.SysGatewayInboundConfig, error) {
	data := gatewayModels.NewSysGatewayInboundConfig(config)
	data.GatewayConfigID = snowflake.GenID()
	data.SetCreateBy(baizeContext.GetUserId(c))
	ret := i.db.Create(data)
	return data, ret.Error
}

func (i *ISysGatewayInboundConfigDaoImpl) Update(c *gin.Context, config *gatewayModels.SysSetGatewayInboundConfig) (*gatewayModels.SysGatewayInboundConfig, error) {
	data := gatewayModels.NewSysGatewayInboundConfig(config)
	data.SetUpdateBy(baizeContext.GetUserId(c))
	ret := i.db.Where("gateway_config_id = ?", config.GatewayConfigID).Updates(data)
	return data, ret.Error
}

func (i *ISysGatewayInboundConfigDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := i.db.Table(i.tableName).Where("gateway_config_id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (i *ISysGatewayInboundConfigDaoImpl) List(c *gin.Context, req *gatewayModels.SysSetGatewayInboundConfigReq) (*gatewayModels.SysSetGatewayInboundConfigList, error) {
	db := i.db.Table(i.tableName).Table(i.tableName)

	if req != nil && req.Name != "" {
		db = db.Where("name like ?", "%"+req.Name+"%")
	}

	if req != nil && req.ProtocolType != "" {
		db = db.Where("protocol_type = ?", req.ProtocolType)
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
	var dto []*gatewayModels.SysGatewayInboundConfig

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &gatewayModels.SysSetGatewayInboundConfigList{
			Rows:  make([]*gatewayModels.SysGatewayInboundConfig, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &gatewayModels.SysSetGatewayInboundConfigList{
			Rows:  make([]*gatewayModels.SysGatewayInboundConfig, 0),
			Total: 0,
		}, ret.Error
	}
	return &gatewayModels.SysSetGatewayInboundConfigList{
		Rows:  dto,
		Total: total,
	}, nil
}
