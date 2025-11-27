package aiDataSetDaoImpl

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/ai/aiDataSetDao"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gateway/v1/config/app/intercept/logalert"
	"nova-factory-server/app/utils/snowflake"
)

type IAiPredictionListDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIAiPredictionListDaoImpl(db *gorm.DB) aiDataSetDao.IAiPredictionListDao {
	return &IAiPredictionListDaoImpl{
		db:    db,
		table: "sys_ai_prediction_list",
	}
}

func (a *IAiPredictionListDaoImpl) Set(c *gin.Context, data *aiDataSetModels.SetSysAiPrediction) (*aiDataSetModels.SysAiPrediction, error) {
	value := aiDataSetModels.ToSysAiPredictionList(data)
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

func (a *IAiPredictionListDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := a.db.Table(a.table).Where("id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (a *IAiPredictionListDaoImpl) List(c *gin.Context, req *aiDataSetModels.SysAiPredictionListReq) (*aiDataSetModels.SysAiPredictionList, error) {
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
	var dto []*aiDataSetModels.SysAiPrediction

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &aiDataSetModels.SysAiPredictionList{
			Rows:  make([]*aiDataSetModels.SysAiPrediction, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &aiDataSetModels.SysAiPredictionList{
			Rows:  make([]*aiDataSetModels.SysAiPrediction, 0),
			Total: 0,
		}, ret.Error
	}

	for k, v := range dto {
		var advanced logalert.Advanced
		err := json.Unmarshal([]byte(v.Advanced), &advanced)
		if err != nil {
			zap.L().Error("json Unmarshal fail", zap.Error(err))
			continue
		}
		dto[k].AdvancedData = &advanced

		var devList []string
		err = json.Unmarshal([]byte(v.Dev), &devList)
		if err != nil {
			zap.L().Error("json Unmarshal fail", zap.Error(err))
			continue
		}
		dto[k].DevList = devList
	}

	return &aiDataSetModels.SysAiPredictionList{
		Rows:  dto,
		Total: uint64(total),
	}, nil
}

func (a *IAiPredictionListDaoImpl) All(c *gin.Context) ([]*aiDataSetModels.SysAiPrediction, error) {
	var list []*aiDataSetModels.SysAiPrediction
	ret := a.db.Table(a.table).Where("state = ?", commonStatus.NORMAL).Find(&list)
	if ret.Error != nil {
		return list, ret.Error
	}
	return list, nil
}
