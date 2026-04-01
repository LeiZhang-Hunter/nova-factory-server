package alertserviceimpl

import (
	"nova-factory-server/app/business/iot/alert/alertdao"
	"nova-factory-server/app/business/iot/alert/alertmodels"
	"nova-factory-server/app/business/iot/alert/alertservice"

	"github.com/gin-gonic/gin"
)

type AlertAiReasonServiceImpl struct {
	aiDao alertdao.AlertAiReasonDao
}

func NewAlertAiReasonServiceImpl(aiDao alertdao.AlertAiReasonDao) alertservice.AlertAiReasonService {
	return &AlertAiReasonServiceImpl{
		aiDao: aiDao,
	}
}

func (a *AlertAiReasonServiceImpl) Set(c *gin.Context, data *alertmodels.SetAlertAiReason) (*alertmodels.SysAlertAiReason, error) {
	return a.aiDao.Set(c, data)
}
func (a *AlertAiReasonServiceImpl) Remove(c *gin.Context, ids []string) error {
	return a.aiDao.Remove(c, ids)
}
func (a *AlertAiReasonServiceImpl) List(c *gin.Context, req *alertmodels.SysAlertAiReasonReq) (*alertmodels.SysAlertReasonList, error) {
	return a.aiDao.List(c, req)
}
