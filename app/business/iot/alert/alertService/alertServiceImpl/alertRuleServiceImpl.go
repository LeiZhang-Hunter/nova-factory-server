package alertServiceImpl

import (
	"errors"
	alertDao2 "nova-factory-server/app/business/iot/alert/alertDao"
	alertModels2 "nova-factory-server/app/business/iot/alert/alertModels"
	"nova-factory-server/app/business/iot/alert/alertService"
	"nova-factory-server/app/business/iot/daemonize/daemonizeDao"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AlertRuleServiceImpl struct {
	rule        alertDao2.AlertRuleDao
	agent       daemonizeDao.IotAgentDao
	templateDao alertDao2.AlertSinkTemplateDao
	actionDao   alertDao2.AlertActionDao
	reasonDao   alertDao2.AlertAiReasonDao
}

func NewAlertRuleServiceImpl(rule alertDao2.AlertRuleDao, agent daemonizeDao.IotAgentDao,
	templateDao alertDao2.AlertSinkTemplateDao, actionDao alertDao2.AlertActionDao,
	reasonDao alertDao2.AlertAiReasonDao) alertService.AlertRuleService {
	return &AlertRuleServiceImpl{
		rule:        rule,
		agent:       agent,
		templateDao: templateDao,
		actionDao:   actionDao,
		reasonDao:   reasonDao,
	}
}

func (a *AlertRuleServiceImpl) Create(c *gin.Context, data *alertModels2.SetSysAlert) (*alertModels2.SysAlert, error) {
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

	// 检查处理id 是否存在
	actionInfo, err := a.actionDao.GetById(c, data.ActionId)
	if err != nil {
		zap.L().Error("get action info", zap.Error(err))
		return nil, errors.New("读取处理通知策略失败")
	}

	if actionInfo == nil {
		return nil, errors.New("处理通知策略不存在")
	}

	// 检查推理id 是否存在
	reasonInfo, err := a.reasonDao.GetById(c, data.ReasonId)
	if err != nil {
		zap.L().Error("get reason info", zap.Error(err))
		return nil, errors.New("读取ai 推理策略失败")
	}

	if reasonInfo == nil {
		return nil, errors.New("ai 推理策略不存在")
	}
	return a.rule.Create(c, data)
}
func (a *AlertRuleServiceImpl) Update(c *gin.Context, data *alertModels2.SetSysAlert) (*alertModels2.SysAlert, error) {
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
func (a *AlertRuleServiceImpl) List(c *gin.Context, req *alertModels2.SysAlertListReq) (*alertModels2.SysAlertList, error) {
	return a.rule.List(c, req)
}
func (a *AlertRuleServiceImpl) Remove(c *gin.Context, ids []string) error {
	return a.rule.Remove(c, ids)
}

func (a *AlertRuleServiceImpl) Change(c *gin.Context, data *alertModels2.ChangeSysAlert) error {
	return a.rule.Change(c, data)
}

func (a *AlertRuleServiceImpl) FindOpen(c *gin.Context, gatewayId int64) (*alertModels2.SysAlert, error) {
	return a.rule.FindOpen(c, gatewayId)
}

// GetReasonByGatewayId 通过网关id读取推理策略
func (a *AlertRuleServiceImpl) GetReasonByGatewayId(c *gin.Context, gatewayId int64) (*alertModels2.SysAlertAiReason, error) {
	info, err := a.rule.FindOpen(c, gatewayId)
	if err != nil {
		return nil, err
	}

	if info == nil {
		return nil, nil
	}

	reasonInfo, err := a.reasonDao.GetById(c, info.ReasonId)
	if err != nil {
		return nil, err
	}

	if reasonInfo == nil {
		return nil, nil
	}
	return reasonInfo, nil
}
