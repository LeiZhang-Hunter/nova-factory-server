package alertdaoimpl

import (
	"encoding/json"
	"errors"
	"nova-factory-server/app/business/iot/alert/alertdao"
	"nova-factory-server/app/business/iot/alert/alertmodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
	"nova-factory-server/app/utils/time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AlertActionDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewAlertActionDaoImpl(db *gorm.DB) alertdao.AlertActionDao {
	return &AlertActionDaoImpl{
		db:    db,
		table: "sys_alert_action",
	}
}

func (a *AlertActionDaoImpl) Set(c *gin.Context, data *alertmodels.SetAlertAction) (*alertmodels.AlertAction, error) {
	value := alertmodels.FromSetAlertActionToData(data)
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

func (a *AlertActionDaoImpl) List(c *gin.Context, req *alertmodels.SysAlertActionListReq) (*alertmodels.SysAlertActionList, error) {
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
	var dto []*alertmodels.AlertAction

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &alertmodels.SysAlertActionList{
			Rows:  make([]*alertmodels.AlertAction, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &alertmodels.SysAlertActionList{
			Rows:  make([]*alertmodels.AlertAction, 0),
			Total: 0,
		}, ret.Error
	}

	for k, v := range dto {
		dto[k].UserNotifyList = make([]alertmodels.UserNotify, 0)
		err := json.Unmarshal([]byte(v.UserNotify), &dto[k].UserNotifyList)
		if err != nil {
			zap.L().Error("json unmarshal error", zap.Error(err))
		}

		var userCountMap map[int64]uint32 = make(map[int64]uint32)
		dto[k].ApiNotifyList = make([]alertmodels.ApiNotify, 0)
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

	return &alertmodels.SysAlertActionList{
		Rows:  dto,
		Total: uint64(total),
	}, nil
}

func (a *AlertActionDaoImpl) GetById(c *gin.Context, id int64) (*alertmodels.SetAlertAction, error) {
	var dto *alertmodels.SetAlertAction
	ret := a.db.Table(a.table).Where("id = ?", id).Where("state = ?", commonStatus.NORMAL).First(&dto)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return dto, ret.Error
}
