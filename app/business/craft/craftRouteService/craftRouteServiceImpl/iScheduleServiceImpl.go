package craftRouteServiceImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/utils/time"
	systime "time"
)

type IScheduleServiceImpl struct {
	scheduleDao    craftRouteDao.IScheduleDao
	scheduleMapDao craftRouteDao.IScheduleMapDao
}

func NewIScheduleServiceImpl(scheduleDao craftRouteDao.IScheduleDao, scheduleMapDao craftRouteDao.IScheduleMapDao) craftRouteService.IScheduleService {
	return &IScheduleServiceImpl{
		scheduleDao:    scheduleDao,
		scheduleMapDao: scheduleMapDao,
	}
}

func (i *IScheduleServiceImpl) getTimeList(monthBegin systime.Time) ([]int64, error) {
	// 搜索特殊日期
	week := time.GetWeek(monthBegin)
	// 计算出本月开始时间是周几
	if week == 7 {
		week = 0
	}

	//计算出格子的第一个日期
	var beginDataTime int64
	beginDataTime = monthBegin.Unix() - int64(week*24*3600)
	if beginDataTime < 0 {
		return nil, errors.New("日期格式错误")
	}
	ret := make([]int64, 42)
	for count := 0; count < 42; count++ {
		ret[count] = beginDataTime + int64(24*count*3600)
	}
	return ret, nil
}

// formatSpecialRet 规整特殊日程,用来判断日期方格里有没有特殊日程
func (i *IScheduleServiceImpl) formatSpecialRet(list []*craftRouteModels.SysProductScheduleMap) map[int64]*craftRouteModels.SysProductScheduleMap {
	var data map[int64]*craftRouteModels.SysProductScheduleMap = make(map[int64]*craftRouteModels.SysProductScheduleMap)
	for _, v := range list {
		recordTime := systime.Unix(int64(v.BeginTime), 0)
		// 提取年、月、日
		year, month, day := recordTime.Date()

		// 构建当天零点的时间
		midnight := systime.Date(year, month, day, 0, 0, 0, 0, recordTime.Location())
		data[midnight.Unix()] = v
	}
	return data
}

func (i *IScheduleServiceImpl) GetMonthSchedule(c *gin.Context, req *craftRouteModels.SysProductScheduleReq) ([]*craftRouteModels.ScheduleStatusData, error) {
	// 搜索普通日程
	scheduleList, err := i.scheduleDao.GetDailySchedule(c)
	if err != nil {
		zap.L().Error("get daily schedule", zap.Error(err))
		return make([]*craftRouteModels.ScheduleStatusData, 0), err
	}

	var scheduleIds []int64
	for _, schedule := range scheduleList {
		scheduleIds = append(scheduleIds, schedule.ID)
	}

	// 读取循环信息列表
	schedules, err := i.scheduleMapDao.GetByScheduleIds(c, scheduleIds)
	if err != nil {
		return make([]*craftRouteModels.ScheduleStatusData, 0), err
	}

	// 循环日期信息
	var dailyIds map[int]int
	for _, schedule := range schedules {
		dailyIds[schedule.Date] = schedule.Date
	}

	monthBegin := time.GetMonthStart(req.Year, req.Month)
	timeList, err := i.getTimeList(monthBegin)
	if err != nil {
		zap.L().Error("get month time list", zap.Error(err))
		return nil, err
	}

	// 读取特殊日程
	specialList, err := i.scheduleMapDao.GetSpecialSchedule(c, monthBegin.Unix())
	if err != nil {
		return make([]*craftRouteModels.ScheduleStatusData, 0), err
	}
	specialMap := i.formatSpecialRet(specialList)

	data := make([]*craftRouteModels.ScheduleStatusData, 42)
	for k, v := range timeList {
		data[k] = &craftRouteModels.ScheduleStatusData{}
		data[k].Time = systime.Unix(v, 0)

		// 检查是否有特殊日程
		_, ok := specialMap[v]
		if ok {
			data[k] = &craftRouteModels.ScheduleStatusData{
				Type: craftRouteModels.SPECIAL,
			}
			continue
		}

		// 检查是否有循环日程
		scheduleWeek := int(data[k].Time.Weekday())
		_, ok = dailyIds[scheduleWeek]
		if !ok {
			continue
		}
		data[k] = &craftRouteModels.ScheduleStatusData{
			Type: craftRouteModels.DAILY,
		}
	}
	return data, nil
}

func (i *IScheduleServiceImpl) List(c *gin.Context, req *craftRouteModels.SysProductScheduleListReq) (*craftRouteModels.SysProductScheduleListData, error) {
	return i.scheduleDao.List(c, req)
}

func (i *IScheduleServiceImpl) Set(c *gin.Context, data *craftRouteModels.SetSysProductSchedule) {
	// 保存调度主表
	value, err := i.scheduleDao.Set(c, data)
	if err != nil {
		return
	}

	// 保存调度子表
	i.scheduleMapDao.Set(c, value)
	return
}
