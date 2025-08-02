package alertServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/alert/alertDao"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/alert/alertService"
)

type AlertActionServiceImpl struct {
	dao alertDao.AlertActionDao
}

func NewAlertActionServiceImpl(dao alertDao.AlertActionDao) alertService.AlertActionService {
	return &AlertActionServiceImpl{
		dao: dao,
	}
}

func (a *AlertActionServiceImpl) Set(c *gin.Context, data *alertModels.SetAlertAction) (*alertModels.AlertAction, error) {
	return a.dao.Set(c, data)
}
func (a *AlertActionServiceImpl) Remove(c *gin.Context, ids []string) error {
	return a.dao.Remove(c, ids)
}
func (a *AlertActionServiceImpl) List(c *gin.Context, req *alertModels.SysAlertActionListReq) (*alertModels.SysAlertActionList, error) {
	return a.dao.List(c, req)
}
