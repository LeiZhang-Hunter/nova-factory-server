package aiDataSetDaoImpl

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type IAiPredictionExceptionDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIAiPredictionExceptionDaoImpl(db *gorm.DB) aidatasetdao.IAiPredictionExceptionDao {
	return &IAiPredictionExceptionDaoImpl{
		db:    db,
		table: "sys_ai_prediction_exception",
	}
}

func (a *IAiPredictionExceptionDaoImpl) Set(c *gin.Context, data *aidatasetmodels.SetSysAiPredictionException) (*aidatasetmodels.SysAiPredictionException, error) {
	value := aidatasetmodels.ToSysAiPredictionException(data)
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

func (a *IAiPredictionExceptionDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := a.db.Table(a.table).Where("id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (a *IAiPredictionExceptionDaoImpl) List(c *gin.Context, req *aidatasetmodels.SysAiPredictionExceptionListReq) (*aidatasetmodels.SysAiPredictionExceptionList, error) {
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
	var dto []*aidatasetmodels.SysAiPredictionException

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &aidatasetmodels.SysAiPredictionExceptionList{
			Rows:  make([]*aidatasetmodels.SysAiPredictionException, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &aidatasetmodels.SysAiPredictionExceptionList{
			Rows:  make([]*aidatasetmodels.SysAiPredictionException, 0),
			Total: 0,
		}, ret.Error
	}

	for k, v := range dto {
		var devList []string
		err := json.Unmarshal([]byte(v.Dev), &devList)
		if err != nil {
			zap.L().Error("json Unmarshal fail", zap.Error(err))
			continue
		}
		dto[k].DevList = devList
	}

	return &aidatasetmodels.SysAiPredictionExceptionList{
		Rows:  dto,
		Total: uint64(total),
	}, nil
}

func (a *IAiPredictionExceptionDaoImpl) All(c *gin.Context) ([]*aidatasetmodels.SysAiPredictionException, error) {
	var list []*aidatasetmodels.SysAiPredictionException
	ret := a.db.Table(a.table).Where("state = ?", commonStatus.NORMAL).Find(&list)
	if ret.Error != nil {
		return list, ret.Error
	}
	return list, nil
}
