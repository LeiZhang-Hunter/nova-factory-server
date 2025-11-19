package alertServiceImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/alert/alertController"
	"nova-factory-server/app/business/alert/alertDao"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/alert/alertService"
	"nova-factory-server/app/business/daemonize/daemonizeDao"
)

type AlertLogServiceImpl struct {
	dao      alertDao.AlertLogClickhouseDao
	ruleDao  alertDao.AlertRuleDao
	agentDao daemonizeDao.IotAgentDao
}

func NewAlertLogServiceImpl(dao alertDao.AlertLogClickhouseDao, ruleDao alertDao.AlertRuleDao, agentDao daemonizeDao.IotAgentDao) alertService.AlertLogService {
	return &AlertLogServiceImpl{
		dao:      dao,
		ruleDao:  ruleDao,
		agentDao: agentDao,
	}
}

// Export 导出告警数据
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

	alertLogList, infos := alertModels.FromDataToNovaAlertLog(&data)
	err = log.dao.Export(c, alertLogList)
	if err != nil {
		return err
	}

	for _, info := range infos {
		alertController.GetAlertRunner().Push(info)
	}
	return nil
}

func (log *AlertLogServiceImpl) List(c *gin.Context, req *alertModels.SysAlertLogListReq) (*alertModels.NovaAlertLogList, error) {
	return log.dao.List(c, req)
}
