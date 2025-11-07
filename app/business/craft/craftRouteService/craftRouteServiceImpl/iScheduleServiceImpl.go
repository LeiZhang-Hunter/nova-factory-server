package craftRouteServiceImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	v1 "nova-factory-server/app/business/craft/craftRouteModels/api/v1"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/utils/time"
	systime "time"
)

type IScheduleServiceImpl struct {
	scheduleDao    craftRouteDao.IScheduleDao
	scheduleMapDao craftRouteDao.IScheduleMapDao
	routeDao       craftRouteDao.ICraftRouteDao
	routeConfigDao craftRouteDao.ISysCraftRouteConfigDao
	db             *gorm.DB
}

func NewIScheduleServiceImpl(scheduleDao craftRouteDao.IScheduleDao, db *gorm.DB,
	scheduleMapDao craftRouteDao.IScheduleMapDao,
	routeDao craftRouteDao.ICraftRouteDao,
	routeConfigDao craftRouteDao.ISysCraftRouteConfigDao) craftRouteService.IScheduleService {
	return &IScheduleServiceImpl{
		scheduleDao:    scheduleDao,
		scheduleMapDao: scheduleMapDao,
		routeDao:       routeDao,
		db:             db,
		routeConfigDao: routeConfigDao,
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
	var dailyIds map[int]int = make(map[int]int)
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
				Time: systime.Unix(v, 0),
				Type: craftRouteModels.SPECIAL,
			}
			continue
		}

		// 检查是否有循环日程
		scheduleWeek := int(data[k].Time.Weekday())
		if scheduleWeek == 0 {
			scheduleWeek = 7
		}
		_, ok = dailyIds[scheduleWeek]
		if !ok {
			continue
		}
		data[k] = &craftRouteModels.ScheduleStatusData{
			Time: systime.Unix(v, 0),
			Type: craftRouteModels.DAILY,
		}
	}
	return data, nil
}

func (i *IScheduleServiceImpl) List(c *gin.Context, req *craftRouteModels.SysProductScheduleListReq) (*craftRouteModels.SysProductScheduleListData, error) {
	return i.scheduleDao.List(c, req)
}

func (i *IScheduleServiceImpl) Set(c *gin.Context, data *craftRouteModels.SetSysProductSchedule) error {
	// 校验控制id
	var routerIds []int64 = make([]int64, 0)
	for _, v := range data.TimeManager {
		if v.RoueId == 0 {
			return errors.New("控制流程id不存在")
		}
		routerIds = append(routerIds, v.RoueId)
	}
	var err error
	var list []*craftRouteModels.SysCraftRoute
	list, err = i.routeDao.GetByIds(c, routerIds)
	if err != nil {
		zap.L().Error("get by id list", zap.Error(err))
		return err
	}
	if len(list) != len(routerIds) {
		return errors.New("控制流程参数错误")
	}

	tx := i.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	// 保存调度主表
	var value *craftRouteModels.SysProductSchedule
	value, err = i.scheduleDao.Set(c, tx, data)
	if err != nil {
		return nil
	}

	data.Id = value.ID
	// 保存调度子表
	err = i.scheduleMapDao.Set(c, tx, data)
	return err
}

func (i *IScheduleServiceImpl) Remove(c *gin.Context, ids []string) error {
	// 保存调度主表
	err := i.scheduleDao.Remove(c, ids)
	if err != nil {
		zap.L().Error("remove schedule", zap.Error(err))
		return err
	}

	// 保存调度子表
	err = i.scheduleMapDao.Remove(c, ids)
	if err != nil {
		zap.L().Error("remove schedule", zap.Error(err))
		return err
	}
	return nil
}

func (i *IScheduleServiceImpl) Detail(c *gin.Context, id int64) (*craftRouteModels.DetailSysProductData, error) {
	scheduleInfo, err := i.scheduleDao.GetById(c, id)
	if err != nil {
		return &craftRouteModels.DetailSysProductData{}, err
	}
	scheduleMapList, err := i.scheduleMapDao.GetByScheduleId(c, id)
	if err != nil {
		return &craftRouteModels.DetailSysProductData{}, err
	}

	return &craftRouteModels.DetailSysProductData{
		Info: scheduleInfo,
		Data: scheduleMapList,
	}, nil
}

func (i *IScheduleServiceImpl) GetRouters(ctx *gin.Context, req *craftRouteModels.ScheduleReq) ([]int64, error) {
	ids := make([]int64, 0)
	//查找特殊日程
	specialList, err := i.scheduleMapDao.GetSpecialScheduleByNow(ctx, req.GatewayId)
	if err != nil {
		return ids, err
	}

	if len(specialList) != 0 {
		for _, v := range specialList {
			ids = append(ids, v.CraftRouteID)
		}
		return ids, nil
	}

	normalList, err := i.scheduleMapDao.GetNormalByTime(ctx, req.GatewayId)
	if err != nil {
		return ids, err
	}

	if normalList == nil {
		for _, v := range normalList {
			ids = append(ids, v.CraftRouteID)
		}
		return ids, nil
	}

	return ids, nil
}

// Schedule 任务调度
func (i *IScheduleServiceImpl) Schedule(ctx *gin.Context, req *craftRouteModels.ScheduleReq) ([]*v1.Router, error) {
	routerIds, err := i.GetRouters(ctx, req)
	if err != nil {
		zap.L().Error("Get routers error", zap.Error(err))
		return make([]*v1.Router, 0), err
	}
	if len(routerIds) == 0 {
		return make([]*v1.Router, 0), nil
	}

	routers, err := i.routeConfigDao.GetConfigByIds(routerIds)
	if err != nil {
		return make([]*v1.Router, 0), err
	}
	return routers, nil
}
