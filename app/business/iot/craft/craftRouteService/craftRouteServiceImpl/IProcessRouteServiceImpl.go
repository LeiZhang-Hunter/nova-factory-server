package craftRouteServiceImpl

import (
	"errors"
	craftRouteDao2 "nova-factory-server/app/business/iot/craft/craftRouteDao"
	"nova-factory-server/app/business/iot/craft/craftRouteModels"
	"nova-factory-server/app/business/iot/craft/craftRouteService"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IProcessRouteServiceImpl struct {
	dao           craftRouteDao2.IRouteProcessDao
	processDao    craftRouteDao2.IProcessDao
	craftRouteDao craftRouteDao2.ICraftRouteDao
}

func NewIProcessRouteServiceImpl(dao craftRouteDao2.IRouteProcessDao, processDao craftRouteDao2.IProcessDao, craftRouteDao craftRouteDao2.ICraftRouteDao) craftRouteService.IProcessRouteService {
	return &IProcessRouteServiceImpl{
		dao:           dao,
		processDao:    processDao,
		craftRouteDao: craftRouteDao,
	}
}

func (i *IProcessRouteServiceImpl) Add(c *gin.Context, request *craftRouteModels.SysProRouteProcessSetRequest) (*craftRouteModels.SysProRouteProcess, error) {
	var data craftRouteModels.SysProRouteProcess
	data.Of(request)
	info, err := i.processDao.GetById(c, data.ProcessID)
	if err != nil {
		zap.L().Error("工序读取错误", zap.Error(err))
		return nil, errors.New("工序不存在")
	}
	if info == nil {
		return nil, errors.New("工序不存在")
	}
	routeInfo, err := i.craftRouteDao.GetById(c, data.RouteID)
	if err != nil {
		zap.L().Error("工艺路线读取错误", zap.Error(err))
		return nil, errors.New("工艺路线不存在")
	}
	if routeInfo == nil {
		return nil, errors.New("工艺路线不存在")
	}
	data.RecordID = snowflake.GenID()
	data.DeptID = baizeContext.GetDeptId(c)
	data.SetCreateBy(baizeContext.GetUserId(c))
	return i.dao.Add(c, &data)
}

func (i *IProcessRouteServiceImpl) Update(c *gin.Context, request *craftRouteModels.SysProRouteProcessSetRequest) (*craftRouteModels.SysProRouteProcess, error) {
	var data craftRouteModels.SysProRouteProcess
	data.Of(request)
	info, err := i.processDao.GetById(c, data.ProcessID)
	if err != nil {
		zap.L().Error("工序读取错误", zap.Error(err))
		return nil, errors.New("工序不存在")
	}
	if info == nil {
		return nil, errors.New("工序不存在")
	}
	routeInfo, err := i.craftRouteDao.GetById(c, data.RouteID)
	if err != nil {
		zap.L().Error("工艺路线读取错误", zap.Error(err))
		return nil, errors.New("工艺路线不存在")
	}
	if routeInfo == nil {
		return nil, errors.New("工艺路线不存在")
	}
	data.DeptID = baizeContext.GetDeptId(c)
	data.SetUpdateBy(baizeContext.GetUserId(c))
	return i.dao.Update(c, &data)
}

func (i *IProcessRouteServiceImpl) List(c *gin.Context, req *craftRouteModels.SysProRouteProcessListReq) (*craftRouteModels.SysProRouteProcessList, error) {
	return i.dao.List(c, req)
}

func (i *IProcessRouteServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
