package craftRouteDaoImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/constant/task"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
	"time"
)

type ISysProTaskDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewISysProTaskDaoImpl(db *gorm.DB) craftRouteDao.ISysProTaskDao {
	return &ISysProTaskDaoImpl{
		db:    db,
		table: "sys_pro_task",
	}
}

func (i *ISysProTaskDaoImpl) Add(c *gin.Context, info *craftRouteModels.SysSetProTask) (*craftRouteModels.SysProTask, error) {
	data := craftRouteModels.OfSysProTask(info)
	data.WorkorderID = snowflake.GenID()
	data.DeptID = baizeContext.GetDeptId(c)
	data.SetCreateBy(baizeContext.GetUserId(c))
	ret := i.db.Table(i.table).Create(data)
	return data, ret.Error
}

func (i *ISysProTaskDaoImpl) Update(c *gin.Context, info *craftRouteModels.SysSetProTask) (*craftRouteModels.SysProTask, error) {
	data := craftRouteModels.OfSysProTask(info)
	data.SetUpdateBy(baizeContext.GetUserId(c))
	ret := i.db.Table(i.table).Where("task_id = ?", data.WorkorderID).Updates(data)
	return data, ret.Error
}

func (i *ISysProTaskDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := i.db.Table(i.table).Where("task_id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (i *ISysProTaskDaoImpl) List(ctx *gin.Context, req *craftRouteModels.SysProTaskReq) (*craftRouteModels.SysProTaskList, error) {
	db := i.db.Table(i.table)

	size := 0

	if req.WorkorderID != 0 {
		db = db.Where("workorder_id = ?", req.WorkorderID)
	} else {
		return &craftRouteModels.SysProTaskList{
			Rows:  make([]*craftRouteModels.SysProTask, 0),
			Total: 0,
		}, nil
	}

	offset := 0
	if req == nil || req.Page <= 0 {
		req.Page = 1
	} else {
		offset = int((req.Page - 1) * req.Size)
	}

	db = db.Where("state", commonStatus.NORMAL)
	db = baizeContext.GetGormDataScope(ctx, db)
	var dto []*craftRouteModels.SysProTask

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &craftRouteModels.SysProTaskList{
			Rows:  make([]*craftRouteModels.SysProTask, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &craftRouteModels.SysProTaskList{
			Rows:  make([]*craftRouteModels.SysProTask, 0),
			Total: 0,
		}, ret.Error
	}
	return &craftRouteModels.SysProTaskList{
		Rows:  dto,
		Total: total,
	}, nil
}

// Schedule 读取最新的任务,获取7天之内的50个任务
func (i *ISysProTaskDaoImpl) Schedule(ctx *gin.Context, req *craftRouteModels.ScheduleReq) ([]*craftRouteModels.SysProTask, error) {
	var dto []*craftRouteModels.SysProTask
	// 计算七天前的时间
	now := time.Now()
	sevenDaysAgo := now.AddDate(0, 0, -7)
	ret := i.db.Table(i.table).Where("gateway_id = ?", req.GatewayId).Where("start_time < ?", now).Where("start_time >= ?", sevenDaysAgo).
		Where("status = ?", task.TASK_STATUS_NORMAL).Limit(int(req.Size)).Where("state = ?", commonStatus.NORMAL).
		Find(&dto)
	if ret.Error != nil {
		return nil, ret.Error
	}

	return dto, nil
}
