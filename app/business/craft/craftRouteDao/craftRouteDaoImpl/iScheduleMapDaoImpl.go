package craftRouteDaoImpl

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
	"strconv"
	"strings"
	"time"
)

type IScheduleMapDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIScheduleMapDaoImpl(db *gorm.DB) craftRouteDao.IScheduleMapDao {
	return &IScheduleMapDaoImpl{
		db:    db,
		table: "sys_product_schedule_map",
	}
}

func (i *IScheduleMapDaoImpl) GetByScheduleIds(c *gin.Context, ids []int64) ([]*craftRouteModels.SysProductScheduleMap, error) {
	var list []*craftRouteModels.SysProductScheduleMap
	ret := i.db.Table(i.table).Where("schedule_id in (?)", ids).Where("state = ?", commonStatus.NORMAL).Find(&list)
	if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
		return list, nil
	}
	return list, ret.Error
}

func (i *IScheduleMapDaoImpl) GetSpecialSchedule(c *gin.Context, beginTime int64) ([]*craftRouteModels.SysProductScheduleMap, error) {
	var list []*craftRouteModels.SysProductScheduleMap
	db := i.db.Table(i.table).Where("begin_time >= ?", beginTime).Where("schedule_type = ?", craftRouteModels.SPECIAL).
		Where("state = ?", commonStatus.NORMAL)
	db = baizeContext.GetGormDataScope(c, db)
	ret := db.Find(&list)
	if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return list, ret.Error
}

// dealDaily 处理循环日程
func (i *IScheduleMapDaoImpl) dealDaily(c *gin.Context, data *craftRouteModels.SetSysProductSchedule) error {
	weeks := strings.Split(data.Time, ",")
	if len(weeks) == 0 {
		return errors.New("循环日期格式错误")
	}

	for _, v := range data.TimeManager {
		beginArr := strings.Split(":", v.BeginTime)
		if len(beginArr) != 2 {
			return errors.New(fmt.Sprintf("日程安排错误:%s", v.BeginTime))
		}
		endArr := strings.Split(":", v.EndTime)
		if len(endArr) != 2 {
			return errors.New(fmt.Sprintf("日程安排错误:%s", v.EndTime))
		}

		beginHour, err := strconv.ParseInt(beginArr[0], 10, 64)
		if err != nil {
			zap.L().Error("dealDaily error", zap.Error(err))
			return err
		}

		beginMinute, err := strconv.ParseInt(beginArr[0], 10, 64)
		if err != nil {
			zap.L().Error("dealDaily error", zap.Error(err))
			return err
		}

		beginUnix := beginHour*3600 + beginMinute*60

		endHour, err := strconv.ParseInt(beginArr[0], 10, 64)
		if err != nil {
			zap.L().Error("dealDaily error", zap.Error(err))
			return err
		}

		endMinute, err := strconv.ParseInt(beginArr[0], 10, 64)
		if err != nil {
			zap.L().Error("dealDaily error", zap.Error(err))
			return err
		}

		endUnix := endHour*3600 + endMinute*60
		var mapList []*craftRouteModels.SysProductScheduleMap = make([]*craftRouteModels.SysProductScheduleMap, len(weeks))
		for dayKey, dayValue := range weeks {
			date, err := strconv.ParseInt(dayValue, 10, 10)
			if err != nil {
				zap.L().Error("dealDaily error", zap.Error(err))
				return err
			}
			mapList[dayKey] = &craftRouteModels.SysProductScheduleMap{
				ID:           snowflake.GenID(),
				ScheduleID:   data.Id,
				BeginTime:    beginUnix,
				EndTime:      endUnix,
				CraftRouteID: v.RoueId,
				ScheduleType: craftRouteModels.DAILY,
				Date:         int(date),
			}
		}
	}

	return nil
}

// dealSpecial 处理特殊日程
func (i *IScheduleMapDaoImpl) dealSpecial(c *gin.Context, data *craftRouteModels.SetSysProductSchedule) error {
	if data.Time == "" {
		return errors.New("执行日期格式错误")
	}

	if len(data.TimeManager) == 0 {
		return errors.New("日程安排格式错误")
	}

	times := strings.Split(data.Time, "~")
	if len(times) != 2 {
		return errors.New("执行日期格式错误")
	}

	beginTIme := times[0]
	endTIme := times[1]
	beginTImeValue, err := time.Parse("2025-07-02", beginTIme)
	if err != nil {
		return err
	}
	endTImeValue, err := time.Parse("2025-07-02", endTIme)
	if err != nil {
		return err
	}
	beginTimeUnix := beginTImeValue.Unix()
	endTimeUnix := endTImeValue.Unix()
	if beginTimeUnix == endTimeUnix {
		return errors.New("执行日期开始和结束不能相同")
	}
	if beginTimeUnix+24*3600 > endTimeUnix {
		return errors.New("日期错误")
	}

	dayCount := 0
	day := (endTimeUnix - beginTimeUnix) / 86400
	dayList := make([]int64, 0)
	for ; dayCount < int(day); dayCount++ {
		var v int64 = beginTimeUnix + int64(dayCount)*86400
		dayList = append(dayList, v)
	}

	for _, v := range data.TimeManager {
		beginArr := strings.Split(":", v.BeginTime)
		if len(beginArr) != 2 {
			return errors.New(fmt.Sprintf("日程安排错误:%s", v.BeginTime))
		}
		endArr := strings.Split(":", v.EndTime)
		if len(endArr) != 2 {
			return errors.New(fmt.Sprintf("日程安排错误:%s", v.EndTime))
		}

		beginHour, err := strconv.ParseInt(beginArr[0], 10, 64)
		if err != nil {
			zap.L().Error("dealDaily error", zap.Error(err))
			return err
		}

		beginMinute, err := strconv.ParseInt(beginArr[0], 10, 64)
		if err != nil {
			zap.L().Error("dealDaily error", zap.Error(err))
			return err
		}

		beginUnix := beginHour*3600 + beginMinute*60

		endHour, err := strconv.ParseInt(beginArr[0], 10, 64)
		if err != nil {
			zap.L().Error("dealDaily error", zap.Error(err))
			return err
		}

		endMinute, err := strconv.ParseInt(beginArr[0], 10, 64)
		if err != nil {
			zap.L().Error("dealDaily error", zap.Error(err))
			return err
		}

		endUnix := endHour*3600 + endMinute*60
		var mapList []*craftRouteModels.SysProductScheduleMap = make([]*craftRouteModels.SysProductScheduleMap, len(dayList))
		for dayKey, dayValue := range dayList {
			mapList[dayKey] = &craftRouteModels.SysProductScheduleMap{
				ID:           snowflake.GenID(),
				ScheduleID:   data.Id,
				BeginTime:    dayValue + beginUnix,
				EndTime:      dayValue + endUnix,
				CraftRouteID: v.RoueId,
				ScheduleType: craftRouteModels.SPECIAL,
			}
		}
		ret := i.db.Table(i.table).Create(mapList)
		if ret.Error != nil {
			zap.L().Error("create error", zap.Error(ret.Error))
			return ret.Error
		}
	}
	return nil
}

func (i *IScheduleMapDaoImpl) Set(c *gin.Context, data *craftRouteModels.SetSysProductSchedule) {
	i.db.Table(i.table).Where("schedule_id = ?", data.Id).Delete(&craftRouteModels.SetSysProductSchedule{})
	if data.Type == craftRouteModels.SPECIAL {
		err := i.dealSpecial(c, data)
		if err != nil {
			return
		}
	} else {
		err := i.dealDaily(c, data)
		if err != nil {
			return
		}
	}
}

func (i *IScheduleMapDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := i.db.Table(i.table).Where("schedule_id in (?)", ids).Delete(&craftRouteModels.SysProductScheduleMap{})
	return ret.Error
}

func (i *IScheduleMapDaoImpl) GetByScheduleId(c *gin.Context, id int64) ([]*craftRouteModels.SysProductScheduleMap, error) {
	var dto []*craftRouteModels.SysProductScheduleMap
	ret := i.db.Table(i.table).Where("schedule_id = ?", id).Find(&dto)
	return dto, ret.Error
}
