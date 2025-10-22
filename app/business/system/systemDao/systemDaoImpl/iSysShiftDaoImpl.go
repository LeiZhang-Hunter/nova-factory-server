package systemDaoImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/business/system/systemModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

// ISysShiftDaoImpl 班次设置
type ISysShiftDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewISysShiftDaoImpl(db *gorm.DB) systemDao.ISysShiftDao {
	return &ISysShiftDaoImpl{
		db:    db,
		table: "sys_work_shift_setting",
	}
}

func (i *ISysShiftDaoImpl) Set(c *gin.Context, valueVO *systemModels.SysWorkShiftSettingVO) (*systemModels.SysWorkShiftSetting, error) {
	value, err := systemModels.ToSysWorkShiftSetting(valueVO)
	if err != nil {
		return nil, err
	}
	if value.ID == 0 {
		value.SetCreateBy(baizeContext.GetUserId(c))
		value.ID = snowflake.GenID()
		value.DeptID = baizeContext.GetDeptId(c)
		ret := i.db.Table(i.table).Create(&value)
		return value, ret.Error
	} else {
		value.SetUpdateBy(baizeContext.GetUserId(c))
		ret := i.db.Table(i.table).Debug().Where("id = ?", value.ID).Updates(&value)
		return value, ret.Error
	}
}

func (i *ISysShiftDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := i.db.Table(i.table).Where("id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (i *ISysShiftDaoImpl) List(c *gin.Context, req *systemModels.SysWorkShiftSettingReq) (*systemModels.SysWorkShiftSettingList, error) {
	db := i.db.Table(i.table)
	if req == nil {
		req = &systemModels.SysWorkShiftSettingReq{}
	}
	if req.Name != "" {
		db = db.Where("name like ?", "%"+req.Name+"%")
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

	db = baizeContext.GetGormDataScope(c, db)
	db = db.Where("state = ?", commonStatus.NORMAL)

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &systemModels.SysWorkShiftSettingList{
			Rows:  make([]*systemModels.SysWorkShiftSetting, 0),
			Total: 0,
		}, ret.Error
	}
	var dto []*systemModels.SysWorkShiftSetting
	ret = db.Offset(offset).Limit(size).Order("create_time desc").Find(&dto)
	if ret.Error != nil {
		return &systemModels.SysWorkShiftSettingList{
			Rows:  make([]*systemModels.SysWorkShiftSetting, 0),
			Total: 0,
		}, ret.Error
	}

	return &systemModels.SysWorkShiftSettingList{
		Rows:  dto,
		Total: total,
	}, nil
}

// Check 查询时间是否冲突
func (i *ISysShiftDaoImpl) Check(c *gin.Context, id int64, startTime int32, endTime int32) *systemModels.SysWorkShiftSetting {
	var info *systemModels.SysWorkShiftSetting
	_ = i.db.Table(i.table).Where("id != ?", id).Where("begin_time < ?", endTime).Where("end_time  > ?", endTime).Where("status = ?", true).Where("state = ?", commonStatus.NORMAL).First(&info)
	if info != nil && info.ID != 0 {
		return info
	}
	_ = i.db.Table(i.table).Where("id != ?", id).Where("begin_time < ?", startTime).Where("end_time  > ?", startTime).Where("status = ?", true).Where("state = ?", commonStatus.NORMAL).First(&info)
	if info != nil && info.ID != 0 {
		return info
	}
	_ = i.db.Table(i.table).Where("id != ?", id).Where("begin_time < ?", startTime).Where("end_time  > ?", endTime).Where("status = ?", true).Where("state = ?", commonStatus.NORMAL).First(&info)
	if info != nil && info.ID != 0 {
		return info
	}
	_ = i.db.Table(i.table).Where("id != ?", id).Where("begin_time > ?", startTime).Where("end_time  < ?", endTime).Where("status = ?", true).Where("state = ?", commonStatus.NORMAL).First(&info)
	if info != nil && info.ID != 0 {
		return info
	}
	if info.ID == 0 {
		return nil
	}
	return nil
}

func (i *ISysShiftDaoImpl) GetEnableShift(c *gin.Context) ([]*systemModels.SysWorkShiftSetting, error) {
	var dto []*systemModels.SysWorkShiftSetting
	ret := i.db.Table(i.table).Where("status = ?", 1).Where("state = ?", commonStatus.NORMAL).Find(&dto)
	if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return dto, ret.Error
}
