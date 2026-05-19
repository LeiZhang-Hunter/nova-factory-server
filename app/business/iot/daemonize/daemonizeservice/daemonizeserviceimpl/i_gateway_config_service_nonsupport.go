//go:build !ai
// +build !ai

package daemonizeserviceimpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/admin/system/systemdao"
	alertDao2 "nova-factory-server/app/business/iot/alert/alertdao"
	"nova-factory-server/app/business/iot/asset/camera/cameradao"
	deviceDao2 "nova-factory-server/app/business/iot/asset/device/devicedao"
	"nova-factory-server/app/business/iot/daemonize/daemonizeservice"
	"nova-factory-server/app/business/iot/system/dao"
	"nova-factory-server/app/utils/gateway/v1/config/pipeline"
)

type iGatewayConfigServiceImpl struct {
	deviceDao         deviceDao2.IDeviceDao
	templateDao       deviceDao2.IDeviceTemplateDao
	templateDataDao   deviceDao2.ISysModbusDeviceConfigDataDao
	cameraDao         cameradao.ICameraDao
	ruleDao           alertDao2.AlertRuleDao
	sinkDao           alertDao2.AlertSinkTemplateDao
	electricConfigDao dao.IDeviceElectricDao
	dictDataDao       systemdao.IDictDataDao
}

func NewIGatewayConfigServiceImpl(
	deviceDao deviceDao2.IDeviceDao,
	templateDao deviceDao2.IDeviceTemplateDao,
	templateDataDao deviceDao2.ISysModbusDeviceConfigDataDao,
	cameraDao cameradao.ICameraDao,
	ruleDao alertDao2.AlertRuleDao,
	sinkDao alertDao2.AlertSinkTemplateDao, electricConfigDao dao.IDeviceElectricDao,
	dictDataDao systemdao.IDictDataDao) daemonizeservice.IGatewayConfigService {
	return &iGatewayConfigServiceImpl{
		deviceDao:         deviceDao,
		templateDao:       templateDao,
		templateDataDao:   templateDataDao,
		cameraDao:         cameraDao,
		ruleDao:           ruleDao,
		sinkDao:           sinkDao,
		electricConfigDao: electricConfigDao,
		dictDataDao:       dictDataDao,
	}
}

func (i *iGatewayConfigServiceImpl) loadPredictionConfig(c *gin.Context, piplines *pipeline.PipelineConfig) error {
	return nil
}
