package craftRouteServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/business/craft/craftRouteService"
)

type ICraftProcessServiceImpl struct {
	dao craftRouteDao.IProcessDao
}

func NewICraftProcessServiceImpl(dao craftRouteDao.IProcessDao) craftRouteService.ICraftProcessService {
	return &ICraftProcessServiceImpl{
		dao: dao,
	}
}

func (i *ICraftProcessServiceImpl) Add(c *gin.Context, process *craftRouteModels.SysProProcess) (*craftRouteModels.SysProProcess, error) {
	process, err := i.dao.Add(c, process)
	if err != nil {
		return nil, err
	}
	return process, nil
}

func (i *ICraftProcessServiceImpl) Update(c *gin.Context, process *craftRouteModels.SysProProcess) (*craftRouteModels.SysProProcess, error) {
	process, err := i.dao.Update(c, process)
	if err != nil {
		return nil, err
	}
	return process, nil
}
func (i *ICraftProcessServiceImpl) Remove(c *gin.Context, processIds []int64) error {
	return i.dao.Remove(c, processIds)
}
