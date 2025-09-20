package aiDataSetDaoImpl

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/ai/aiDataSetDao"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type IAiPredictionExceptionDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIAiPredictionExceptionDaoImpl(db *gorm.DB) aiDataSetDao.IAiPredictionExceptionDao {
	return &IAiPredictionExceptionDaoImpl{
		db:    db,
		table: "sys_ai_prediction_exception",
	}
}

func (a *IAiPredictionExceptionDaoImpl) Set(c *gin.Context, data *aiDataSetModels.SetSysAiPredictionException) (*aiDataSetModels.SysAiPredictionException, error) {
	value := aiDataSetModels.ToSysAiPredictionException(data)
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

func (a *IAiPredictionExceptionDaoImpl) List(c *gin.Context, req *aiDataSetModels.SysAiPredictionExceptionListReq) (*aiDataSetModels.SysAiPredictionExceptionList, error) {
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
	var dto []*aiDataSetModels.SysAiPredictionException

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &aiDataSetModels.SysAiPredictionExceptionList{
			Rows:  make([]*aiDataSetModels.SysAiPredictionException, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &aiDataSetModels.SysAiPredictionExceptionList{
			Rows:  make([]*aiDataSetModels.SysAiPredictionException, 0),
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

	return &aiDataSetModels.SysAiPredictionExceptionList{
		Rows:  dto,
		Total: uint64(total),
	}, nil
}
