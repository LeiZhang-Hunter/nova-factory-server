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

type IScheduleDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIScheduleDaoImpl(db *gorm.DB) craftRouteDao.IScheduleDao {
	return &IScheduleDaoImpl{
		db:    db,
		table: "sys_product_schedule",
	}
}

func (i *IScheduleDaoImpl) GetDailySchedule(c *gin.Context) ([]*craftRouteModels.SysProductSchedule, error) {
	var list []*craftRouteModels.SysProductSchedule
	db := i.db.Table(i.table).Where("schedule_type = ?", craftRouteModels.DAILY).
		Where("state = ?", commonStatus.NORMAL)
	db = baizeContext.GetGormDataScope(c, db)
	ret := db.Find(&list)
	if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return list, ret.Error
}

func (i *IScheduleDaoImpl) List(c *gin.Context, req *craftRouteModels.SysProductScheduleListReq) (*craftRouteModels.SysProductScheduleListData, error) {
	db := i.db.Table(i.table)

	if req != nil && req.Name != "" {
		db = db.Where("schedule_name = ?", "%"+req.Name+"%")
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
	var dto []*craftRouteModels.SysProductSchedule

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &craftRouteModels.SysProductScheduleListData{
			Rows:  make([]*craftRouteModels.SysProductSchedule, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Debug().Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &craftRouteModels.SysProductScheduleListData{
			Rows:  make([]*craftRouteModels.SysProductSchedule, 0),
			Total: 0,
		}, ret.Error
	}
	return &craftRouteModels.SysProductScheduleListData{
		Rows:  dto,
		Total: total,
	}, nil
}

func (i *IScheduleDaoImpl) Set(c *gin.Context, tx *gorm.DB, data *craftRouteModels.SetSysProductSchedule) (*craftRouteModels.SysProductSchedule, error) {
	value := craftRouteModels.ToSysProductSchedule(data)
	if data.Id == 0 {
		value.ID = snowflake.GenID()
		value.DeptID = baizeContext.GetDeptId(c)
		value.SetCreateBy(baizeContext.GetUserId(c))
		ret := tx.Table(i.table).Create(value)
		return value, ret.Error
	}
	value.SetUpdateBy(baizeContext.GetUserId(c))
	ret := tx.Table(i.table).Updates(value)
	return value, ret.Error
}

func (i *IScheduleDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := i.db.Table(i.table).Where("id in (?)", ids).Delete(&craftRouteModels.SysProductSchedule{})
	return ret.Error
}

func (i *IScheduleDaoImpl) GetById(c *gin.Context, id int64) (*craftRouteModels.SysProductSchedule, error) {
	var data *craftRouteModels.SysProductSchedule
	ret := i.db.Table(i.table).Where("id = ?", id).First(&data)
	return data, ret.Error
}
