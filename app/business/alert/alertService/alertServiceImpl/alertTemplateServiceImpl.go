package alertServiceImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/alert/alertDao"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/alert/alertService"
)

type AlertTemplateServiceImpl struct {
	dao alertDao.AlertSinkTemplateDao
}

func NewAlertTemplateServiceImpl(dao alertDao.AlertSinkTemplateDao) alertService.AlertTemplateService {
	return &AlertTemplateServiceImpl{
		dao: dao,
	}
}

func (a *AlertTemplateServiceImpl) Create(c *gin.Context, data *alertModels.SetSysAlertSinkTemplate) (*alertModels.SysAlertSinkTemplate, error) {
	gatewayInfo, err := a.dao.GetByGatewayId(c, uint64(data.GatewayID))
	if err != nil {
		zap.L().Error("get gateway info error", zap.Error(err))
		return nil, err
	}
	if gatewayInfo != nil {
		return nil, errors.New("网关已经存在不能重复设置")
	}
	info, err := a.dao.Create(c, data)
	if err != nil {
		zap.L().Error("create sink error", zap.Error(err))
		return nil, err
	}
	return info, nil
}
func (a *AlertTemplateServiceImpl) Update(c *gin.Context, data *alertModels.SetSysAlertSinkTemplate) (*alertModels.SysAlertSinkTemplate, error) {
	info, err := a.dao.GetById(c, uint64(data.GatewayID))
	if err != nil {
		zap.L().Error("get gateway info error", zap.Error(err))
		return nil, err
	}
	if info == nil {
		return nil, errors.New("数据不存在")
	}
	data.GatewayID = info.GatewayID
	value, err := a.dao.Update(c, data)
	if err != nil {
		zap.L().Error("update sink error", zap.Error(err))
		return nil, err
	}
	return value, nil
}
func (a *AlertTemplateServiceImpl) List(c *gin.Context, req *alertModels.SysAlertSinkTemplateReq) (*alertModels.SysAlertSinkTemplateListData, error) {
	return a.dao.List(c, req)
}
func (a *AlertTemplateServiceImpl) Remove(c *gin.Context, ids []string) error {
	return a.dao.Remove(c, ids)
}
