package aiDataSetDaoImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type IAiPredictionControlDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIAiPredictionControlDaoImpl(db *gorm.DB) aidatasetdao.IAiPredictionControlDao {
	return &IAiPredictionControlDaoImpl{
		db:    db,
		table: "sys_ai_prediction_control",
	}
}

func (a *IAiPredictionControlDaoImpl) Set(c *gin.Context, data *aidatasetmodels.SetSysAiPredictionControl) (*aidatasetmodels.SysAiPredictionControl, error) {
	value := aidatasetmodels.ToSysAiPredictionControl(data)
	if value.ID == 0 {
		value.SetCreateBy(baizeContext.GetUserId(c))
		value.ID = snowflake.GenID()
		value.DeptID = baizeContext.GetDeptId(c)
		ret := a.db.Table(a.table).Create(&value)
		return value, ret.Error
	} else {
		value.SetUpdateBy(baizeContext.GetUserId(c))
		ret := a.db.Table(a.table).Where("id = ?", value.ID).Updates(&value)
		return value, ret.Error
	}
}

func (a *IAiPredictionControlDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := a.db.Table(a.table).Where("id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (a *IAiPredictionControlDaoImpl) List(c *gin.Context, req *aidatasetmodels.SysAiPredictionControlListReq) (*aidatasetmodels.SysAiPredictionControlList, error) {
	db := a.db.Table(a.table)

	if req != nil && req.Name != "" {
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
	db = db.Where("state", commonStatus.NORMAL)
	db = baizeContext.GetGormDataScope(c, db)
	var dto []*aidatasetmodels.SysAiPredictionControl

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &aidatasetmodels.SysAiPredictionControlList{
			Rows:  make([]*aidatasetmodels.SysAiPredictionControl, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &aidatasetmodels.SysAiPredictionControlList{
			Rows:  make([]*aidatasetmodels.SysAiPredictionControl, 0),
			Total: 0,
		}, ret.Error
	}

	return &aidatasetmodels.SysAiPredictionControlList{
		Rows:  dto,
		Total: uint64(total),
	}, nil
}

func (a *IAiPredictionControlDaoImpl) Find(c *gin.Context) (*aidatasetmodels.SysAiPredictionControl, error) {
	var info *aidatasetmodels.SysAiPredictionControl
	ret := a.db.Table(a.table).Where("state = ?", commonStatus.NORMAL).First(&info)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, ret.Error
	}
	return info, nil
}
