package craftRouteServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
)

type ISysProTaskServiceImpl struct {
	dao     craftRouteDao.ISysProWorkorderDao
	taskDao craftRouteDao.ISysProTaskDao
}

func NewISysProTaskServiceImpl(dao craftRouteDao.ISysProWorkorderDao, taskDao craftRouteDao.ISysProTaskDao) craftRouteService.ISysProTaskService {
	return &ISysProTaskServiceImpl{
		dao:     dao,
		taskDao: taskDao,
	}
}

func (i *ISysProTaskServiceImpl) Add(ctx *gin.Context, req *craftRouteModels.SysSetProTask) (*craftRouteModels.SysProTask, error) {
	return i.taskDao.Add(ctx, req)
}

func (i *ISysProTaskServiceImpl) Update(ctx *gin.Context, req *craftRouteModels.SysSetProTask) (*craftRouteModels.SysProTask, error) {
	return i.taskDao.Update(ctx, req)
}

func (i *ISysProTaskServiceImpl) Remove(ctx *gin.Context, ids []string) error {
	return i.taskDao.Remove(ctx, ids)
}

func (i *ISysProTaskServiceImpl) List(ctx *gin.Context, req *craftRouteModels.SysProTaskReq) (*craftRouteModels.SysProTaskList, error) {
	return i.taskDao.List(ctx, req)
}

func (i *ISysProTaskServiceImpl) Schedule(ctx *gin.Context, req *craftRouteModels.ScheduleReq) ([]*craftRouteModels.SysProTask, error) {
	// 读取最新的任务,获取7天之内的50个任务
	return i.taskDao.Schedule(ctx, req)
}
