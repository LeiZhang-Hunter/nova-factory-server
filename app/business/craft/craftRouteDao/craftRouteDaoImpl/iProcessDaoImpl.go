package craftRouteDaoImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type IProcessDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewIProcessDaoImpl(db *gorm.DB) craftRouteDao.IProcessDao {
	return &IProcessDaoImpl{
		db:        db,
		tableName: "sys_pro_process",
	}
}

func (i *IProcessDaoImpl) Add(c *gin.Context, req *craftRouteModels.SysProSetProcessReq) (*craftRouteModels.SysProProcess, error) {
	data := craftRouteModels.NewSysProProcess(req)
	data.ProcessID = snowflake.GenID()
	data.SetCreateBy(baizeContext.GetUserId(c))
	var info *craftRouteModels.SysProProcess
	ret := i.db.Table(i.tableName).Where("process_code = ?", req.ProcessCode).Where("state = ?", commonStatus.NORMAL).First(&info)
	if ret.Error != nil {
		if !errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return nil, ret.Error
		} else {
			ret = i.db.Table(i.tableName).Create(data)
			return data, ret.Error
		}
	}
	if info == nil {
		ret = i.db.Table(i.tableName).Create(data)
		return data, ret.Error
	}
	data.ProcessID = info.ProcessID
	data.CreateTime = info.CreateTime
	ret = i.db.Table(i.tableName).Where("process_id = ?", info.ProcessID).Updates(data)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return data, ret.Error
}

func (i *IProcessDaoImpl) Update(c *gin.Context, req *craftRouteModels.SysProSetProcessReq) (*craftRouteModels.SysProProcess, error) {
	data := craftRouteModels.NewSysProProcess(req)
	data.SetUpdateBy(baizeContext.GetUserId(c))
	ret := i.db.Table(i.tableName).Where("process_id = ?", data.ProcessID).Updates(data)
	return data, ret.Error
}

func (i *IProcessDaoImpl) Remove(c *gin.Context, processIds []int64) error {
	ret := i.db.Table(i.tableName).Where("process_id in (?)", processIds).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (i *IProcessDaoImpl) List(c *gin.Context, req *craftRouteModels.SysProProcessListReq) (*craftRouteModels.SysProProcessListData, error) {
	db := i.db.Table(i.tableName).Table(i.tableName)

	if req != nil && req.ProcessCode != "" {
		db = db.Where("process_code = ?", req.ProcessCode)
	}

	if req != nil && req.ProcessName != "" {
		db = db.Where("process_name = ?", req.ProcessName)
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
	var dto []*craftRouteModels.SysProProcess

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &craftRouteModels.SysProProcessListData{
			Rows:  make([]*craftRouteModels.SysProProcess, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &craftRouteModels.SysProProcessListData{
			Rows:  make([]*craftRouteModels.SysProProcess, 0),
			Total: 0,
		}, ret.Error
	}
	return &craftRouteModels.SysProProcessListData{
		Rows:  dto,
		Total: total,
	}, nil
}

func (i *IProcessDaoImpl) GetById(c *gin.Context, id int64) (*craftRouteModels.SysProProcess, error) {
	var data *craftRouteModels.SysProProcess
	ret := i.db.Table(i.tableName).Where("process_id = ?", id).Where("state = ?", commonStatus.NORMAL).First(&data)
	return data, ret.Error
}

func (i *IProcessDaoImpl) GetByIds(c *gin.Context, ids []int64) ([]*craftRouteModels.SysProProcess, error) {
	var data []*craftRouteModels.SysProProcess
	ret := i.db.Table(i.tableName).Where("process_id in (?)", ids).Where("state = ?", commonStatus.NORMAL).Find(&data)
	return data, ret.Error
}
