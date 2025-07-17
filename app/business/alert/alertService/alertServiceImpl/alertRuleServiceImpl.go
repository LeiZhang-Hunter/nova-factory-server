package alertServiceImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/alert/alertDao"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/alert/alertService"
	"nova-factory-server/app/business/daemonize/daemonizeDao"
)

type AlertRuleServiceImpl struct {
	rule        alertDao.AlertRuleDao
	agent       daemonizeDao.IotAgentDao
	templateDao alertDao.AlertSinkTemplateDao
}

func NewAlertRuleServiceImpl(rule alertDao.AlertRuleDao, agent daemonizeDao.IotAgentDao, templateDao alertDao.AlertSinkTemplateDao) alertService.AlertRuleService {
	return &AlertRuleServiceImpl{
		rule:        rule,
		agent:       agent,
		templateDao: templateDao,
	}
}

func (a *AlertRuleServiceImpl) Create(c *gin.Context, data *alertModels.SetSysAlert) (*alertModels.SysAlert, error) {
	gatewayInfo, err := a.agent.GetByObjectId(c, uint64(data.GatewayID))
	if err != nil {
		return nil, err
	}
	if gatewayInfo == nil {
		return nil, errors.New("网关不存在")
	}

	templateInfo, err := a.templateDao.GetById(c, uint64(data.TemplateID))
	if err != nil {
		return nil, err
	}
	if templateInfo == nil {
		return nil, errors.New("告警模板不存在")
	}
	return a.rule.Create(c, data)
}
func (a *AlertRuleServiceImpl) Update(c *gin.Context, data *alertModels.SetSysAlert) (*alertModels.SysAlert, error) {
	gatewayInfo, err := a.agent.GetByObjectId(c, uint64(data.GatewayID))
	if err != nil {
		return nil, err
	}
	if gatewayInfo == nil {
		return nil, errors.New("网关不存在")
	}

	templateInfo, err := a.templateDao.GetById(c, uint64(data.TemplateID))
	if err != nil {
		return nil, err
	}
	if templateInfo == nil {
		return nil, errors.New("告警模板不存在")
	}
	return a.rule.Update(c, data)
}
func (a *AlertRuleServiceImpl) List(c *gin.Context, req *alertModels.SysAlertListReq) (*alertModels.SysAlertList, error) {
	return a.rule.List(c, req)
}
func (a *AlertRuleServiceImpl) Remove(c *gin.Context, ids []string) error {
	return a.rule.Remove(c, ids)
}

func (a *AlertRuleServiceImpl) Change(c *gin.Context, data *alertModels.ChangeSysAlert) error {
	return a.rule.Change(c, data)
}
