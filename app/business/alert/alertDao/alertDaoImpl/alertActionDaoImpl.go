package alertDaoImpl

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/alert/alertDao"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
	"nova-factory-server/app/utils/time"
)

type AlertActionDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewAlertActionDaoImpl(db *gorm.DB) alertDao.AlertActionDao {
	return &AlertActionDaoImpl{
		db:    db,
		table: "sys_alert_action",
	}
}

func (a *AlertActionDaoImpl) Set(c *gin.Context, data *alertModels.SetAlertAction) (*alertModels.AlertAction, error) {
	value := alertModels.FromSetAlertActionToData(data)
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

func (a *AlertActionDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := a.db.Table(a.table).Where("id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (a *AlertActionDaoImpl) List(c *gin.Context, req *alertModels.SysAlertActionListReq) (*alertModels.SysAlertActionList, error) {
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
	var dto []*alertModels.AlertAction

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &alertModels.SysAlertActionList{
			Rows:  make([]*alertModels.AlertAction, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &alertModels.SysAlertActionList{
			Rows:  make([]*alertModels.AlertAction, 0),
			Total: 0,
		}, ret.Error
	}

	for k, v := range dto {
		dto[k].UserNotifyList = make([]alertModels.UserNotify, 0)
		err := json.Unmarshal([]byte(v.UserNotify), &dto[k].UserNotifyList)
		if err != nil {
			zap.L().Error("json unmarshal error", zap.Error(err))
		}
		var userCountMap map[int64]uint32 = make(map[int64]uint32)

		dto[k].ApiNotifyList = make([]alertModels.ApiNotify, 0)
		err = json.Unmarshal([]byte(v.ApiNotify), &dto[k].ApiNotifyList)
		if err != nil {
			zap.L().Error("json unmarshal error", zap.Error(err))
		}

		for dataK, data := range dto[k].UserNotifyList {
			start := time.SecondsToTime(data.TimeStart)
			end := time.SecondsToTime(data.TimeEnd)
			dto[k].UserNotifyList[dataK].TimeRange = make([]string, 2)
			dto[k].UserNotifyList[dataK].TimeRange[0] = start
			dto[k].UserNotifyList[dataK].TimeRange[1] = end
			for _, dtoKData := range data.Receiver {
				userCountMap[dtoKData.UserId] = 0
			}
		}

		for dataK, data := range dto[k].ApiNotifyList {
			start := time.SecondsToTime(data.TimeStart)
			end := time.SecondsToTime(data.TimeEnd)
			dto[k].ApiNotifyList[dataK].TimeRange = make([]string, 2)
			dto[k].ApiNotifyList[dataK].TimeRange[0] = start
			dto[k].ApiNotifyList[dataK].TimeRange[1] = end
		}

		dto[k].UserCount = uint32(len(userCountMap))
	}

	return &alertModels.SysAlertActionList{
		Rows:  dto,
		Total: uint64(total),
	}, nil
}
