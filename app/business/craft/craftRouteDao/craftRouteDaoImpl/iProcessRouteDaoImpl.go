package craftRouteDaoImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
)

type IProcessRouteDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewIProcessRouteDaoImpl(db *gorm.DB) craftRouteDao.IRouteProcessDao {
	return &IProcessRouteDaoImpl{
		db:        db,
		tableName: "sys_pro_route_process",
	}
}

func (i *IProcessRouteDaoImpl) Add(c *gin.Context, data *craftRouteModels.SysProRouteProcess) (*craftRouteModels.SysProRouteProcess, error) {
	ret := i.db.Table(i.tableName).Create(data)
	return data, ret.Error
}

func (i *IProcessRouteDaoImpl) Update(c *gin.Context, data *craftRouteModels.SysProRouteProcess) (*craftRouteModels.SysProRouteProcess, error) {
	ret := i.db.Table(i.tableName).Where("record_id = ?", data.RecordID).Updates(data)
	return data, ret.Error

}

func (i *IProcessRouteDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := i.db.Table(i.tableName).Where("record_id in (?)", ids).Update("state", -1)
	return ret.Error
}

func (i *IProcessRouteDaoImpl) List(c *gin.Context, req *craftRouteModels.SysProRouteProcessListReq) (*craftRouteModels.SysProRouteProcessList, error) {
	db := i.db.Table(i.tableName).Table(i.tableName)

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
	var dto []*craftRouteModels.SysProRouteProcess

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &craftRouteModels.SysProRouteProcessList{
			Rows:  make([]*craftRouteModels.SysProRouteProcess, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &craftRouteModels.SysProRouteProcessList{
			Rows:  make([]*craftRouteModels.SysProRouteProcess, 0),
			Total: 0,
		}, ret.Error
	}
	return &craftRouteModels.SysProRouteProcessList{
		Rows:  dto,
		Total: total,
	}, nil
}
