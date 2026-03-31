package alertserviceimpl

import (
	"errors"
	"nova-factory-server/app/business/iot/alert/alertdao"
	"nova-factory-server/app/business/iot/alert/alertmodels"
	"nova-factory-server/app/business/iot/alert/alertservice"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AlertTemplateServiceImpl struct {
	dao alertdao.AlertSinkTemplateDao
}

func NewAlertTemplateServiceImpl(dao alertdao.AlertSinkTemplateDao) alertservice.AlertTemplateService {
	return &AlertTemplateServiceImpl{
		dao: dao,
	}
}

func (a *AlertTemplateServiceImpl) Create(c *gin.Context, data *alertmodels.SetSysAlertSinkTemplate) (*alertmodels.SysAlertSinkTemplate, error) {
	info, err := a.dao.Create(c, data)
	if err != nil {
		zap.L().Error("create sink error", zap.Error(err))
		return nil, err
	}
	return info, nil
}
func (a *AlertTemplateServiceImpl) Update(c *gin.Context, data *alertmodels.SetSysAlertSinkTemplate) (*alertmodels.SysAlertSinkTemplate, error) {
	info, err := a.dao.GetById(c, uint64(data.ID))
	if err != nil {
		zap.L().Error("get gateway info error", zap.Error(err))
		return nil, err
	}
	if info == nil {
		return nil, errors.New("数据不存在")
	}
	value, err := a.dao.Update(c, data)
	if err != nil {
		zap.L().Error("update sink error", zap.Error(err))
		return nil, err
	}
	return value, nil
}
func (a *AlertTemplateServiceImpl) List(c *gin.Context, req *alertmodels.SysAlertSinkTemplateReq) (*alertmodels.SysAlertSinkTemplateListData, error) {
	return a.dao.List(c, req)
}
func (a *AlertTemplateServiceImpl) Remove(c *gin.Context, ids []string) error {
	return a.dao.Remove(c, ids)
}
