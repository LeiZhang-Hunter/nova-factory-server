package alertServiceImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/alert/alertDao"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/alert/alertService"
	"nova-factory-server/app/business/daemonize/daemonizeDao"
)

type AlertLogServiceImpl struct {
	dao      alertDao.AlertLogDao
	ruleDao  alertDao.AlertRuleDao
	agentDao daemonizeDao.IotAgentDao
}

func NewAlertLogServiceImpl(dao alertDao.AlertLogDao, ruleDao alertDao.AlertRuleDao, agentDao daemonizeDao.IotAgentDao) alertService.AlertLogService {
	return &AlertLogServiceImpl{
		dao:      dao,
		ruleDao:  ruleDao,
		agentDao: agentDao,
	}
}

func (log *AlertLogServiceImpl) Export(c *gin.Context, data alertModels.AlertLogData) error {
	// 校验网关帐号密码
	gatewayInfo, err := log.agentDao.GetByObjectId(c, uint64(data.GatewayId))
	if err != nil {
		zap.L().Error("export data error", zap.Error(err))
		return err
	}

	if gatewayInfo == nil {
		return errors.New("gateway not found")
	}

	if gatewayInfo.Username != data.Username || gatewayInfo.Password != data.Password {
		return errors.New("username or password error")
	}

	alertLogList := alertModels.FromDataToSysAlertLog(&data, uint64(gatewayInfo.DeptID), c)
	for _, alertLog := range alertLogList {
		alertLog.DeptID = gatewayInfo.DeptID
	}
	err = log.dao.Export(c, alertLogList)
	if err != nil {
		return err
	}
	return nil
}
