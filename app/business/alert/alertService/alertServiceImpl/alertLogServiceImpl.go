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

	alertLogList, _ := alertModels.FromDataToNovaAlertLog(&data)
	err = log.dao.Export(c, alertLogList)
	if err != nil {
		return err
	}

	//for _, info := range infos {
	//	alertController.GetAlertRunner().Push(info)
	//}
	return nil
}

func (log *AlertLogServiceImpl) List(c *gin.Context, req *alertModels.SysAlertLogListReq) (*alertModels.NovaAlertLogList, error) {
	list, err := log.dao.List(c, req)
	if err != nil {
		return &alertModels.NovaAlertLogList{
			Rows: make([]*alertModels.NovaAlertLog, 0),
		}, err
	}

	if list == nil {
		return &alertModels.NovaAlertLogList{
			Rows: make([]*alertModels.NovaAlertLog, 0),
		}, nil
	}

	var ruleIdMap map[uint64]bool = make(map[uint64]bool)
	for _, v := range list.Rows {
		ruleIdMap[v.AlertID] = true
	}

	var ruleIds []uint64 = make([]uint64, 0)
	for k, _ := range ruleIdMap {
		ruleIds = append(ruleIds, k)
	}

	rules, err := log.ruleDao.GetByIds(c, ruleIds)
	if err != nil {
		return &alertModels.NovaAlertLogList{
			Rows: make([]*alertModels.NovaAlertLog, 0),
		}, err
	}

	var rulesMap map[uint64]*alertModels.SysAlert
	for _, v := range rules {
		rulesMap[uint64(v.ID)] = v
	}

	for k, v := range list.Rows {
		rule, ok := rulesMap[v.AlertID]
		if !ok {
			list.Rows[k].RuleName = "系统告警"
			continue
		}
		list.Rows[k].RuleName = rule.Name
	}
	return list, nil
}

// Count 统计clickhouse数据库
func (log *AlertLogServiceImpl) Count(c *gin.Context) (int64, error) {
	return log.dao.Count(c)
}
