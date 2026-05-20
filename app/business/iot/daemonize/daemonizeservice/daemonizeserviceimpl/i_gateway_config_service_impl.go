//go:build ai
// +build ai

package daemonizeserviceimpl

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/admin/system/systemdao"
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	alertDao2 "nova-factory-server/app/business/iot/alert/alertdao"
	"nova-factory-server/app/business/iot/asset/camera/cameradao"
	deviceDao2 "nova-factory-server/app/business/iot/asset/device/devicedao"
	"nova-factory-server/app/business/iot/daemonize/daemonizeservice"
	"nova-factory-server/app/business/iot/system/dao"
	"nova-factory-server/app/utils/gateway/v1/config/app/source/prediction"
	"nova-factory-server/app/utils/gateway/v1/config/cfg"
	"nova-factory-server/app/utils/gateway/v1/config/pipeline"
	source2 "nova-factory-server/app/utils/gateway/v1/config/source"
	"time"
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
	predictionDao     aidatasetdao.IAiPredictionControlDao
}

func NewIGatewayConfigServiceImpl(
	deviceDao deviceDao2.IDeviceDao,
	templateDao deviceDao2.IDeviceTemplateDao,
	templateDataDao deviceDao2.ISysModbusDeviceConfigDataDao,
	cameraDao cameradao.ICameraDao,
	ruleDao alertDao2.AlertRuleDao,
	sinkDao alertDao2.AlertSinkTemplateDao, electricConfigDao dao.IDeviceElectricDao,
	dictDataDao systemdao.IDictDataDao, predictionDao aidatasetdao.IAiPredictionControlDao) daemonizeservice.IGatewayConfigService {
	return &iGatewayConfigServiceImpl{
		deviceDao:         deviceDao,
		templateDao:       templateDao,
		templateDataDao:   templateDataDao,
		cameraDao:         cameraDao,
		ruleDao:           ruleDao,
		sinkDao:           sinkDao,
		electricConfigDao: electricConfigDao,
		dictDataDao:       dictDataDao,
		predictionDao:     predictionDao,
	}
}

// loadPredictionConfig 加载预测配置
func (i *iGatewayConfigServiceImpl) loadPredictionConfig(c *gin.Context, piplines *pipeline.PipelineConfig) error {
	var scheduleEnabled bool = true
	predictionInfo, err := i.predictionDao.Find(c)
	if err != nil {
		zap.L().Error("i.predictionDao.Find() failed", zap.Error(err))
		return err
	}
	if predictionInfo != nil {
		//  查询预测配置
		predictionPipeline := pipeline.NewConfig()
		predictionPipeline.Name = "prediction"
		predictionPipelineSource := source2.Config{
			Enabled: &scheduleEnabled,
			Name:    "prediction",
			Type:    "prediction",
		}
		var predictionConfig prediction.Config
		predictionConfig.Parallelism = uint16(predictionInfo.Parallelism)
		predictionConfig.Model = predictionInfo.Model
		predictionConfig.TimeWindow = time.Duration(predictionInfo.Interval)
		predictionConfig.PredictLength = uint64(predictionInfo.PredictLength)
		//predictionConfig.
		pack, err := cfg.Pack(&predictionConfig)
		if err != nil {
			zap.L().Error("cfg.Pack() failed", zap.Error(err))
			return err
		}
		predictionPipelineSource.Properties = pack

		predictionPipeline.Sources = append(predictionPipeline.Sources, &predictionPipelineSource)
		piplines.Pipelines = append(piplines.Pipelines, *predictionPipeline)

	}
	return nil
}
