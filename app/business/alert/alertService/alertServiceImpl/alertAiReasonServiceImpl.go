package alertServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/alert/alertDao"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/alert/alertService"
)

type AlertAiReasonServiceImpl struct {
	aiDao alertDao.AlertAiReasonDao
}

func NewAlertAiReasonServiceImpl(aiDao alertDao.AlertAiReasonDao) alertService.AlertAiReasonService {
	return &AlertAiReasonServiceImpl{
		aiDao: aiDao,
	}
}

func (a *AlertAiReasonServiceImpl) Set(c *gin.Context, data *alertModels.SetAlertAiReason) (*alertModels.SysAlertAiReason, error) {
	return a.aiDao.Set(c, data)
}
func (a *AlertAiReasonServiceImpl) Remove(c *gin.Context, ids []string) error {
	return a.aiDao.Remove(c, ids)
}
func (a *AlertAiReasonServiceImpl) List(c *gin.Context, req *alertModels.SysAlertAiReasonReq) (*alertModels.SysAlertReasonList, error) {
	return a.aiDao.List(c, req)
}
