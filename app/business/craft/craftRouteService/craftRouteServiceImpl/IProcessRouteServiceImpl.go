package craftRouteServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type IProcessRouteServiceImpl struct {
	dao craftRouteDao.IRouteProcessDao
}

func NewIProcessRouteServiceImpl(dao craftRouteDao.IRouteProcessDao) craftRouteService.IProcessRouteService {
	return &IProcessRouteServiceImpl{
		dao: dao,
	}
}

func (i *IProcessRouteServiceImpl) Add(c *gin.Context, request *craftRouteModels.SysProRouteProcessSetRequest) (*craftRouteModels.SysProRouteProcess, error) {
	var data craftRouteModels.SysProRouteProcess
	data.Of(request)
	data.RecordID = snowflake.GenID()
	data.DeptID = baizeContext.GetDeptId(c)
	data.SetCreateBy(baizeContext.GetUserId(c))
	return i.dao.Add(c, &data)
}

func (i *IProcessRouteServiceImpl) Update(c *gin.Context, request *craftRouteModels.SysProRouteProcessSetRequest) (*craftRouteModels.SysProRouteProcess, error) {
	var data craftRouteModels.SysProRouteProcess
	data.Of(request)
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
