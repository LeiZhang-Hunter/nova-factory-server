package alertserviceimpl

import (
	"nova-factory-server/app/business/iot/alert/alertdao"
	"nova-factory-server/app/business/iot/alert/alertmodels"
	"nova-factory-server/app/business/iot/alert/alertservice"

	"github.com/gin-gonic/gin"
)

type AlertActionServiceImpl struct {
	dao alertdao.AlertActionDao
}

func NewAlertActionServiceImpl(dao alertdao.AlertActionDao) alertservice.AlertActionService {
	return &AlertActionServiceImpl{
		dao: dao,
	}
}

func (a *AlertActionServiceImpl) Set(c *gin.Context, data *alertmodels.SetAlertAction) (*alertmodels.AlertAction, error) {
	return a.dao.Set(c, data)
}
func (a *AlertActionServiceImpl) Remove(c *gin.Context, ids []string) error {
	return a.dao.Remove(c, ids)
}
func (a *AlertActionServiceImpl) List(c *gin.Context, req *alertmodels.SysAlertActionListReq) (*alertmodels.SysAlertActionList, error) {
	return a.dao.List(c, req)
}
