package aiDataSetDaoImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/ai/aiDataSetDao"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type IAiPredictionControlDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIAiPredictionControlDaoImpl(db *gorm.DB) aiDataSetDao.IAiPredictionControlDao {
	return &IAiPredictionControlDaoImpl{
		db:    db,
		table: "sys_ai_prediction_control",
	}
}

func (a *IAiPredictionControlDaoImpl) Set(c *gin.Context, data *aiDataSetModels.SetSysAiPredictionControl) (*aiDataSetModels.SysAiPredictionControl, error) {
	value := aiDataSetModels.ToSysAiPredictionControl(data)
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

func (a *IAiPredictionControlDaoImpl) List(c *gin.Context, req *aiDataSetModels.SysAiPredictionControlListReq) (*aiDataSetModels.SysAiPredictionControlList, error) {
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
	var dto []*aiDataSetModels.SysAiPredictionControl

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &aiDataSetModels.SysAiPredictionControlList{
			Rows:  make([]*aiDataSetModels.SysAiPredictionControl, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &aiDataSetModels.SysAiPredictionControlList{
			Rows:  make([]*aiDataSetModels.SysAiPredictionControl, 0),
			Total: 0,
		}, ret.Error
	}

	return &aiDataSetModels.SysAiPredictionControlList{
		Rows:  dto,
		Total: uint64(total),
	}, nil
}
